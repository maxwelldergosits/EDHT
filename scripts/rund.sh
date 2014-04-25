#!/bin/bash

control_c()
# run if user hits control-c
{
  ps aux | grep 'daemon' | awk '{print $2}' | xargs kill
  exit $?
}
# trap keyboard interrupt (control-c)
trap control_c SIGINT

a=$2
for i in $(seq 1 $1);
  do
    a=$(($a + 1))
    daemon -port=$a -group-port=$3 -group-address=$4&
done
while true; do sleep 10000; done

