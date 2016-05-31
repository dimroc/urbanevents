# Urban Events

Tool to record and search tweeted media across cities. http://urbanevents.dimroc.com/?q=graffiti

## High Level

1. City Recorder: Classifies tweets from Twitter's [Public Streaming API](https://dev.twitter.com/streaming/reference/post/statuses/filter)
with a neighborhood by running Elasticsearch (ES) Geospatial percolations against an index of
[city neighborhood GeoJSON files](https://github.com/dimroc/neighborhoods).
2. City Web: Searches across cities using Elasticsearch's `top_hits` metric aggregator and displays results in React JS.
3. Currently listening to NYC, London, Paris, Austin, Miami, and Los Angeles.

## City Recorder (cityservice/)

* Listens to geotagged tweets in real-time with a Golang service
* Unfurls instagram links to get actual media
* Uses Elasticsearch percolator against an index of geoshapes to classify tweets
* Pushes events using Server Side Events (SSE)
* Golang test framework [Convey](http://goconvey.co/) and mock framework [GoMock](https://github.com/golang/mock).
* http://nyc.urbanevents.dimroc.com
* http://london.urbanevents.dimroc.com

## City Web

* Isomorphic javascript using [Golang with Duktape](https://github.com/dimroc/urbanevents/tree/master/cityweb) for server side JS rendering
* React JS with Redux for the client
* Webpack with ES6, HMR, and all that good stuff
* Elasticsearch Aggregation queries to search across cities
* http://urbanevents.dimroc.com

## Kibana

* Useful for adhoc Elasticsearch queries and analytics.

## Deployment

* Docker all day
* Links and Volumes all day
* [Docker Cloud](http://cloud.docker.com) (formerly [tutum.co](http://www.tutum.co)) to deploy and host all those docker containers
* Extensive use of [Stackfile](https://github.com/dimroc/urbanevents/blob/master/Stackfile.example).

## To Do

1. Ability to play videos
2. More Design Love
3. Image classification using machine learning or developer friendly service like 
  [https://imagga.com/](https://imagga.com/) - But can I afford it?

