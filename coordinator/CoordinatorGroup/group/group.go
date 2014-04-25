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
  coordinators   map[uint64]RemoteServer
  daemons        map[uint64]RemoteServer
  pendingCommits map[uint64]RemoteServer
  nshards      uint
  nfailures    uint
  ml           mlog.MLog
  id           uint64
  newDaemonCallBack func(uint64)
}

func NewGroup(shards, failures uint,localPort,localAddress string,logger mlog.MLog, cb func(uint64)) Group{
  newGroup := Group{}
  newGroup.coordinators = make(map[uint64]RemoteServer)
  newGroup.daemons = make(map[uint64]RemoteServer)
  newGroup.id = utils.GenId(utils.GenMachineId(),true)
  newGroup.nshards = shards
  newGroup.nfailures = failures
  newGroup.ml = logger
  newGroup.pendingCommits = make(map[uint64]RemoteServer)

  me := RemoteServer{
    Address:localAddress,
    Port:localPort,
    ID:newGroup.id,
    Coordinator:true}

  newGroup.coordinators[me.ID] = me
  newGroup.newDaemonCallBack = cb

  return newGroup
}

func JoinGroup(regReply RegisterReply, logger mlog.MLog, cb func(uint64)) (Group) {
  newGroup := Group{}
  newGroup.coordinators = regReply.Coordinators
  newGroup.daemons = regReply.Daemons
  newGroup.ml = logger
  newGroup.id = regReply.ID
  newGroup.nshards = regReply.Nshards
  newGroup.nfailures = regReply.Nfailures
  newGroup.pendingCommits = make(map[uint64]RemoteServer)
  newGroup.newDaemonCallBack = cb

  return newGroup
}


func (g * Group) Delete(rs RemoteServer) {
  if rs.Coordinator {
    delete(g.coordinators,rs.ID)
  } else {
    delete(g.daemons,rs.ID)
  }
}

func (g * Group)GetNFailures() uint {
  return g.nfailures
}

func (g * Group)GetNShards() uint {
  return g.nshards
}

func (g * Group) GetDaemon(d uint64) RemoteServer{
  return g.daemons[d]
}
func (g * Group) GetID() uint64{
  return g.id
}

func (g * Group) GetCoordinator(d uint64) RemoteServer{
  return g.coordinators[d]
}

func (g * Group) Coordinators() map[uint64]RemoteServer {
  return g.coordinators
}
