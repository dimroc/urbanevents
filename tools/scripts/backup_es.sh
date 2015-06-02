#!/bin/bash

cd /ruby/urbanevents/tools && bundle exec rake es:snapshot >> /var/log/cron.log 2>&1

