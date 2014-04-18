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
