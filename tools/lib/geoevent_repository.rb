class GeoeventRepository
  include Elasticsearch::Persistence::Repository
  attr_reader :builder

  def initialize(options={})
    self.client = Elasticsearch::Client.new url: options[:url], log: options[:log]
    self.index = options[:index].presence || "#{GO_ENV}-geoevents-#{Time.now.strftime("%Y%m%d-%H%M%S")}"
    @builder = Builder.new(client)
  end

  # Set a custom document type
  type :geoevent

  settings number_of_shards: 3 do
    mapping do
      indexes :createdAt, type: 'date'
      indexes :text, analyzer: 'snowball'
      indexes :text_fr, analyzer: 'snowball', language: 'French'
      indexes :city, type: 'string', index: 'not_analyzed'
      indexes :username, type: 'string', index: 'not_analyzed'
      indexes :fullName, type: 'string'
      indexes :place, type: 'string'
      indexes :service, type: 'string', index: 'not_analyzed'

      indexes :mediaType, type: 'string', index: 'not_analyzed'
      indexes :thumbnailUrl, type: 'string', index: 'no'
      indexes :mediaUrl, type: 'string', index: 'no'
      indexes :link, type: 'string', index: 'no'

      indexes :hashtags, type: 'string', index: 'not_analyzed'

      indexes :geojson, type: 'geo_shape'
      indexes :point, type: 'geo_point', geo_hash: true, geohash_prefix: true, geohash_precision: '1m'

      indexes :neighborhoods, type: 'string', index: 'not_analyzed'

      # Rest of the attributes are created lazily
    end
  end

  def copy_from(source_index)
    Geoevent.gateway.index = source_index
    Geoevent.gateway.client = self.client
    Geoevent.find_in_batches(size: 100) do |batch|
      insert_batch(batch)
    end
  end

  def register_percolator(city, geojson_path)
    feature_collection = FeatureCollection.new(geojson_path)
    feature_collection.clean_duplicate_coordinates!
    feature_collection.each do |feature|
      self.client.index({
        index: index,
        type: '.percolator',
        id: city.to_s + "," + feature.name,
        body: { query: { geo_shape: { geojson: { shape: feature.geometry } } } }
      })
    end
  end

  private

  def insert_batch(batch)
    bulk_insertion = batch.map do |geoevent|
      { index: { _id: geoevent.id.to_s, data: geoevent.to_hash } }
    end

    self.client.bulk({
      index: index,
      type: type,
      body: bulk_insertion
    })
  end

  class Builder
    include Elasticsearch::DSL
    attr_accessor :client

    def initialize(client)
      @client = client
    end

    def city_count_for_service_since(cityKeys, service, time)
      city_filters = cityKeys.inject({}) do |memo, cityKey|
        memo[cityKey] = {
          bool: { must: [
            { term: { city: cityKey } },
            { term: { service: service } }
          ]}
        }
        memo
      end

      definition = search do
        size 0
        aggregation :city_counts do
          filters do
            filters city_filters
            aggregation :since do
              date_range do
                field :createdAt
                ranges [
                  { from: time }
                ]
              end
            end
          end
        end
      end

      response = client.search body: definition
      rval = response["aggregations"]["city_counts"]["buckets"]
      cityKeys.inject({}) do |memo, cityKey|
        memo[cityKey] = rval[cityKey.to_s]["since"]["buckets"][0]["doc_count"]
        memo
      end
    end

    def city_count_since(cityKeys, time)
      city_filters = cityKeys.inject({}) do |memo, cityKey|
        memo[cityKey] = {terms: { city: [cityKey] } }
        memo
      end

      definition = search do
        size 0
        aggregation :city_counts do
          filters do
            filters city_filters
            aggregation :since do
              date_range do
                field :createdAt
                ranges [
                  { from: time }
                ]
              end
            end
          end
        end
      end

      response = client.search body: definition
      rval = response["aggregations"]["city_counts"]["buckets"]
      cityKeys.inject({}) do |memo, cityKey|
        memo[cityKey] = rval[cityKey]["since"]["buckets"][0]["doc_count"]
        memo
      end
    end
  end
end
