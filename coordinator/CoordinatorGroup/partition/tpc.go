package partition
import (
    "sync"
  . "EDHT/common"
    "net/rpc"
)

var tpcLock sync.Mutex

func (pts * PartitionSet) PreCommit(r Ranges, updateID uint64) bool {
  tpcLock.Lock()
  defer tpcLock.Unlock()

  if pts.tpcInProgress {
    return false
  } else {
    pts.newRanges = r
    pts.tpcInProgress = true
    pts.updateID = updateID
    return true
  }
}

type PreCommitRequest struct {
  R Ranges
  Id uint64
}
func (pts * PartitionSet) Commit(updateID uint64) bool {

  tpcLock.Lock()
  defer tpcLock.Unlock()

  if !pts.tpcInProgress || updateID != pts.updateID {
    return false
  } else {
    for i,v := range pts.newRanges.Rs{

      pts.Shards[i].Start = conv(v.Start)
      pts.Shards[i].End = conv(v.End)
    }
    pts.tpcInProgress = false
    return true
  }

}

func (pts * PartitionSet) Abort(updateID uint64) bool {

  tpcLock.Lock()
  defer tpcLock.Unlock()

  if !pts.tpcInProgress || updateID != pts.updateID {
    return false
  } else {
    pts.tpcInProgress = false
    return true
  }

}

func PreCommitPartition(r Ranges, id uint64, dest RemoteServer) (bool,error) {

  // connect to client // TODO: make a connection caching service that will create
  // a connection or recycle an old one.
  addr := dest.Address +":"+dest.Port
  client, err := rpc.DialHTTP("tcp", addr)
  if err != nil {
    return false,err
  }
  // Synchronous call
  args := PreCommitRequest{r,id}
  var reply bool
  err = client.Call("Coordinator.PreCommitPartition", args, &reply)

  return reply,err

}

func CommitPartition(id uint64, dest RemoteServer) (bool,error) {

  // connect to client // TODO: make a connection caching service that will create
  // a connection or recycle an old one.
  addr := dest.Address +":"+dest.Port
  client, err := rpc.DialHTTP("tcp", addr)
  if err != nil {
    return false,err
  }
  // Synchronous call
  args := id
  var reply bool
  err = client.Call("Coordinator.CommitPartition", args, &reply)

  return reply,err

}
func AbortPartition(id uint64, dest RemoteServer) (bool,error) {

  // connect to client // TODO: make a connection caching service that will create
  // a connection or recycle an old one.
  addr := dest.Address +":"+dest.Port
  client, err := rpc.DialHTTP("tcp", addr)
  if err != nil {
    return false,err
  }
  // Synchronous call
  args := id
  var reply bool
  err = client.Call("Coordinator.AbortPartition", args, &reply)

  return reply,err

}


