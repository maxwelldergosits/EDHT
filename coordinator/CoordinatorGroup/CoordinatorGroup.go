package CoordinatorGroup

import (
  "EDHT/coordinator/CoordinatorGroup/partition"
  "EDHT/coordinator/CoordinatorGroup/group"
  "EDHT/common/rpc_stubs"
  . "EDHT/common"
  "github.com/mad293/mlog"
  "EDHT/utils"
)
type CoordinatorGroup struct {
  Pts *partition.PartitionSet
  Gms group.Group
  ml mlog.MLog
}


func NewCoodinatorGroup(nshards,failures int, port,addr string, logger mlog.MLog,cb func(uint64),dataDir string) CoordinatorGroup{

  newG := group.NewGroup(uint(nshards),uint(failures),port,addr,logger,cb,dataDir)

  del := &PD{newG,logger}
  pts := partition.MakeKeySpace(nshards,del)

  return CoordinatorGroup{pts,newG,logger}
}

func ConnectToGroup(groupAddress, groupPort, address, port string, logger mlog.MLog,cb func(uint64),dataDir string) (CoordinatorGroup, error) {

  mid := utils.GenMachineId()
  rs := RemoteServer{
    groupAddress,
    groupPort,
    0,
    false}

  rr,pr, err := rpc_stubs.AttachToGroupRPC(true,address,port,mid,rs)
  if (err != nil) {
    return CoordinatorGroup{},err
  }
  g := group.JoinGroup(rr,logger,cb,dataDir)
  del := &PD{g,logger}
  ps := partition.MakePartitionSet(pr,del)
  return CoordinatorGroup{
    ps,
    g,logger},nil
}

func (cg * CoordinatorGroup) UpdatePartitions(diffs []partition.Diff, newPTS * partition.PartitionSet) (error){
  // Two phase commit
  updateID := utils.GetTimeNano()
  var succ bool = false

  localPreCommit := func() {
    cg.Pts.PreCommit(*newPTS,updateID)
  }

  localCommit := func() {
    succ = cg.Pts.ApplyCopyDiffs(diffs)
    cg.Pts.Commit(updateID)
  }
  localAbort := func() {
    cg.Pts.Abort(updateID)
  }
  remotePreCommit := func(rs RemoteServer) (bool,error) {
    return partition.PreCommitPartition(*newPTS,updateID,rs)
  }
  remoteCommit := func(rs RemoteServer) (map[string]string) {
    partition.CommitPartition(updateID,rs)
    return nil
  }
  remoteAbort := func(rs RemoteServer) (map[string]string) {
    partition.AbortPartition(updateID,rs)
    return nil
  }
  failure := func(rs RemoteServer) {
    cg.Gms.Delete(rs)
  }

  acceptors := cg.Gms.Coordinators()
  tpc := utils.InitTPC(acceptors,cg.Gms.GetID(),
                  localPreCommit,localCommit,localAbort,
                  remotePreCommit,remoteCommit,remoteAbort,
                  failure,true)

  err, _ := tpc.Run()
  if (err == nil && succ) {
    cg.Pts.ApplyDeleteDiffs(diffs)
  }
  return err
}

func (cg * CoordinatorGroup) GetPartitions() (*partition.PartitionSet) {
  return cg.Pts
}
