package main

import(
  . "EDHT/common"
  "EDHT/common/rpc_stubs"
  "EDHT/utils"
  "EDHT/common/group"
  "errors"
)

type Shard struct {
  start uint64
  end uint64
  daemons map[uint64]bool
  keys uint
}

func PutKV(key string, value string) (bool,string){
  shard := partitions.GetShardForKey(key)
  nD := uint(len(shard.daemons))
  if (nD < group.GetNFailures() + 1) {
    return false,""
  }
  succ,info := tryTPC(shard,key,value)
  return succ,info["ov"] // returns the old value of the key "" if there was no key
}


func GetK(key string) (string,error) {
  shard := partitions.GetShardForKey(key)
  return getValue(shard,key)
}

func tryTPC(shard *Shard, key string, value string) (bool,map[string]string){

  noop := func() {}
  rc   := func(v RemoteServer)map[string]string {
    info, err := rpc_stubs.CommitDaemonRPC(key,v)
    if (err == nil) {
      return info
    } else {
      return nil
    }
  }
  ra  := func(v RemoteServer)map[string]string {
    info, err := rpc_stubs.AbortDaemonRPC(key,v)
    if (err != nil) {
      return info
    } else {
      return nil
    }
  }
  rpc := func(v RemoteServer)(bool,error) {
    succ,err := rpc_stubs.PreCommitDaemonRPC(key,value,v)
    return succ,err
  }

  id := group.GetLocalID()

  acceptors := make(map[uint64]RemoteServer)
  for k,_ := range shard.daemons {
    acceptors[k] = group.GetDaemon(k)
  }

  var failure = func(v RemoteServer) {
    group.DeleteDaemon(v.ID)
    delete(shard.daemons,v.ID)
  }

  tpc := utils.InitTPC(acceptors,id,noop,noop,noop,rpc,rc,ra,failure)
  succ,info := tpc.Run()
  if succ {
    ml.VPrintf("kv","Commited key %s\n",key)
  } else {
    ml.VPrintf("kv","Aborted key %s\n",key)
  }
  return succ,info
}

func getValue(shard * Shard, key string) (string,error) {
  time := utils.GetTimeNano()
  d := int(time % uint64(len(shard.daemons)))

  i := 0
  for k,_ := range shard.daemons {
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

