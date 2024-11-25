#!/bin/bash
echo "Removing IP Addresses"
ip addr flush eth0

trap "echo Shutting down; exit 0" SIGTERM SIGINT SIGKILL
/bin/sleep infinity &
wait
