package main

import (
  "EDHT/coordinator/partition"
  . "EDHT/common"
)
type CoordinatorGroup struct {
}


func (cg * CoordinatorGroup) Add(c RemoteServer) (bool,error) {
  return false,nil
}


func (cg * CoordinatorGroup) UpdatePartitions(diffs []partition.Diff) (bool,error){
  return false,nil
}

func (cg * CoordinatorGroup) getPartitions() (partition.PartitionSet) {
  ps := partition.PartitionSet{}
  return ps
}
