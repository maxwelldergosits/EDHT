package main

import(
  . "EDHT/common"
  "EDHT/common/rpc_stubs"
  "EDHT/utils"
  "EDHT/common/group"
  "errors"
)

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


func GetK(key string) (string,error) {
  shard := getShardForKey(key)
  return getValue(shard,key)
}

func tryTPC(shard *Shard, key string, value string) bool{

  noop := func() {}
  rc   := func(v RemoteServer) {
    rpc_stubs.CommitDaemonRPC(key,v)
  }
  ra  := func(v RemoteServer) {
    rpc_stubs.AbortDaemonRPC(key,v)
  }
  rpc := func(v RemoteServer)(bool) {
    succ,err := rpc_stubs.PreCommitDaemonRPC(key,value,v)
    return (succ || err!=nil)
  }

  id := group.GetLocalID()

  acceptors := make(map[uint64]RemoteServer)
  for k,_ := range shard.Daemons {
    acceptors[k] = group.GetDaemon(k)
  }
  tpc := utils.InitTPC(acceptors,id,noop,noop,noop,rpc,rc,ra)
  return tpc.Run()
}

func getValue(shard * Shard, key string) (string,error) {
  time := utils.GetTimeNano()
  d := int(time % uint64(len(shard.Daemons)))

  i := 0
  for k,_ := range shard.Daemons {
    if i==d {
      rs := group.GetDaemon(k)
      rep, err := rpc_stubs.GetKeyDaemonRPC(key,rs)
      if err != nil {
        return "",err
      } else {
        return rep,nil
      }
    } else {
      i++
    }
  }
  return "",errors.New("No key found")
}

