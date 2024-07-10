package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Email counter metric
var EmailMeter = otel.Meter("email")

// Shutdown handler is responsible for finishing trace.
type ShutdownHandler func(context.Context) error

// Regular email data
type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Regular customer
type Customer struct {
	Id       string // the customer unique id
	Name     string // the customer name
	Document string // the document number
	Email    string // the customer email
}

// This function is responsible for setting up the program before it runs
func init() {
	gofakeit.Seed(0)
}

// Build redis client connection
func setupRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	return rdb
}

// Initializes the open telemetry tracer.
func setupTracer(ctx context.Context) (ShutdownHandler, error) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	tp := buildTracer(ctx, exporter)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

// Build an open telemetry tracer
func buildTracer(ctx context.Context, exporter *otlptrace.Exporter) *sdktrace.TracerProvider {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", os.Getenv("OTEL_SERVICE_NAME")),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
}

// Responsible for finalizing trace context.
func doShutdown(ctx context.Context, shutdown ShutdownHandler) {
	func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down tracer: %v", err)
		}
	}()

}

// Add span attributes values
func setupSpanValues(span trace.Span) {
	span.SetAttributes(
		attribute.String("customer.id", gofakeit.UUID()),
		attribute.String("customer.email", gofakeit.Email()),
		attribute.String("customer.password", gofakeit.Password(true, true, true, true, true, 10)),
		attribute.String("customer.vatnumber", gofakeit.SSN()),
		attribute.String("customer.credit_card", gofakeit.CreditCard().Number),
		attribute.String("db.user", gofakeit.Username()),
		attribute.String("db.password", gofakeit.Password(true, true, true, true, true, 10)),
		attribute.String("account.email", gofakeit.Email()),
	)
}

// Returns an internal server error
func writeHttpError(span trace.Span, w http.ResponseWriter, errorMessage string) {
	span.AddEvent("error",
		trace.WithAttributes(
			attribute.String("value", errorMessage),
		),
	)
	span.End()

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errorMessage))
}

// Route to generate stats for every request.
func sendEmailRoute(rdb *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("go-tracer")
		_, span := tracer.Start(r.Context(), "send-email")

		message, err := gofakeit.EmailText(&gofakeit.EmailOptions{})

		if err != nil {
			writeHttpError(span, w, fmt.Sprintf("failed to fetch random message: %v", err))
			return
		}

		customer := Customer{
			Id:       gofakeit.UUID(),
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Document: gofakeit.SSN(),
		}

		email := Email{
			From:    fmt.Sprintf("no-reply@%v", gofakeit.DomainName()),
			To:      customer.Email,
			Subject: gofakeit.BookTitle(),
			Body:    message,
		}

		span.SetAttributes(
			attribute.String("customer.id", customer.Id),
			attribute.String("customer.email", customer.Email),
			attribute.String("customer.document", customer.Document),
		)

		setupSpanValues(span)

		jsonEmail, err := json.Marshal(email)

		if err != nil {
			writeHttpError(span, w, fmt.Sprintf("failed to parse email message: %v", err))
			return
		}

		err = rdb.SPublish(r.Context(), "email", jsonEmail).Err()

		if err != nil {
			writeHttpError(span, w, fmt.Sprintf("failed to queue email message: %v", err))
			return
		}

		span.AddEvent("email",
			trace.WithAttributes(
				attribute.String("subject", email.Subject),
				attribute.String("content", email.Body),
			),
		)

		counter, _ := EmailMeter.Int64Counter(
			"send_email",
			metric.WithUnit("1"),
			metric.WithDescription("Just a email counter"),
		)

		counter.Add(r.Context(), 1, metric.WithAttributes(attribute.String("foo", "bar")))

		response := fmt.Sprintf("📨 The email was queued successfully: %v", email.Subject)

		span.AddEvent("log-message", trace.WithAttributes(attribute.String("message", message)))

		log.Println(message)

		span.End()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(response))
	}
}

func main() {
	rdb := setupRedis()
	ctx := context.Background()
	shutdown, err := setupTracer(ctx)

	if err != nil {
		log.Fatalf("failed to initialize open telemetry tracer: %v", err)
	}

	defer doShutdown(ctx, shutdown)

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(sendEmailRoute(rdb)), "SendEmail")

	http.Handle("/send-email", otelHandler)

	http.ListenAndServe(":8001", nil)
}
