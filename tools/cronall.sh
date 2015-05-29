#!/bin/bash
cd /ruby/urbanevents/tools/scripts
. ../set_env.sh /tmp/.env
list=`ls *.sh`
for file in $list
do
   ./$file
done
