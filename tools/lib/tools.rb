GO_ENV = ENV['GO_ENV'] || 'development'

require 'bundler'
Bundler.require(:default, GO_ENV)

require_relative 'tweet_repository.rb'
require_relative 'cloud_monitor.rb'

