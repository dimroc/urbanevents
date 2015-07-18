class HoodRepository
  include Elasticsearch::Persistence::Repository

  def initialize(options={})
    self.client = Elasticsearch::Client.new url: options[:url], log: options[:log]
    self.index = options[:index].presence || "#{GO_ENV}-geoevents-hoods"
  end

  type :hood

  settings number_of_shards: 3 do
    mapping do
      indexes :shape, type: 'geo_shape'

      # Rest of the attributes are created lazily
    end
  end

  def import_from(geojson_path)
    feature_collection = FeatureCollection.new(geojson_path)
    feature_collection.clean_duplicate_coordinates!
    feature_collection.each do |feature|
      puts "Adding #{feature.name}"
      self.client.index({
        index: index,
        type: type,
        body: { id: feature.name, shape: feature.geometry, file: File.basename(geojson_path) }
      })
    end
  end

  def update_mapping!
    client.indices.put_mapping({
      index: index,
      type: type,
      body: mappings.to_hash
    })
  end
end
