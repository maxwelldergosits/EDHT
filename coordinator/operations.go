package main

import "errors"

func GetKey(key string) (string,error) {
  pts := gc.GetPartitions()
  return pts.Get(key)
}

func PutKey(key string, value string, options map[string]bool) (error,map[string]string) {
  pts := gc.GetPartitions()
  if pts.CanCommit() {
    err, info := pts.Put(key,value,options)
    return err,info
  } else {
    ml.NPrintln("dumping commit")
    return errors.New("partition update in progress"),map[string]string{}
  }
}

func GetInfo(i int) []uint64 {

  pts := gc.GetPartitions()

  if (i == 1) {
    // for shard get number of keys being held
    keys,err := pts.GetNKeysForEachShard()

    if err != nil {
      return []uint64{}
    }
    out := make([]uint64,len(keys))
    for i := range keys {
      out[i] = uint64(keys[i])
    }
    return out
  } else if (i == 2) {
    tmp := pts.Ranges()
    out := make([]uint64,len(tmp)/2)
    for i:=0; i < len(tmp)/2; i++ {
      out[i] = tmp[(2*i) + 1] - tmp[(2*i)]
    }
    return out
  } else {
    return []uint64{}
  }
}

func getTopology() ([][]uint64, []uint64) {
  daemons := gc.GetPartitions().IDs()
  coordinators := gc.Gms.Coordinators()
  outC,i := make([]uint64, len(coordinators)),0
  for k, _ := range(coordinators){
    outC[i] = k
    i++
  }
  return daemons,outC
}
