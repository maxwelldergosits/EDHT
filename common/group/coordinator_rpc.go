package group

import (
  "log"
  "net/rpc"
  "net/http"
  "net"
  . "EDHT/common"
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

func coordinatorRPCstub(methodName string, ns * RemoteServer,addr string) int {


  // connect to client // TODO: make a connection caching service that will create
  // a connection or recycle an old one.
  client, err := rpc.DialHTTP("tcp", addr)
  if err != nil {
  log.Fatal("dialing:", err)
  }
  // Synchronous call
  args := ns
  var reply int
  err = client.Call(methodName, args, &reply)

  return reply

}

//***********//
//precommit
//***********//

//RPC stub//
func propseRegisterRPC(ns * RemoteServer, addr string) int {
  return coordinatorRPCstub("Coordinator.ProposeRegister",ns,addr)
}

// RPC method//
func (t * Coordinator) ProposeRegister(ns * RemoteServer, res * int) error {

  if !preCommit(*ns) {
    *res = 0 //**responding NO**//
  } else {
    *res = 1 //** responding YES**//
  }
  return nil
}

//***********//
// do commit //
//***********//

//RPC stub//
func registerRPC(ns * RemoteServer, addr string) int {
  return coordinatorRPCstub("Coordinator.Register",ns,addr)
}

// RPC method//
func (t *Coordinator) Register(ns * RemoteServer, res * int) error{

  localCommit(*ns)

  return nil

}

//*****************//
// rollback commit //
//****************//

//RPC stub//
func RollBackRegisterRPC(ns * RemoteServer, addr string) int{
  return coordinatorRPCstub("Coordinator.RollbackRegister",ns,addr)
}

// RPC method//
func (t* Coordinator) RollbackRegister(ns * RemoteServer, res * int) error {

  localAbort(*ns)

  return nil

}

