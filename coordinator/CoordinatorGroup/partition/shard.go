package partition
import . "EDHT/common"
import "github.com/mad293/mlog"

type PartitionDelegate interface {
  GetDaemon(uint64) RemoteServer
  DeleteDaemon(uint64)
  GetLocalID() uint64
  Logger() *mlog.MLog
}

type Shard struct {
  Start uint64
  End uint64
  Daemons map[uint64]bool
  Keys uint
  delegate PartitionDelegate
}

type PartitionSet struct {
  Shards []Shard
  d PartitionDelegate
  tpcInProgress bool
  newRanges Ranges
  updateID uint64
}

type Diff struct {
  From int
  To   int
  Start uint64
  End uint64
}


