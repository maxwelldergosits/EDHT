package partition


import (
  "EDHT/common/rpc_stubs"
  . "EDHT/common"
  "errors"
  "EDHT/utils"
)

func (shard * Shard) GetInfoForShard() (error,uint){
  time := utils.GetTimeNano()
  if (len(shard.Daemons)) <= 0 {return errors.New("No Daemons"),0}
  d := int(time % uint64(len(shard.Daemons)))

  i := 0
  for k,_ := range shard.Daemons {
    if i==d {
      rs := shard.delegate.GetDaemon(k)
      rep, err := rpc_stubs.GetInfoDaemonRPC(1,rs)
      if err != nil {
        return err,0
      } else {
        return nil,uint(rep)
      }
    } else {
      i++
    }
  }
  return errors.New("bad index"),0
}
// returns a daemon in the group
//preferably returning a different one everytime to load balance
func (shard * Shard) getDaemon() uint64 {
  time := utils.GetTimeNano()
  if (len(shard.Daemons)) <= 0 {return 0}
  d := int(time % uint64(len(shard.Daemons)))
  i := 0
  for k,_ := range shard.Daemons {
    if i==d {
      return k
    } else {
      i++
    }
  }
  return 0
}

func (shard * Shard) IDs() []uint64 {
  ids := make([]uint64, len(shard.Daemons))
  i := 0
  for k,_ := range shard.Daemons {
    ids[i] = k
    i++
  }
  return ids
}

func (pts * PartitionSet) IDs() [][]uint64 {
  ids := make([][]uint64,len(pts.Shards))
  for i := range pts.Shards {
    ids[i] = pts.Shards[i].IDs()
  }
  return ids
}

func (shard * Shard) getValue(key string) (string,error) {
  time := utils.GetTimeNano()
  d := int(time % uint64(len(shard.Daemons)))

  i := 0
  for k,_ := range shard.Daemons{
    if i==d {
      rs := shard.delegate.GetDaemon(k)
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

func (shard * Shard) Put(key,value string, options map[string]bool) (error,map[string]string) {
  if (options["unsafe"]) {
    for k,_ := range shard.Daemons {
      go func(){
        rs := shard.delegate.GetDaemon(k)
        rpc_stubs.DaemonPutRPC(key,value,rs)
      }()
    }
  return nil,map[string]string{"unsafe":"true"}
  } else {
    err,info:= shard.tryTPC(key,value,options)
    if err==nil {
      info["succ"] = "true"
    }
    return err,info
  }
}


func (pts * PartitionSet) Get(key string) (string,error) {
  shard := pts.GetShardForKey(key)
  return shard.getValue(key)
}

func (pts * PartitionSet) Put(key, value string, options map[string]bool) (error,map[string]string) {
  shard := pts.GetShardForKey(key)
  return shard.Put(key,value,options)
}

func (pts * PartitionSet) GetNKeysForEachShard() ([]uint,error) {
  keys := make([]uint,len(pts.Shards))
  for i := range pts.Shards {
    err,info := pts.Shards[i].GetInfoForShard()
    keys[i] = info
    if err != nil {
      return []uint{},err
    }
  }
  return keys,nil
}


func (pts * PartitionSet) ApplyCopyDiffs(diffs []Diff) bool{

  copyDiffs := make([]Diff,0)

  for i := range diffs {
    if diffs[i].From == -1 {
    } else {
      copyDiffs = append(copyDiffs,diffs[i])
    }
  }

  done := make(chan bool, len(copyDiffs))
  for i := range copyDiffs {
    diff := copyDiffs[i]

    go func(diff Diff) {
      done <- pts.Shards[diff.To].CopyKVs(pts.Shards[diff.From],diff.Start,diff.End)
    }(diff)

  }

  succ := true
  for i:= 0; i < len(copyDiffs); i++ {
    succ =(succ && <-done)
  }

  return succ
}


func (shard * Shard) CopyKVs(otherShard Shard, start, end uint64) bool {

  num_daemons := len(shard.Daemons)
  results := make(chan int,len(shard.Daemons))
  for k,_ := range shard.Daemons {
    fromServer := shard.delegate.GetDaemon(otherShard.getDaemon())
    toServer   := shard.delegate.GetDaemon(k)
    go func(fromServer,toServer RemoteServer) {
      keys, _ := rpc_stubs.RetrieveKeysInRangeDaemonRPC(unconv(start),unconv(end),fromServer,toServer)
      results <- len(keys)
    }(fromServer,toServer)
  }
  num := <- results
  succ := true
  for i := 1; i < num_daemons; i++ {
    res := <- results
    succ = succ && (num == res)
  }
  return succ
}


func (pts * PartitionSet) GarbageCollect() {

  for i := range pts.Shards {
    pts.Shards[i].GarbageCollect()
  }
}
func (shard * Shard) GarbageCollect() {

  for k,_ := range shard.Daemons {
    fromServer := shard.delegate.GetDaemon(k)

    rpc_stubs.DeleteKeysNotInRangeDaemonRPC(unconv(shard.Start),unconv(shard.End),fromServer)
  }

}
