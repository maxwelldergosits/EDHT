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


func getShardForKey(key string) *Shard{

  var n = conv(key)

  for _,shard := range shards {
    if shard.Start <= n && n <= shard.End {
      return shard
    }
  }
  return nil // This should never happen

}


func PutKV(key string, value string) bool{
  shard := getShardForKey(key)
  succ := tryTPC(shard,key,value)
  return succ
}


func getK(key string) (string,error) {
  shard := getShardForKey(key)
  return getValue(shard,key)
}

func tryTPC(shard *Shard, key string, value string) bool{

}

func getValue(shard * Shard, key string) (string,error) {
  return "",nil
}

