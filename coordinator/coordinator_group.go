package main

import (
  "log"
  "net/rpc"
  . "EDHT/utils"
)

var (
  localAddress string
  localPort    string
  remoteServers map[string]RemoteServer
)


func InitLocalState(alocalAddress string, alocalPort string) {

  // validate IP and port
  if !ValidateIP(alocalAddress){
    log.Panic("Error: invalid IP address: ", alocalAddress)
  }


  if !ValidatePort(alocalPort){
    log.Panic("Error: invalid port: ", alocalPort)
  }

  localAddress = alocalAddress
  localPort    = alocalPort

  remoteServers = map[string]RemoteServer{}

}

type RemoteServer struct{
  Address string
  Port string
}



func addToServers(rs * RemoteServer) {
  var ns = *rs
  remoteServers[ns.Address+":"+ns.Port] = ns
  log.Println("New coordinator at",ns.Address+":"+ns.Port)
}

// this is an RPC method so we can remotely tell a coordinator to add a new remoteserver
func (t *Member) RegisterCoordinator(ns * RemoteServer, res * int) error{

  addToServers(ns)

  *res = 1
  return nil

}

// this rpc call prompts reciever to add this new server to the group
func (t * Member) AttachRSToGroup(ns RemoteServer, res * map[string]RemoteServer) error {

  *res = AttachRSToGroup_local(ns)
  return nil

}

func AttachRSToGroup_local(ns RemoteServer) map[string]RemoteServer {

  // channel for barrier
  done := make(chan bool)
  var num = 0


  for _,v := range remoteServers {
    num++

    // spawn off other goroutines to tell other coordinators
    go func(rs RemoteServer) {


      // start connection to remote Server
      client, err := rpc.DialHTTP("tcp", rs.Address + ":" + rs.Port)
      if err != nil {
        log.Fatal("dialing:", err)
      }

      // make the rpc call
      args := &ns
      var reply int
      err = client.Call("Member.RegisterCoordinator", args, &reply)
      if err != nil {
        log.Fatal("register error:", err)
      }

      client.Close()
      done<-true
    }(v)

  }

  dest := make(map[string]RemoteServer)

  for k,v := range remoteServers {
    dest[k] = v
  }
  // locally register new server
  addToServers(&ns)

  //wait for servers to be updated.
  for i := 0; i < num; i++ {
    <-done
  }

  return dest
}



func AttachToGroup(groupAddress string, groupPort string) {


  var rs RemoteServer = RemoteServer{localAddress,localPort}
  var res map[string]RemoteServer

  // start connection to remote Server
  client, err := rpc.DialHTTP("tcp", groupAddress + ":" + groupPort)
  if err != nil {
    log.Fatal("dialing:", err)
  }
  // make the rpc call
  err = client.Call("Member.AttachRSToGroup", rs, &res)
  if err != nil {
    log.Fatal("attach error:", err)
  }
  remoteServers = res
  remoteServers[groupAddress+":"+groupPort]=RemoteServer{groupAddress,groupPort}

}

