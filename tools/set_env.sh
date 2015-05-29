#!/bin/bash

# Proper use:
# source ./bin/set_env .env
if [ -z "$1" ]
  then
    echo "No .env supplied"
  else
    for line in $(cat "$1"); do echo $line; export $line; done;
fi
