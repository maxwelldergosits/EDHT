#!/bin/bash
a=$2
for i in $(seq 1 $1);
  do
    a=$(($a + 1))
    daemon -port=$a -group-port=$3 -group-address=$4&
done


