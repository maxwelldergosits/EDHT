/*
All group memebership chagnes must go through this group
*/
package group
import (
  . "EDHT/common"
  . "EDHT/utils"
  "net/rpc"
  "log"
  )

type Group struct {
  Coordinators map[int64]RemoteServer
  Daemons      map[int64]RemoteServer
}




var defaultGroup Group
var id int64

var verboseLog func(a...interface{})
var normalLog func(a...interface{})

//Creates a new group
//Returns your id in the group

func InitGroup(vl func(a...interface{}), nl func(a...interface{})) {

  verboseLog = vl
  normalLog = nl


  pendingCommits = map[int64]RemoteServer{}
}

func CreateGroup(ip string, port string) int64 {

  self := RemoteServer{
    Address:ip,
    Port:port,
    ID:GenId(GenMachineId(),true),
    Coordinator:true}

  defaultGroup = Group{
    Coordinators:map[int64]RemoteServer{
      self.ID:self},
    Daemons:map[int64]RemoteServer{}}


  id = self.ID
  return self.ID
}

//Joins group at <ip>:<port> as a Daemon
//Handles rpc to coordinator specifed at <ip>:<port>
//Returns your id in the group
func JoinGroupAsDaemon(ip string, port string, localIP string, localPort string) int64 {

  me := RemoteServer{
    Address:localIP,
    Port:localPort,
    ID:GenMachineId(),
    Coordinator:false}


  var res RegisterReply

  // start connection to remote Server
  client, err := rpc.DialHTTP("tcp", ip + ":" + port)
  if err != nil {
    log.Fatal("dialing:", err)
  }
  // make the rpc call
  err = client.Call("Coordinator.AttachRSToGroup", me, &res)
  if err != nil {
    log.Fatal("attach error:", err)
  }

  id := res.ID
  verboseLog("id:",res.ID)

  return id

}

//Joins group at <ip>:<port> as a Coordinator
//Handles rpc to coordinator specifed at <ip>:<port>
//Returns your id in the group
func JoinGroupAsCoordinator(ip string, port string,localAddress string, localPort string) int64 {

  var rs RemoteServer = RemoteServer{localAddress,localPort,GenMachineId(),true}
  var res RegisterReply

  // start connection to remote Server
  client, err := rpc.DialHTTP("tcp", ip + ":" + port)
  if err != nil {
    log.Fatal("dialing:", err)
  }
  // make the rpc call
  err = client.Call("Coordinator.AttachRSToGroup", rs, &res)
  if err != nil {
    log.Fatal("attach error:", err)
  }
  defaultGroup.Coordinators = res.Coordinators
  defaultGroup.Daemons      = res.Daemons
  id                        = res.ID

  verboseLog("id:",res.ID)

  return id
}

