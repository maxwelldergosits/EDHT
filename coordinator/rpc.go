package main

import (
  . "EDHT/common"
  "net/rpc"
  "net"
  "net/http"
  "log"
)

type Coordinator int

func CoordinatorStartServer(ip string, port string) {

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
func (t * Coordinator) AttachRSToGroup(ns RemoteServer, res *ConnectReply) error {

  rr := gc.Gms.AttachRSToGroup_local(ns)
  scs := gc.GetPartitions().GetShardCopies()
  *res = ConnectReply{rr,scs}
  return nil

}

//rpc to precommit
func (t * Coordinator) ProposeRegister(ns * RemoteServer, res * bool) error {

  *res = gc.Gms.PreCommit(*ns)
  return nil

}

//rpc for commit
func (t * Coordinator) Register(ns * RemoteServer, res * bool) error{

  gc.Gms.LocalCommit(*ns)
  *res = true
  return nil

}

//RPC to abort
func (t * Coordinator) RollbackRegister(ns * RemoteServer, res * bool) error {

  gc.Gms.LocalAbort(*ns)
  *res = true
  return nil

}
