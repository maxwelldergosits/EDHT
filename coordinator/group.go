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
<<<<<<< HEAD
  id                  int64
)



func addToServers(rs * RemoteServer) {
  var ns = *rs
  remoteServers[ns.Address+":"+ns.Port] = ns
  log.Println("New coordinator at",ns.Address+":"+ns.Port)
}
=======
)


>>>>>>> two phase commit is done?


// this rpc call prompts reciever to add this new server to the group
func (t * Coordinator) AttachRSToGroup(ns RemoteServer, res * RegisterReply) error {

  *res = AttachRSToGroup_local(ns)
  return nil

}


func preCommit(rs RemoteServer) int {
<<<<<<< HEAD

  verboseLog("precommiting:",rs)
=======
  verboseLog("precommiting:",rs)

>>>>>>> two phase commit is done?
  pendingCommits[rs.ID] = rs

  return 1
}

func localCommit(rs RemoteServer) int {
<<<<<<< HEAD

  verboseLog("commiting:",rs)

=======
  verboseLog("commiting:",rs)
>>>>>>> two phase commit is done?
  if rs.Coordinator {
    remoteCoordinators[rs.ID]=rs
  } else {
    remoteDaemons[rs.ID]=rs
  }
<<<<<<< HEAD

=======
>>>>>>> two phase commit is done?
  delete(pendingCommits,rs.ID)

  return 1
}

func localAbort(rs RemoteServer) int {
<<<<<<< HEAD

=======
>>>>>>> two phase commit is done?
  verboseLog("aborting:",rs)
  delete(pendingCommits,rs.ID)
  return 1
}

func AttachRSToGroup_local(rs RemoteServer) RegisterReply {
<<<<<<< HEAD
  verboseLog("attaching:",rs)

  mid := rs.ID
  rs.ID = utils.GenId(mid,rs.Coordinator)
=======
  verboseLog("attaching to:",rs)
>>>>>>> two phase commit is done?

  done := make(chan bool)
  var num = 0

  //precommit to all nodes
  for _,v := range remoteCoordinators {

    num++

    go func(v RemoteServer) {

      addr := v.Address+":"+v.Port

      done <- (propseRegisterRPC(&rs,addr) == 1)

    }(v)

  }

  commit := true
  //wait for respsonses
  for i := 0; i < num; i++ {
    commit = (commit && <-done)
  }

  if commit {

    localCommit(rs)

<<<<<<< HEAD
<<<<<<< HEAD
    for _,v := range remoteCoordinators {

      if v == rs || v.ID == id {num--; continue} // we shouldn't wait for a response from
=======
    for _,v := range remoteServers {
>>>>>>> work on the two phase commit
=======
    for _,v := range remoteCoordinators {
>>>>>>> two phase commit is done?

      go func(v RemoteServer) {

        addr := v.Address+":"+v.Port

        var commit = (registerRPC(&rs,addr) == 1)

        done <- commit

      }(v)

    }

  } else {

    localAbort(rs)

<<<<<<< HEAD
<<<<<<< HEAD
    for _,v := range remoteCoordinators {
      if v == rs {continue}
=======
    for _,v := range remoteServers {
=======
    for _,v := range remoteCoordinators {
>>>>>>> two phase commit is done?

>>>>>>> work on the two phase commit
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
  // make the rpc call
  err = client.Call("Coordinator.AttachRSToGroup", rs, &res)
  if err != nil {
    log.Fatal("attach error:", err)
  }
<<<<<<< HEAD

  remoteCoordinators = res.Coordinators
  remoteDaemons      = res.Daemons
  id = res.ID
  verboseLog("id:",res.ID)

=======
  remoteCoordinators = res.Coordinators
>>>>>>> two phase commit is done?
}

