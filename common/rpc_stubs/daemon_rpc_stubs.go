package rpc_stubs
import (
  . "EDHT/common"
  "EDHT/utils"
  "errors"
)

func PreCommitDaemonRPC(key string, value string, rs RemoteServer) (bool,error) {
  client, err := utils.MakeConnection(rs)
  if err != nil {
    return false,errors.New("dialing error:"+err.Error())
  }

  t:= Tuple{key,value}
  var reply bool
  err = client.Call("Daemon.PreCommit",t,&reply)
  return reply,err



}

func CommitDaemonRPC(key string, rs RemoteServer) error {

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return errors.New("dialing error:"+err.Error())
  }

  var reply bool
  err = client.Call("Daemon.Commit",key,&reply)
  return err

}


func AbortDaemonRPC(key string,rs RemoteServer) error{

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return errors.New("dialing error:"+err.Error())
  }

  var reply bool
  err = client.Call("Daemon.Abort",key,&reply)
  return err

}

func GetKeyDaemonRPC(key string, rs RemoteServer) (string,error) {

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return "",errors.New("dialing error:"+err.Error())
  }

  var reply string
  err = client.Call("Daemon.Get",key,&reply)
  return reply,nil

}
