from random import randint
from flask import Flask

from opentelemetry import trace

tracer = trace.get_tracer("python.tracer")

app = Flask(__name__)


@app.route("/")
def hello():
    __send_metrics__()

    return "The metrics are delivered."


def __send_metrics__():
    with tracer.start_as_current_span("customer") as span:
        id = randint(1, 1000)

        span.set_attribute("customer.id", id)
        span.set_attribute("customer.email", "foo@bar.com")
        span.set_attribute("customer.password", "foo-bar-123456789")
        span.set_attribute("customer.vatnumber", "0000000191")
        span.set_attribute("customer.credit_card", "371449635398431")