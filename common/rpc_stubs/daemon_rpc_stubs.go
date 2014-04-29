package rpc_stubs
import (
  . "EDHT/common"
  "EDHT/utils"
  "errors"
)

func DaemonPutRPC(key, value string,rs RemoteServer) (bool,error) {
  client,err := utils.MakeConnection(rs)
  if err != nil {
    return false,errors.New("dialing error:"+err.Error())
  }

  tuple := Tuple{key,value}

  var reply bool
  err = client.Call("Daemon.Put",tuple,&reply)
  return reply,err

}

func PreCommitDaemonRPC(req PutRequest, rs RemoteServer) (bool,error) {
  client, err := utils.MakeConnection(rs)
  if err != nil {
    return false,errors.New("dialing error:"+err.Error())
  }

  var reply bool
  err = client.Call("Daemon.PreCommit",req,&reply)
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


func GetKVsInRangeDaemonRPC(start, end string, rs RemoteServer) (map[string]string,error) {

  keys := Range{start,end}

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return nil,errors.New("dialing error:"+err.Error())
  }

  var reply map[string]string
  err = client.Call("Daemon.GetAllKVsInRange",keys,&reply)
  return reply,err

}

func RetrieveKeysInRangeDaemonRPC(start, end string, od RemoteServer, rs RemoteServer) ([]string,error) {


  ra := Range{start,end}
  arg := ServerRange{od,ra}

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return nil,errors.New("dialing error:"+err.Error())
  }

  var reply []string
  err = client.Call("Daemon.RetrieveKeysInRange",arg,&reply)
  return reply,err


}
func DeleteKeysInRangeDaemonRPC(start, end string, rs RemoteServer) (error) {


  ra := Range{start,end}

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return errors.New("dialing error:"+err.Error())
  }

  var reply bool
  err = client.Call("Daemon.DeleteKeysInRange",ra,&reply)
  return err


}

func CommitKeysDaemonRPC(keys []string, rs RemoteServer) error {

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return errors.New("dialing error:"+err.Error())
  }

  var reply bool
  err = client.Call("Daemon.CommitKeys",keys,reply)
  return err

}

func AbortKeysDaemonRPC(keys []string, rs RemoteServer) error {

  client, err := utils.MakeConnection(rs)
  if err != nil {
    return errors.New("dialing error:"+err.Error())
  }

  var reply bool
  err = client.Call("Daemon.AbortKeys",keys,reply)
  return err

}
