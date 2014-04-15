package main
import (
  "EDHT/coordinator/partition"
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
