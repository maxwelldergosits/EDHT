#!/bin/bash

num_coordinators=$1
num_shards=$2
start_port=1457

control_c()
# run if user hits control-c
{
  ps aux | grep 'coordinator' | awk '{print $2}' | xargs kill
  ps aux | grep 'daemon' | awk '{print $2}' | xargs kill
  exit $?
}
# trap keyboard interrupt (control-c)
trap control_c SIGINT

coordinator -shards=$num_shards -recalc-time=$3&
sleep 1

for i in $(seq 1 $num_coordinators);
  do
    coordinator -port=$start_port -connect-to-group -group-port=1456 -group-address=127.0.0.1 -recalc-time=$3&
    start_port=$(($start_port + 1))
    sleep .2
done
sleep 2
echo "starting 4* $num_shards daemons"
a=4000
for i in $(seq 1 $num_shards);
  do
    daemon -port=$a -group-port=1456 -group-address=127.0.0.1&
    sleep .2
    daemon -port=$(($a+1)) -group-port=1456 -group-address=127.0.0.1&
    sleep .2
    daemon -port=$(($a+2)) -group-port=1456 -group-address=127.0.0.1&
    sleep .2
    daemon -port=$(($a+3)) -group-port=1456 -group-address=127.0.0.1&
    sleep .2
    a=$(($a + 4))
done
echo "done making daemons"
while true; do sleep 10000; done

