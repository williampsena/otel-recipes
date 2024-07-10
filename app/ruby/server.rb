# frozen_string_literal: true

require 'sinatra'
require 'sinatra/required_params'
require_relative 'config/initializers/opentelemetry'
require 'prometheus/client'
require "prometheus/middleware/collector"
require "prometheus/middleware/exporter"
require './client'

prometheus = Prometheus::Client.registry
poke_image_requests = Prometheus::Client::Counter.new(:poke_image_requests,
                                                      docstring: 'A counter of HTTP pokemon image requests made')
prometheus.register(poke_image_requests)

use Prometheus::Middleware::Collector
use Prometheus::Middleware::Exporter

set :bind, '0.0.0.0'
set :port, ENV.fetch('PORT', '8002')
set :server_settings, timeout: 10
set :public_folder, 'public'

ENV['OTEL_TRACES_EXPORTER'] ||= 'console'

enable :logging

get '/pokemon/:name/image' do
  required_params :name

  client = PokeAPI.new
  response = client.pokemon(params[:name])

  if response.ok?
    image_url = response.parsed_response.dig('sprites', 'front_default')
    image_response = client.download_image(image_url)

    if image_response.ok?
      poke_image_requests.increment
      content_type 'image/png'
      return image_response.body
    end
  end

  content_type :json
  status 500
  return { messge: "👀 Oops, you can't get any pokemon photo." }.to_json
end

get '/pokemon/:name/details' do
  content_type :json
  required_params :name

  client = PokeAPI.new
  response = client.pokemon(params[:name])

  return response.body if response.ok?

  status 500
  return { messge: "👀 Oops, you can't get any pokemon photo." }.to_json
end

get '/pokemon/:name' do
  content_type :json
  required_params :name

  client = PokeAPI.new
  response = client.pokemon(params[:name])

  return response.parsed_response.slice('id', 'name').to_json if response.ok?

  status 500
  return { messge: "👀 Oops, you can't get any pokemon." }.to_json
end

get '/' do
  return '🤓 Welcome to the lightweight Pokemon API.'
end

error do
  OpenTelemetry::Trace.current_span.record_exception(env['sinatra.error'])
end
