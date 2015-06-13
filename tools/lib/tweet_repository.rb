class TweetRepository
  include Elasticsearch::Persistence::Repository

  def initialize(options={})
    self.client = Elasticsearch::Client.new url: options[:url], log: options[:log]
    self.index = options[:index].presence || "#{GO_ENV}-geoevents-#{Time.now.strftime("%Y%m%d-%H%M%S")}"
  end

  # Set a custom document type
  type :tweet

  settings number_of_shards: 3 do
    mapping do
      indexes :createdAt, type: 'date'
      indexes :payload, analyzer: 'snowball'
      indexes :city, type: 'string', index: 'not_analyzed'
      indexes :metadata, type: 'object' do
        indexes :screenName, type: 'string'
        indexes :hashtags, type: 'string'
        indexes :mediaTypes, type: 'string', index: 'not_analyzed'
        indexes :mediaUrls, type: 'string', index: 'not_analyzed'
        indexes :expandedUrls, type: 'string', index: 'not_analyzed'
      end

      indexes :geojson, type: 'geo_shape', "tree": "quadtree", "precision": "1m"
      indexes :point, type: 'geo_point', geo_hash: true, geohash_prefix: true, geohash_precision: '1m'

      # Rest of the attributes are created lazily
    end
  end

  def copy_from(source_index)
    Tweet.gateway.index = source_index
    Tweet.gateway.client = self.client
    Tweet.find_in_batches(size: 100) do |batch|
      insert_batch(batch)
    end
  end

  def city_count_since(cityKeys, time)
    Builder.new(client).city_count_since(cityKeys, time)
  end

  private

  def insert_batch(batch)
    bulk_insertion = batch.map do |tweet|
      { index: { _id: tweet.id.to_s, data: tweet.to_hash } }
    end

    self.client.bulk({
      index: index,
      type: 'tweet',
      body: bulk_insertion
    })
  end

  class Builder
    include Elasticsearch::DSL
    attr_accessor :client

    def initialize(client)
      @client = client
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
