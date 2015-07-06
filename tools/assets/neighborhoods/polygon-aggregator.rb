#!/usr/bin/env ruby

require 'json'
require 'active_support/all'

"Aggregating #{ARGV[0]} based on name"

filename = ARGV[0]

json = JSON.parse(File.read filename)
output = {
  type: "FeatureCollection",
  features: []
}

features = json["features"]
uniq_names = features.map do |feature| feature["properties"]["name"] end
uniq_names.uniq!

uniq_features = uniq_names.map do |name|
  relevants = features.select { |f| f["properties"]["name"] == name }
  {
    type: "Feature",
    properties: relevants.first["properties"],
    geometry: {
      type: "MultiPolygon",
      coordinates: relevants.map { |r| r["geometry"]["coordinates"] }
    }
  }
end

json["features"] = uniq_features

puts json.to_json
