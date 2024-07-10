# frozen_string_literal: true

require 'opentelemetry/sdk'
require 'opentelemetry/exporter/otlp'
require 'opentelemetry-instrumentation-all'
require 'opentelemetry/instrumentation/base'
require 'logger'

OpenTelemetry.logger = Logger.new($stdout, level: Logger::DEBUG)

OpenTelemetry::SDK.configure do |c|
  c.service_name = ENV.fetch('OTEL_SERVICE_NAME', 'ruby-otlp')
  c.use_all
end
