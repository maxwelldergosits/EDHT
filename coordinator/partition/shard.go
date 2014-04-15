package partition

type PartitionDelegate interface {
  UpdateShard(*Shard)
  GetInfo(*Shard)
  CopyDiff(*Shard, *Shard,uint64,uint64) bool
  DeleteDiff(*Shard,uint64,uint64)
}

type Shard struct {
  Start uint64
  End uint64
  daemons map[uint64]bool
  Keys uint
}

type PartitionSet struct {
  shards []*Shard
  d PartitionDelegate
}

type Diff struct {
  From int
  To   int
  Start uint64
  End uint64

}
