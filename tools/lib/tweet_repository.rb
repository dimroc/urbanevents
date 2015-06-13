class TweetRepository
  include Elasticsearch::Persistence::Repository
  include Elasticsearch::DSL

  def initialize(options={})
    client Elasticsearch::Client.new url: options[:url], log: true
    index "#{GO_ENV}-geoevents-#{Time.now.strftime("%Y%m%d-%H%M%S")}"
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

    response = self.client.search body: definition
    rval = response["aggregations"]["city_counts"]["buckets"]
    cityKeys.inject({}) do |memo, cityKey|
      memo[cityKey] = rval[cityKey]["since"]["buckets"][0]["doc_count"]
      memo
    end
  end
end
