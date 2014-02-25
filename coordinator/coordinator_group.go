package main

import (
  "log"
  "container/list"
  "net/rpc"
  . "EDHT/utils"
)

var (
  localAddress string
  localPort    string
  remoteServers list.List
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

  remoteServers = *list.New()

}

type RemoteServer struct{
  Address string
  Port string
}


func addToServers(rs * RemoteServer) {
  remoteServers.PushFront(*rs)
  var ns = *rs
  log.Println("New Coordinator at ",ns.Address+":"+ns.Port,"is now part of group")
}

// this is an RPC method so we can remotely tell a coordinator to add a new remoteserver
func (t *Member) RegisterCoordinator(ns * RemoteServer, res * int) error{

  addToServers(ns)

  *res = 1
  return nil

}

func (t *Member) RegisterCoordinatorList(ns * []RemoteServer, res * int) error{

  servers := *ns
  for i := range servers {
    addToServers(&servers[i])
  }
  *res = 1
  return nil

}

// this rpc call prompts reciever to add this new server to the group
func (t * Member) AttachRSToGroup(ns RemoteServer, res * int) error {
  // channel for barrier
  done := make(chan bool)
  num := remoteServers.Len()
//  sl := remoteServers


  for e := remoteServers.Front(); e != nil; e = e.Next() {


    // spawn off other goroutines to tell other coordinators
    log.Println("about to make rpc call to",e.Value)
    if e.Value == nil {break}
    func(rs RemoteServer) {


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
    }(e.Value.(RemoteServer))

  }


  //copy servers to array
  var numServers = remoteServers.Len()

  var sa  = make([]RemoteServer,numServers)
  var saIndex = 0
  for e:= remoteServers.Front(); e != nil; e = e.Next() {
    sa[saIndex] = e.Value.(RemoteServer)
    saIndex++
  }

  // start connection to remote Server
  client, err := rpc.DialHTTP("tcp", ns.Address + ":" + ns.Port)
  if err != nil {
    log.Fatal("dialing:", err)
  }

  // make the rpc call
  args := &sa
  var reply int
  err = client.Call("Member.RegisterCoordinatorList", args, &reply)
  if err != nil {
    log.Fatal("register error:", err)
  }

  // locally register new server
  addToServers(&ns)

  //wait for servers to be updated.
  for i := 0; i < num; i++ {
    <-done
  }

  return nil
}



func AttachToGroup(groupAddress string, groupPort string) {


  var rs RemoteServer = RemoteServer{localAddress,localPort}
  var res int

  // start connection to remote Server
  log.Println("making rpc call to", groupAddress+":"+groupPort)
  client, err := rpc.DialHTTP("tcp", groupAddress + ":" + groupPort)
  if err != nil {
    log.Fatal("dialing:", err)
  }
  // make the rpc call
  err = client.Call("Member.AttachRSToGroup", rs, &res)
  if err != nil {
    log.Fatal("attach error:", err)
  }

}

