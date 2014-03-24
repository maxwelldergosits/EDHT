package main
import (
  . "EDHT/common"
)
// n must be a postive power of two, 2 4 8 16 32 etc
func MakeKeySpace(n int) map[int]Shard {

  var shards_map map[int]Shard = make(map[int]Shard)
  var size uint64 = 1 << 63
  var chunk uint64 = size / uint64(n)

  for i:=0; i < n; i++ {
    ns := Shard{Start:(uint64(i) * chunk),End:(uint64(i+1) * chunk) -1,Daemons:make(map[int64]RemoteServer)}
    shards_map[i] = ns
  }
  return shards_map
}


