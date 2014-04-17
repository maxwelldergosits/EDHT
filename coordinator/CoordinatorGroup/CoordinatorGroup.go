package CoordinatorGroup

import (
  "EDHT/coordinator/CoordinatorGroup/partition"
  "EDHT/coordinator/CoordinatorGroup/group"
  "EDHT/common/rpc_stubs"
  "github.com/mad293/mlog"
  "EDHT/utils"
)
type CoordinatorGroup struct {
  pts *partition.PartitionSet
  gms group.Group
}


func NewCoodinatorGroup(nshards,failures int, logger mlog.MLog) CoordinatorGroup{

  return CoordinatorGroup{}
}

func ConnectToGroup(groupAddress, groupPort, address, port string) (CoordinatorGroup, error) {

  mid := utils.GenMachineId()

  rr,pr, err := rpc_stubs.AttachToGroupRPC(true,address,port,mid,address+":"+groupPort)
  if (err != nil) {
    return CoordinatorGroup{},err
  }
  g, nil := group.JoinGroup(rr)
  ps := partition.MakePartitionSet(pr,new(PD))
  return CoordinatorGroup{
    ps,
    g},nil
}

func (cg * CoordinatorGroup) UpdatePartitions(diffs []partition.Diff) (error){
  // Two phase commit
  return nil
}

func (cg * CoordinatorGroup) GetPartitions() (*partition.PartitionSet) {
  return cg.pts
}
