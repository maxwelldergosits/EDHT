package partition
import (
    "bytes"
    "encoding/binary"
    . "EDHT/common"
)

func (t * PartitionSet) GetShardForKey(key string) *Shard{

  var n = conv(key)

  for _,shard := range t.shards {
    if shard.Start <= n && n <= shard.End {
      return shard
    }
  }
  return nil // This should never happen

}

func (t * PartitionSet) CanCommit() bool {
  return !t.tpcInProgress
}

func (t * Shard) Daemons() *map[uint64]bool{
  return &t.daemons
}


func MakePartitionSet(ns []ShardCopy, del PartitionDelegate) *PartitionSet{
  shards := make([]*Shard,len(ns))
  for i := range ns {
    sc := ns[i]
    shards[i] = &Shard{
      Start:sc.Start,
      End:sc.End,
      daemons:sc.Daemons,
      Keys:0,
      delegate:del}
  }
  return &PartitionSet{
    shards,
    del,false,nil,0}
}

func (pts * PartitionSet) GetShardCopies() []ShardCopy {
  scs := make([]ShardCopy,len(pts.shards))
  for i := range pts.shards {
    shard := pts.shards[i]
    scs[i] = ShardCopy {
      Start:shard.Start,
      End:shard.End,
      Daemons:shard.daemons}
  }
  return scs
}

// n must be a postive power of two, 2 4 8 16 32 etc
func MakeKeySpace(n int, del PartitionDelegate) *PartitionSet {

  var shards_map = make([]*Shard,n,n)
  var size uint64 = 1 << 63
  var chunk uint64 = size / uint64(n)

  for i:=0; i < n; i++ {
    ns := &Shard{Start:(uint64(i) * chunk),End:(uint64(i+1) * chunk) -1,daemons:make(map[uint64]bool),Keys:0,delegate:del}
    shards_map[i] = ns
  }

  return &PartitionSet{shards_map,del,false,nil,0}

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


func unconv(n uint64) string {

  arr := make([]byte,8)
   for i:= 0; i < 8; i++ {
     arr[i] = uint8(n>>uint(8*(7-i))&0xff)
  }
  return string(arr[:])

}

// function that gets called when a new daemon is commited to the system.
// Argument : id of the daemon
func ( t * PartitionSet) AddDaemon(id uint64) {
  slot := int(djb2(id) % uint64(len(t.shards)))
  t.d.Logger().VPrintln("gms","added daemon to shard",slot)
  t.shards[slot].daemons[id]= true
}

