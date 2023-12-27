#!/bin/bash

PORT=${PORT:-8000}

echo "Server listening at port $PORT"

opentelemetry-instrument \
    --traces_exporter console,otlp \
    --metrics_exporter console,otlp \
    --logs_exporter console,otlp \
    --exporter_otlp_endpoint $OTEL_EXPORTER_OTLP_ENDPOINT \
    --service_name python \
    flask run -p $PORT --host=0.0.0.0
