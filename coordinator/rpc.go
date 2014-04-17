package main

import (
  . "EDHT/common"
  "EDHT/coordinator/CoordinatorGroup/group"
  "net/rpc"
  "net"
  "net/http"
  "log"
)
var g group.Group // one piece of state, becuase of rpc

type Coordinator int

func CoordinatorStartServer(ip string, port string, gg group.Group) {

  g = gg
  reg := new(Coordinator)
  rpc.Register(reg)
  rpc.HandleHTTP()
  l,e := net.Listen("tcp",ip+":"+port)
  if e != nil {
    log.Fatal("listen error:", e)
  }
  http.Serve(l, nil)

}

// this rpc call prompts reciever to add this new server to the group
func (t * Coordinator) AttachRSToGroup(ns RemoteServer, res * RegisterReply) error {

  *res = g.AttachRSToGroup_local(ns)
  return nil

}

//rpc to precommit
func (t * Coordinator) ProposeRegister(ns * RemoteServer, res * bool) error {

  *res = g.PreCommit(*ns)
  return nil

}

//rpc for commit
func (t * Coordinator) Register(ns * RemoteServer, res * bool) error{

  g.LocalCommit(*ns)
  *res = true
  return nil

}

//RPC to abort
func (t * Coordinator) RollbackRegister(ns * RemoteServer, res * bool) error {

  g.LocalAbort(*ns)
  *res = true
  return nil

}
