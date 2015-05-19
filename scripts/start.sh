#!/bin/bash

echo 'starting tweet producer'
./scripts/produce_tweets.sh &

echo 'starting web watcher'
cd web/
gulp watch
