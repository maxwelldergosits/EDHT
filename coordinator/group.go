package main

import (
  "log"
  "net/rpc"
  . "EDHT/common"
)

var (
  localAddress string
  localPort    string
  remoteServers map[string]RemoteServer
)


func addToServers(rs * RemoteServer) {
  var ns = *rs
  remoteServers[ns.Address+":"+ns.Port] = ns
  log.Println("New coordinator at",ns.Address+":"+ns.Port)
}


// this rpc call prompts reciever to add this new server to the group
func (t * Coordinator) AttachRSToGroup(ns RemoteServer, res * map[string]RemoteServer) error {

  *res = AttachRSToGroup_local(ns)
  return nil

}


func AttachRSToGroup_local(rs RemoteServer) map[string]RemoteServer {

  done := make(chan bool)
  var num = 0
  //precommit to all nodes
  for _,v := range remoteServers {

    num++

    go func(v RemoteServer) {

      addr := v.Address+":"+v.Port

      var commit = (propseRegisterRPC(&rs,addr) == 1)

      done <- commit

    }(v)

  }

  commit := true
  //wait for respsonses
  for i := 0; i < num; i++ {
    commit = (commit && <-done)
  }

  if commit {

    localCommit(rs)

    for _,v := range remoteServers {

      go func(v RemoteServer) {

        addr := v.Address+":"+v.Port

        var commit = (registerRPC(&rs,addr) == 1)

        done <- commit

      }(v)

    }

  } else {

    localAbort(rs)

    for _,v := range remoteServers {

      go func(v RemoteServer) {

        addr := v.Address+":"+v.Port

        var commit = (registerRPC(&rs,addr) == 1)

        done <- commit

      }(v)

    }

  }

  for i := 0; i < num; i++ {
    <-done
  }

  return remoteServers

}


func AttachToGroup(groupAddress string, groupPort string) {

  var rs RemoteServer = RemoteServer{localAddress,localPort,0}
  var res map[string]RemoteServer

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
  remoteServers = res
  remoteServers[groupAddress+":"+groupPort]=RemoteServer{groupAddress,groupPort,0}

}

