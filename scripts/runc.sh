#!/bin/bash

control_c()
# run if user hits control-c
{
  ps aux | grep 'coordinator' | awk '{print $2}' | xargs kill
  exit $?
}
# trap keyboard interrupt (control-c)
trap control_c SIGINT


a=$2
for i in $(seq 1 $1);
  do
    a=$(($a + 1))
    coordinator -port=$a -connect-to-group -group-port=$3 -group-address=$4&
    sleep .2
done
while true; do sleep 10000; done

