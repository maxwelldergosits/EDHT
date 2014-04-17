package CoordinatorGroup
import (
  "EDHT/coordinator/CoordinatorGroup/partition"
  . "EDHT/common"
)
type PD struct {
}

func (pd * PD) GetInfo(s *partition.Shard) {

}


func (pd * PD) UpdateShard(s *partition.Shard) {

}

func (pd * PD) CopyDiff(from,to *partition.Shard,a,b uint64) bool{
  return false
}
func (pd * PD) DeleteDiff(from*partition.Shard,a,b uint64) {

}

func (pd * PD) GetDaemon(uint64) RemoteServer {
  return RemoteServer{}
}

func (pd * PD) GetLocalID() uint64 {
return 0
}

func (pd * PD) DeleteDaemon(uint64) {

}
