package partition
import . "EDHT/common"

type PartitionDelegate interface {
  GetDaemon(uint64) RemoteServer
  DeleteDaemon(uint64)
  GetLocalID() uint64
}

type Shard struct {
  Start uint64
  End uint64
  daemons map[uint64]bool
  Keys uint
  delegate PartitionDelegate
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


