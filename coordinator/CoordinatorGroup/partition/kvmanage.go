package partition

import(
  . "EDHT/common"
  "EDHT/common/rpc_stubs"
  "EDHT/utils"
)

func (shard * Shard) tryTPC(key string, value string, options map[string]bool) (error,map[string]string){
  request := PutRequest{
    key,
    value,
    options}

  noop := func() {}
  rc   := func(v RemoteServer)map[string]string {
    info, err := rpc_stubs.CommitDaemonRPC(key,v)
    if (err == nil) {
      return info
    } else {
      return nil
    }
  }
  ra  := func(v RemoteServer)map[string]string {
    info, err := rpc_stubs.AbortDaemonRPC(key,v)
    if (err != nil) {
      return info
    } else {
      return nil
    }
  }
  rpc := func(v RemoteServer)(bool,error) {
    succ,err := rpc_stubs.PreCommitDaemonRPC(request,v)
    return succ,err
  }

  id := shard.delegate.GetLocalID()

  acceptors := make(map[uint64]RemoteServer)
  for k,_ := range *shard.Daemons() {
    acceptors[k] = shard.delegate.GetDaemon(k)
  }

  var failure = func(v RemoteServer) {
    shard.delegate.DeleteDaemon(v.ID)
    delete(*shard.Daemons(),v.ID)
  }

  tpc := utils.InitTPC(acceptors,id,noop,noop,noop,rpc,rc,ra,failure)
  err,info := tpc.Run()
  return err,info
}




