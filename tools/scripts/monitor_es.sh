#!/bin/bash

cd /ruby/urbanevents/tools && bundle exec rake es:report_health_last_hour >> /var/log/cron.log 2>&1

