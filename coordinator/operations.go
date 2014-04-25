package main

func GetKey(key string) (string,error) {
  pts := gc.GetPartitions()
  return pts.Get(key)
}

func PutKey(key string, value string, options map[string]bool) (error,map[string]string) {
  pts := gc.GetPartitions()
  err, info := pts.Put(key,value,options)
  return err,info
}

func GetInfo(i int) []uint {

  pts := gc.GetPartitions()

  // for shard get number of keys being held
  keys,err := pts.GetNKeysForEachShard()
  if err != nil {
    return []uint{}
  }

  return keys
}
