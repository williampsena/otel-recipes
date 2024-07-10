# frozen_string_literal: true

require 'httparty'
require_relative 'config/initializers/opentelemetry'

# This class is Pokemon Http client implementation
class PokeAPI
  include HTTParty
  default_timeout 5

  base_uri 'https://pokeapi.co/api/v2'

  def pokemon(name)
    self.class.get("/pokemon/#{name}")
  end

  def download_image(image_url)
    HTTParty.get(image_url)
  end
end
