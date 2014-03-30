package rpc_stubs
import (
  . "EDHT/common"
  "log"
  "net/rpc"
)
func RollBackRegisterRPC(ns * RemoteServer, addr string) int{
  return coordinatorRPCstub("Coordinator.RollbackRegister",ns,addr)
}

//RPC stub//
func RegisterRPC(ns * RemoteServer, addr string) int {
  return coordinatorRPCstub("Coordinator.Register",ns,addr)
}

//RPC stub//
func PropseRegisterRPC(ns * RemoteServer, addr string) int {
  return coordinatorRPCstub("Coordinator.ProposeRegister",ns,addr)
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

func AttachToGroupRPC(rs RemoteServer,addr string) RegisterReply {

  var res RegisterReply
  // start connection to remote Server
  client, err := rpc.DialHTTP("tcp", addr)
  if err != nil {
    log.Fatal("dialing:", err)
  }
  // make the rpc call
  err = client.Call("Coordinator.AttachRSToGroup", rs, &res)
  if err != nil {
    log.Fatal("attach error:", err)
  }
  return res
}
