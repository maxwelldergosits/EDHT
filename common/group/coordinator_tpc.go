package group
import (
  "EDHT/utils"
 . "EDHT/common"
  "EDHT/common/rpc_stubs"
)
var pendingCommits      map[uint64]RemoteServer

// this method does the two phase commit for a new server
func AttachRSToGroup_local(rs RemoteServer) RegisterReply {
  ml.VPrintln("debug","attaching:",rs)
  //update id
  machineid := rs.ID
  rs.ID = utils.GenId(machineid,rs.Coordinator)

  var lc = func(){
    localCommit(rs)
  }
  var la = func(){
    localAbort(rs)
  }

  var lpc = func(){
    preCommit(rs)
  }

  var rpc = func(v RemoteServer)(bool,error){
    return rpc_stubs.PropseRegisterRPC(&rs,v.Address+":"+v.Port)
  }

  var rc = func(v RemoteServer)map[string]string{
    rpc_stubs.RegisterRPC(&rs,v.Address+":"+v.Port)
    return nil
  }

  var ra = func(v RemoteServer)map[string]string {
    rpc_stubs.RollBackRegisterRPC(&rs,v.Address+":"+v.Port)
    return nil
  }
  var acceptors map[uint64]RemoteServer = make(map[uint64]RemoteServer)
  for k,v := range defaultGroup.Coordinators {
    acceptors[k] = v
  }

  var failure = func(rs RemoteServer) {
    DeleteCoordinator(rs.ID)
  }
  t := utils.InitTPC(acceptors,id,lpc,lc,la,rpc,rc,ra,failure)

  ok,_ := t.Run()
  if (ok){
    if rs.Coordinator {
      return RegisterReply{defaultGroup.Coordinators,defaultGroup.Daemons,rs.ID,defaultGroup.Nshards,defaultGroup.Nfailures}
    }else {
      return RegisterReply{nil,nil,rs.ID,0,0}
    }
  }
  return RegisterReply{nil,nil,0,0,0}
}



func preCommit(rs RemoteServer)bool{
  ml.VPrintln("debug","precommiting:",rs)
  pendingCommits[rs.ID] = rs
  return true
}

func localCommit(rs RemoteServer){
  ml.VPrintln("debug","commiting:",rs)
  if rs.Coordinator {
    defaultGroup.Coordinators[rs.ID]=rs
  } else {
    ml.VPrintln("debug","added new Daemon")
    newDaemonCallback(rs.ID)
    defaultGroup.Daemons[rs.ID]=rs
  }
  delete(pendingCommits,rs.ID)
}

func localAbort(rs RemoteServer) {
  ml.VPrintln("debug","aborting:",rs)
  delete(pendingCommits,rs.ID)
}
