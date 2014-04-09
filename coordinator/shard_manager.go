package main
import (
  . "EDHT/common"
    "bytes"
    "encoding/binary"
)

var shards map[int]*Shard

// n must be a postive power of two, 2 4 8 16 32 etc
func MakeKeySpace(n int) map[int]*Shard {

  var shards_map map[int]*Shard = make(map[int]*Shard)
  var size uint64 = 1 << 63
  var chunk uint64 = size / uint64(n)

  for i:=0; i < n; i++ {
    ns := &Shard{Start:(uint64(i) * chunk),End:(uint64(i+1) * chunk) -1,Daemons:make(map[uint64]bool)}
    shards_map[i] = ns
  }

  shards = shards_map

  return shards_map
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

// function that gets called when a new daemon is commited to the system.
// Argument : id of the daemon
func NewDaemon(id uint64) {
  slot := int(djb2(id) % uint64(len(shards)))
  ml.NPrintln("daemon", id, "added to shard",slot)
  shards[slot].Daemons[id]= true
}


