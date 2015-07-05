#!/usr/bin/env ruby

require 'json'
require 'active_support/all'

"Cleaning #{ARGV[0]}"

filename = ARGV[0]

json = JSON.parse(File.read filename)
output = {
  type: "FeatureCollection",
  features: []
}

#{
  #"type": "FeatureCollection",
  #"features": [
    #{
      #"type": "Feature",
      #"properties": {
        #"name": "Bourse",
        #"cartodb_id": 2,
        #"created_at": "2013-02-26T07:07:16.384Z",
        #"updated_at": "2013-02-26T18:36:18.682Z"
      #},
      #"geometry": {
        #"type": "MultiPolygon",

json.each do |entry|
  reformatted = {
    type: "Feature",
    properties: entry.except("geoJson"),
    geometry: entry["geoJson"]
  }

  output[:features] << reformatted
end

puts output.to_json
