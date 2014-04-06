package group

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
func (t * Coordinator) AttachRSToGroup(ns RemoteServer, res * RegisterReply) error {
  *res = AttachRSToGroup_local(ns)
  return nil
}

//***********//
//precommit
//***********//


// RPC method//
func (t * Coordinator) ProposeRegister(ns * RemoteServer, res * bool) error {

  *res = preCommit(*ns)
  return nil
}

//***********//
// do commit //
//***********//


// RPC method//
func (t *Coordinator) Register(ns * RemoteServer, res * bool) error{

  localCommit(*ns)
  *res = true

  return nil

}

//*****************//
// rollback commit //
//****************//

// RPC method//
func (t* Coordinator) RollbackRegister(ns * RemoteServer, res * bool) error {

  localAbort(*ns)
  *res = true
  return nil

}

