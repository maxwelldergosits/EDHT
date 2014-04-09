/*
All group memebership chagnes must go through this group
*/
package group
import (
  . "EDHT/common"
  . "EDHT/utils"
  "EDHT/common/rpc_stubs"
  "mlog"
  )

type Group struct {
  Coordinators map[uint64]RemoteServer
  Daemons      map[uint64]RemoteServer
  Nshards      uint
  Nfailures    uint
}




var defaultGroup Group
var id uint64
var ml mlog.MLog

var newDaemonCallback func(d uint64)

//Creates a new group
//Returns your id in the group

func InitGroup(vl mlog.MLog,newDaemon func(d uint64)) {

  ml  = vl
  newDaemonCallback = newDaemon


  pendingCommits = map[uint64]RemoteServer{}
}

func GetDaemon(d uint64) RemoteServer{
  return defaultGroup.Daemons[d]
}

func GetCoordinator(d uint64) RemoteServer{
  return defaultGroup.Coordinators[d]
}

func DeleteDaemon(d uint64) {
  ml.VPrintln("debug","Deleting Coordinator",d)
  delete(defaultGroup.Daemons,d)
}
func DeleteCoordinator(d uint64) {
  ml.VPrintln("debug","Deleting Coordinator",d)
  delete(defaultGroup.Coordinators,d)
}

func GetLocalID() uint64 {
  return id
}
func GetNFailures() uint {
  return defaultGroup.Nfailures
}
func GetNShards() uint {
  return defaultGroup.Nshards
}

func CreateGroup(ip string, port string,nshards uint, nfailures uint) uint64 {

  self := RemoteServer{
    Address:ip,
    Port:port,
    ID:GenId(GenMachineId(),true),
    Coordinator:true}

  defaultGroup = Group{
    Coordinators:map[uint64]RemoteServer{
      self.ID:self},
    Daemons:map[uint64]RemoteServer{},
    Nshards:nshards,
    Nfailures:nfailures}


  id = self.ID
  return self.ID
}

//Joins group at <ip>:<port> as a Daemon
//Handles rpc to coordinator specifed at <ip>:<port>
//Returns your id in the group
func JoinGroupAsDaemon(ip string, port string, localIP string, localPort string) uint64 {

  me := RemoteServer{
    Address:localIP,
    Port:localPort,
    ID:GenMachineId(),
    Coordinator:false}


  var res RegisterReply

  res = rpc_stubs.AttachToGroupRPC(me,ip+":"+port)

  id := res.ID
  ml.VPrintln("debug","id:",res.ID)

  return id

}

//Joins group at <ip>:<port> as a Coordinator
//Handles rpc to coordinator specifed at <ip>:<port>
//Returns your id in the group
func JoinGroupAsCoordinator(ip string, port string,localAddress string, localPort string) Group {

  var rs RemoteServer = RemoteServer{localAddress,localPort,GenMachineId(),true}
  var res RegisterReply

  res = rpc_stubs.AttachToGroupRPC(rs,ip+":"+port)

  defaultGroup = Group{nil,nil,0,0}
  defaultGroup.Coordinators = res.Coordinators
  defaultGroup.Daemons      = res.Daemons
  defaultGroup.Nshards      = res.Nshards
  defaultGroup.Nfailures    = res.Nfailures
  id                        = res.ID

  ml.VPrintln("debug","id:",res.ID)

  return defaultGroup
}

