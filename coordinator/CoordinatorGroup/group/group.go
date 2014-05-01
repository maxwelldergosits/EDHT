/*
All group memebership chagnes must go through this group
*/
package group
import (
  . "EDHT/common"
  "EDHT/utils"
  "github.com/mad293/mlog"
  )

type Group struct {
  coordinators   utils.Uint64ServerMap
  daemons   utils.Uint64ServerMap
  pendingCommits   utils.Uint64ServerMap
  nshards      uint
  nfailures    uint
  ml           mlog.MLog
  id           uint64
  newDaemonCallBack func(uint64)
}

func NewGroup(shards, failures uint,localPort,localAddress string,logger mlog.MLog, cb func(uint64), dataDir string) Group{
  newGroup := Group{}
  newGroup.coordinators = utils.NewUint64ServerMap(dataDir+"coordinators")
  newGroup.daemons = utils.NewUint64ServerMap(dataDir+"daemons")
  newGroup.id = utils.GenId(utils.GenMachineId(),true)
  newGroup.nshards = shards
  newGroup.nfailures = failures
  newGroup.ml = logger
  newGroup.pendingCommits = utils.NewUint64ServerMap(dataDir+"pendingCommits")

  me := RemoteServer{
    Address:localAddress,
    Port:localPort,
    ID:newGroup.id,
    Coordinator:true}

  newGroup.coordinators.Put(me.ID,me)
  newGroup.newDaemonCallBack = cb

  return newGroup
}

func JoinGroup(regReply RegisterReply, logger mlog.MLog, cb func(uint64),  dataDir string) (Group) {
  newGroup := Group{}
  newGroup.coordinators = utils.NewUint64ServerFromMap(dataDir+"coordinators",regReply.Coordinators)
  newGroup.daemons = utils.NewUint64ServerFromMap(dataDir+"daemons",regReply.Daemons)
  newGroup.ml = logger
  newGroup.id = regReply.ID
  newGroup.nshards = regReply.Nshards
  newGroup.nfailures = regReply.Nfailures
  newGroup.pendingCommits = utils.NewUint64ServerMap(dataDir+"pendingCommits")
  newGroup.newDaemonCallBack = cb

  return newGroup
}


func (g * Group) Delete(rs RemoteServer) {
  if rs.Coordinator {
    g.coordinators.Delete(rs.ID)
  } else {
    g.daemons.Delete(rs.ID)
  }
}

func (g * Group)GetNFailures() uint {
  return g.nfailures
}

func (g * Group)GetNShards() uint {
  return g.nshards
}

func (g * Group) GetDaemon(d uint64) RemoteServer{
  rs, _ := g.daemons.Get(d)
  return rs
}
func (g * Group) GetID() uint64{
  return g.id
}

func (g * Group) GetCoordinator(d uint64) RemoteServer{
  rs, _:= g.coordinators.Get(d)
  return rs
}

func (g * Group) Coordinators() map[uint64]RemoteServer {
  return g.coordinators.Map()
}
