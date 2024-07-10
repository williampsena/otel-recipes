# Preparing environment

```shell
asdf install
```

# Installing packages

```shell
pnpm install

# or

npm install
```

# How to run web app?

```shell
export PORT=8003
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:14318
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:14318/v1/traces
export OTEL_RESOURCE_ATTRIBUTES=team=dev,cluster-name=local,env=dev
export OTEL_SERVICE_NAME=node-otlp
export OTEL_TRACES_EXPORTER=otlp
export OTEL_LOG_LEVEL=info

npm start
```