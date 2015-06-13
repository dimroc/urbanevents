class Tweet
  include Elasticsearch::Persistence::Model
  attr_reader :attributes

  def initialize(attributes={})
    @attributes = attributes
  end

  def to_hash
    @attributes
  end
end
