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

func CommitDaemonRPC(key string, rs RemoteServer)(map[string]string,error) {

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return nil,errors.New("dialing error:"+err.Error())
  }

  var reply map[string]string
  err = client.Call("Daemon.Commit",key,&reply)
  if err != nil{
    return nil,err
  } else {
    return reply,nil
  }

}


func AbortDaemonRPC(key string,rs RemoteServer) (map[string]string,error){

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return nil,errors.New("dialing error:"+err.Error())
  }

  var reply map[string]string
  err = client.Call("Daemon.Abort",key,&reply)
  if err != nil{
    return nil,err
  } else {
    return reply,nil
  }

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

func GetInfoDaemonRPC(arg uint, rs RemoteServer) (int,error) {

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return -1,errors.New("dialing error:"+err.Error())
  }

  var reply int
  err = client.Call("Daemon.GetInfo",arg,&reply)
  return reply,err

}

