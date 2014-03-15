package main

import (
  "log"
  . "EDHT/common"
  "flag"
  "net/rpc"
  "EDHT/utils"
)
var (
  port string
  ip   string
  groupPort string
  id int64
  groupAddress string
  verbose bool
  verboseLog func(a...interface{})
  normalLog func(a...interface{})

  local_state DaemonData
)


func registerCLA(){

  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")

  flag.BoolVar(&verbose, "verbose", false, "verbose output")
  flag.StringVar(&groupAddress, "group-address", "", "Address of any node in a group to connect to")
  flag.StringVar(&groupPort, "group-port", "", "Port of that the node in the group is on")

  flag.Parse()
}


func AttachToGroup(groupAddress string, groupPort string) {

  var rs RemoteServer = RemoteServer{ip,port,utils.GenMachineId(),false}
  var res RegisterReply

  // start connection to remote Server
  client, err := rpc.DialHTTP("tcp", groupAddress + ":" + groupPort)
  if err != nil {
    log.Fatal("dialing:", err)
  }
  // make the rpc call
  err = client.Call("Coordinator.AttachRSToGroup", rs, &res)
  if err != nil {
    log.Fatal("attach error:", err)
  }

  id = res.ID
  verboseLog("id:",res.ID)

}


func main() {

 registerCLA()

  local_state = SpawnDaemon(ip,port)

  if(verbose) {
    log.Println("port:",port)
    log.Println("ip-address:",ip)
  }

  normalLog,verboseLog = utils.GenLogger(verbose,"",true)
  AttachToGroup(groupAddress,groupPort)
  startServer(ip,port)

}
