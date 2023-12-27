# Preparing environment

```shell
python3 -m venv .venv
. .venv/bin/activate 
```

# Installing packages

```shell
pip install -r requirements.txt 
```

# How to run web app?

```shell
#activate env
. .venv/bin/activate 

opentelemetry-instrument \
    --traces_exporter console,otlp \
    --metrics_exporter console,otlp \
    --logs_exporter console,otlp \
    --exporter_otlp_endpoint http://localhost:14317 \
    --service_name python \
    flask run -p $PORT --host=0.0.0.0
```