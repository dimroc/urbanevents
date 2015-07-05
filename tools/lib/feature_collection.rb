class FeatureCollection
  include Enumerable
  attr_reader :geojson

  def initialize(geojson_path)
    @geojson = JSON.parse(File.read(geojson_path))
    raise ArgumentError, "Not a feature collection" unless geojson["type"] == "FeatureCollection"
  end

  def clean_duplicate_coordinates!
    each do |feature|
      feature.remove_duplicate_coordinates!
    end
  end

  def each
    geojson["features"].each do |feature|
      yield Feature.new feature
    end
  end

  def as_json(options={})
    geojson.as_json(options)
  end

  class Feature
    def initialize(feature_hash)
      @hash = feature_hash
    end

    def name
      @hash["properties"]["name"]
    end

    def geometry
      @hash["geometry"]
    end

    def geometry_type
      geometry["type"]
    end

    def remove_duplicate_coordinates!
      self.send("remove_duplicate_coordinates_from_#{geometry_type.downcase}")
    end

    private

    def remove_duplicate_coordinates_from_multipolygon
      geometry["coordinates"] = geometry["coordinates"].map do |outer|
        outer.map do |inner|
          prev = nil
          cleaned = inner.map do |coordinate|
            if coordinate == prev
              prev = coordinate
              #puts "removing duplicate: #{coordinate}"
              nil
            else
              prev = coordinate
              coordinate
            end
          end.compact

          #puts "first: #{cleaned.first} last: #{cleaned.last}"
          #puts "count: #{cleaned.count}"
          cleaned
        end
      end
    end

    def remove_duplicate_coordinates_from_polygon
      geometry["coordinates"] = geometry["coordinates"].map do |polygon|
        prev = nil
        cleaned = polygon.map do |coordinate|
          if coordinate == prev
            prev = coordinate
            #puts "removing duplicate: #{coordinate}"
            nil
          else
            prev = coordinate
            coordinate
          end
        end.compact

        #puts "first: #{cleaned.first} last: #{cleaned.last}"
        #puts "count: #{cleaned.count}"
        cleaned
      end
    end
  end
end
