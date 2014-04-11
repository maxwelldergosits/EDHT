package main
import (
    "bytes"
    "encoding/binary"
    "EDHT/common/rpc_stubs"
    "EDHT/common/group"
)

type PartitionSet struct {
  shards map[int]*Shard
}


func (t * PartitionSet) GetShardForKey(key string) *Shard{

  var n = conv(key)

  for _,shard := range t.shards {
    if shard.start <= n && n <= shard.end {
      return shard
    }
  }
  return nil // This should never happen

}

// n must be a postive power of two, 2 4 8 16 32 etc
func MakeKeySpace(n int) *PartitionSet {

  var shards_map map[int]*Shard = make(map[int]*Shard)
  var size uint64 = 1 << 63
  var chunk uint64 = size / uint64(n)

  for i:=0; i < n; i++ {
    ns := &Shard{start:(uint64(i) * chunk),end:(uint64(i+1) * chunk) -1,daemons:make(map[uint64]bool)}
    shards_map[i] = ns
  }

  return &PartitionSet{shards_map}

}

// general purpose hash function for distributing nodes into the shards
func djb2(n uint64) uint64{

  var hash uint64= 5381

  h := func(c uint64,hash uint64) (uint64){
    return (((hash << 5) + hash) + c)
  }

  hash = h((n >> (7*8) & 0xFF),hash)
  hash = h((n >> (6*8) & 0xFF),hash)
  hash = h((n >> (5*8) & 0xFF),hash)
  hash = h((n >> (4*8) & 0xFF),hash)
  hash = h((n >> (3*8) & 0xFF),hash)
  hash = h((n >> (2*8) & 0xFF),hash)
  hash = h((n >> (1*8) & 0xFF),hash)
  hash = h((n >> (0*8) & 0xFF),hash)

  return hash;
}

// return the integer representation of a string
func conv(key string) uint64 {

  var n uint64
  b := []byte(key)

  if len(b) < 8 { // make sure there are at least 8 bytes
    l := 8 - len(b)
    for i := 0; i < l; i++ {
      b = append(b,'\x00')
    }
  }

  buf := bytes.NewBuffer(b)
  binary.Read(buf, binary.BigEndian, &n)

  return n
}

// function that gets called when a new daemon is commited to the system.
// Argument : id of the daemon
func ( t * PartitionSet) AddDaemon(id uint64) {
  slot := int(djb2(id) % uint64(len(t.shards)))
  ml.VPrintln("gms","daemon", id, "added to shard",slot)
  t.shards[slot].daemons[id]= true
}



func (t * PartitionSet) GatherInfo() {

  sum := 0
  for _,shard := range t.shards {
    for id := range shard.daemons {
      rs := group.GetDaemon(id)
      ml.VPrintln("ps","rs = ",rs)
      keys, err := rpc_stubs.GetInfoDaemonRPC(1,rs)
      if (err != nil) {
        ml.VPrintln("ps","error:",err.Error())
      }
      sum += keys
    }
    avg := sum
    shard.keys = uint(avg)
    ml.VPrintln("ps",shard.keys)
  }

}

