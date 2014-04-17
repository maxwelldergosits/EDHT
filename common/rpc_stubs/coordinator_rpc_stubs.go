package rpc_stubs
import (
  . "EDHT/common"
  "net/rpc"
)
func RollBackRegisterRPC(ns * RemoteServer, addr string) (bool,error){
  return coordinatorRPCstub("Coordinator.RollbackRegister",ns,addr)
}

//RPC stub//
func RegisterRPC(ns * RemoteServer, addr string) (bool,error) {
  return coordinatorRPCstub("Coordinator.Register",ns,addr)
}

//RPC stub//
func PropseRegisterRPC(ns * RemoteServer, addr string) (bool,error) {
  return coordinatorRPCstub("Coordinator.ProposeRegister",ns,addr)
}

func coordinatorRPCstub(methodName string, ns * RemoteServer,addr string) (bool,error) {

  // connect to client // TODO: make a connection caching service that will create
  // a connection or recycle an old one.
  client, err := rpc.DialHTTP("tcp", addr)
  if err != nil {
    return false,err
  }
  // Synchronous call
  args := ns
  var reply bool
  err = client.Call(methodName, args, &reply)

  return reply,err

}

func AttachToGroupRPC(coordinator bool, laddr, lport string, mid uint64, addr string) (RegisterReply,[]ShardCopy,error) {

  rs := RemoteServer{laddr,lport,mid,coordinator}

  var res ConnectReply
  // start connection to remote Server
  client, err := rpc.DialHTTP("tcp", addr)
  if err != nil {
    return RegisterReply{},[]ShardCopy{},err
  }
  // make the rpc call
  err = client.Call("Coordinator.AttachRSToGroup", rs, &res)
  if (err != nil) {
    return RegisterReply{},[]ShardCopy{},err
  }
  return res.RegReply, res.Partitions, err
}
