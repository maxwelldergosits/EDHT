package main
import (
  . "EDHT/common"
    "bytes"
    "encoding/binary"
)
// n must be a postive power of two, 2 4 8 16 32 etc
func MakeKeySpace(n int) map[int]*Shard {

  var shards_map map[int]*Shard = make(map[int]*Shard)
  var size uint64 = 1 << 63
  var chunk uint64 = size / uint64(n)

  for i:=0; i < n; i++ {
    ns := &Shard{Start:(uint64(i) * chunk),End:(uint64(i+1) * chunk) -1,Daemons:make(map[int64]RemoteServer)}
    shards_map[i] = ns
  }
  return shards_map
}



func conv(key string) uint64 {
  var n uint64
  b := []byte(key)

  if len(b) < 8 {
    l := 8 - len(b)
    for i := 0; i < l; i++ {
    b = append(b,'\x00')
    }
  }
  buf := bytes.NewBuffer(b)
  binary.Read(buf, binary.BigEndian, &n)
  return n
}

func getShardForKey(shards map[int]*Shard, key string) *Shard{

  var n = conv(key)

  for _,shard := range shards {
    if shard.Start <= n && n <= shard.End {
      return shard
    }
  }
  return nil // This should never happen

}


func PutKV(shards map[int]*Shard, key string, value string) bool{

  shard := getShardForKey(shards,key)

  succ := tryTPC(shard,key,value)
  return succ
}


func getK(key string) {


}
