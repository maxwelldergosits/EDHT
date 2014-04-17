package partition


import (
  "EDHT/common/rpc_stubs"
  "errors"
  "EDHT/utils"
)

func (shard * Shard) GetInfoForShard() error{
    sum := 0
    for id := range *shard.Daemons() {
      rs := shard.delegate.GetDaemon(id)
      keys, err := rpc_stubs.GetInfoDaemonRPC(1,rs)
      if (err != nil) {
        return err
      }
      sum += keys
    }
    avg := sum
    shard.Keys = uint(avg)
    return nil
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
