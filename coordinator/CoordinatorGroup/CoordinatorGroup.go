package CoordinatorGroup

import (
  "EDHT/coordinator/CoordinatorGroup/partition"
  "EDHT/coordinator/CoordinatorGroup/group"
  "EDHT/common/rpc_stubs"
  "github.com/mad293/mlog"
  "EDHT/utils"
)
type CoordinatorGroup struct {
  Pts *partition.PartitionSet
  Gms group.Group
}


func NewCoodinatorGroup(nshards,failures int, port,addr string, logger mlog.MLog) CoordinatorGroup{

  newG := group.NewGroup(uint(nshards),uint(failures),port,addr,logger)

  del := &PD{newG}
  pts := partition.MakeKeySpace(nshards,del)

  return CoordinatorGroup{pts,newG}
}

func ConnectToGroup(groupAddress, groupPort, address, port string, logger mlog.MLog) (CoordinatorGroup, error) {

  mid := utils.GenMachineId()

  rr,pr, err := rpc_stubs.AttachToGroupRPC(true,address,port,mid,address+":"+groupPort)
  if (err != nil) {
    return CoordinatorGroup{},err
  }
  g := group.JoinGroup(rr,logger)
  del := &PD{g}
  ps := partition.MakePartitionSet(pr,del)
  return CoordinatorGroup{
    ps,
    g},nil
}

func (cg * CoordinatorGroup) UpdatePartitions(diffs []partition.Diff, newPTS * partition.PartitionSet) (error){
  // Two phase commit
  return nil
}

func (cg * CoordinatorGroup) GetPartitions() (*partition.PartitionSet) {
  return cg.Pts
}
