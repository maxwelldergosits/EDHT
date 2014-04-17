/*
All group memebership chagnes must go through this group
*/
package group
import (
  . "EDHT/common"
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

func NewGroup(shards, failures uint) Group {
    //TODO implement
  return Group{}
}

func JoinGroup(regReply RegisterReply) (Group, error) {
  newGroup := Group{}
  newGroup.coordinators = regReply.Coordinators
  newGroup.daemons = regReply.Daemons
  newGroup.id = regReply.ID
  newGroup.nshards = regReply.Nshards
  newGroup.nfailures = regReply.Nfailures
  newGroup.pendingCommits = make(map[uint64]RemoteServer)

  return newGroup,nil
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

func (g * Group) GetCoordinator(d uint64) RemoteServer{
  return g.coordinators[d]
}
