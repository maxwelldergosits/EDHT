package partition


import (
  "EDHT/common/rpc_stubs"
  "errors"
  "EDHT/utils"
)

func (shard * Shard) GetInfoForShard() (error,uint){
  time := utils.GetTimeNano()
  if (len(*shard.Daemons())) <= 0 {return errors.New("No Daemons"),0}
  d := int(time % uint64(len(*shard.Daemons())))

  i := 0
  for k,_ := range *shard.Daemons() {
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
  if (len(*shard.Daemons())) <= 0 {return 0}
  d := int(time % uint64(len(*shard.Daemons())))
  i := 0
  for k,_ := range *shard.Daemons() {
    if i==d {
      return k
    } else {
      i++
    }
  }
  return 0
}

func (shard * Shard) getValue(key string) (string,error) {
  time := utils.GetTimeNano()
  d := int(time % uint64(len(*shard.Daemons())))

  i := 0
  for k,_ := range *shard.Daemons() {
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
 return shard.tryTPC(key,value,options)
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
  keys := make([]uint,len(pts.shards))
  for i := range pts.shards {
    err,info := pts.shards[i].GetInfoForShard()
    keys[i] = info
    if err != nil {
      return []uint{},err
    }
  }
  return keys,nil
}


func (pts * PartitionSet) ApplyDiffs(diffs []Diff) {

  copyDiffs := make([]Diff,0)
  deleteDiffs := make([]Diff,0)

  for i := range diffs {
    if diffs[i].To == -1 {
      deleteDiffs = append(deleteDiffs,diffs[i])
    } else {
      copyDiffs = append(copyDiffs,diffs[i])
    }
  }

  for i := range copyDiffs {
    diff := copyDiffs[i]
    pts.shards[diff.To].CopyKVs(pts.shards[diff.From],diff.Start,diff.End)
  }

  for i := range deleteDiffs {
    diff := deleteDiffs[i]
    pts.shards[diff.From].DeleteKs(diff.Start,diff.End)
  }

}

func (shard * Shard) CopyKVs(otherShard * Shard, start, end uint64) {

  for k,_ := range *shard.Daemons() {
    fromServer := shard.delegate.GetDaemon(otherShard.getDaemon())
    toServer   := shard.delegate.GetDaemon(k)
    rpc_stubs.RetrieveKeysInRangeDaemonRPC(unconv(start),unconv(end),fromServer,toServer)
  }

}

func (shard * Shard) DeleteKs(start, end uint64) {

  for k,_ := range *shard.Daemons() {
    fromServer := shard.delegate.GetDaemon(k)

    rpc_stubs.DeleteKeysInRangeDaemonRPC(unconv(start),unconv(end),fromServer)
  }

}
