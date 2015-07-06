GO_ENV = ENV['GO_ENV'] || 'development'

require 'bundler'
Bundler.require(:default, GO_ENV)

require 'elasticsearch/persistence/model'

require_relative 'geoevent.rb'
require_relative 'feature_collection.rb'
require_relative 'hood_repository.rb'
require_relative 'geoevent_repository.rb'
require_relative 'cloud_monitor.rb'

