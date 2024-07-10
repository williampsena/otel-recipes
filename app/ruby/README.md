# Preparing environment

```shell
asdf install
```

# Installing packages

```shell
bundle install
```

# How to run web app?

```shell
export PORT=8002
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:14318/v1/traces
export OTEL_RESOURCE_ATTRIBUTES=team=dev,cluster-name=local,env=dev
export OTEL_SERVICE_NAME=ruby-otlp
export OTEL_TRACES_EXPORTER=console,otlp

ruby server.rb
```