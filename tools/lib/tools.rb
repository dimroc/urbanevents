GO_ENV = ENV['GO_ENV'] || 'development'

require 'bundler'
Bundler.require(:default, GO_ENV)

require 'elasticsearch/persistence/model'

require_relative 'tweet.rb'
require_relative 'tweet_repository.rb'
require_relative 'cloud_monitor.rb'

