#!/bin/bash
while true
do
    dig gravity.beryju.io. @127.0.0.1 -p 53
    sleep $[ ( $RANDOM % 5 )  + 1 ]
done
