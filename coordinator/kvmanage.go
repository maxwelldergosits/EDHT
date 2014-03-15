package main

import(
  "net/rpc"
  . "EDHT/common"
  "log"
)

func kv_get(key string, daemon RemoteServer) (string) {

  client, err :=rpc.DialHTTP("tcp",daemon.Address+":"+daemon.Port)
  if err != nil {
    log.Fatal("dialing error:",err)
  }

  var reply string
  err = client.Call("Daemon.Get",key,&reply)
  if err != nil {
    log.Fatal("calling error:",err)
  }
  return reply
}


func kv_put(key string, value string, daemon RemoteServer) (error) {

  client, err :=rpc.DialHTTP("tcp",daemon.Address+":"+daemon.Port)
  if err != nil {
    log.Fatal("dialing error:",err)
  }
  var args = Tuple{key,value}
  var reply string
  err = client.Call("Daemon.Put",args,&reply)
  if err != nil {
    log.Fatal("calling error:",err)
  }
  return err
}


