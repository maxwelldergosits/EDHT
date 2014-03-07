package main

import (
  "log"
  "net/rpc"
  . "EDHT/common"
  "EDHT/utils"
)

var (
  localAddress        string
  localPort           string
  remoteCoordinators  map[int64]RemoteServer
  remoteDaemons       map[int64]RemoteServer
  pendingCommits      map[int64]RemoteServer
)




// this rpc call prompts reciever to add this new server to the group
func (t * Coordinator) AttachRSToGroup(ns RemoteServer, res * RegisterReply) error {

  *res = AttachRSToGroup_local(ns)
  return nil

}


func preCommit(rs RemoteServer) int {
  verboseLog("precommiting:",rs)

  pendingCommits[rs.ID] = rs

  return 1
}

func localCommit(rs RemoteServer) int {
  verboseLog("commiting:",rs)
  if rs.Coordinator {
    remoteCoordinators[rs.ID]=rs
  } else {
    remoteDaemons[rs.ID]=rs
  }
  delete(pendingCommits,rs.ID)

  return 1
}

func localAbort(rs RemoteServer) int {
  verboseLog("aborting:",rs)
  delete(pendingCommits,rs.ID)
  return 1
}

func AttachRSToGroup_local(rs RemoteServer) RegisterReply {
  verboseLog("attaching to:",rs)

  done := make(chan bool)
  var num = 0
  //precommit to all nodes
  for _,v := range remoteCoordinators {

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

    for _,v := range remoteCoordinators {

      go func(v RemoteServer) {

        addr := v.Address+":"+v.Port

        var commit = (registerRPC(&rs,addr) == 1)

        done <- commit

      }(v)

    }

  } else {

    localAbort(rs)

    for _,v := range remoteCoordinators {

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

  if rs.Coordinator {
    return RegisterReply{remoteCoordinators,remoteDaemons,rs.ID}
  }else {
    return RegisterReply{nil,nil,rs.ID}
  }

}


func AttachToGroup(groupAddress string, groupPort string) {

  var rs RemoteServer = RemoteServer{localAddress,localPort,utils.GenMachineId(),true}
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
  remoteCoordinators = res.Coordinators
}

