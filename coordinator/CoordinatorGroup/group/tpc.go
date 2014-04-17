package group

import (
  "EDHT/utils"
 . "EDHT/common"
  "EDHT/common/rpc_stubs"
)

// this method does the two phase commit for a new server
func (g * Group) AttachRSToGroup_local(rs RemoteServer) RegisterReply {
  g.ml.VPrintln("gms","attaching:",rs)
  //update id
  machineid := rs.ID
  rs.ID = utils.GenId(machineid,rs.Coordinator)

  var lc = func(){
    g.LocalCommit(rs)
  }
  var la = func(){
    g.LocalAbort(rs)
  }

  var lpc = func(){
    g.PreCommit(rs)
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
  for k,v := range g.coordinators {
    acceptors[k] = v
  }

  var failure = func(rs RemoteServer) {
    g.Delete(rs)
  }
  t := utils.InitTPC(acceptors,g.id,lpc,lc,la,rpc,rc,ra,failure)

  ok,_ := t.Run()
  if (ok!=nil){
    if rs.Coordinator {
      return RegisterReply{g.coordinators,g.daemons,rs.ID,g.nshards,g.nfailures}
    }else {
      return RegisterReply{nil,nil,rs.ID,0,0}
    }
  }
  return RegisterReply{nil,nil,0,0,0}
}



func (g * Group) PreCommit(rs RemoteServer)bool{
  g.ml.VPrintln("gms","precommiting:",rs)
  g.pendingCommits[rs.ID] = rs
  return true
}

func (g * Group) LocalCommit(rs RemoteServer){
  g.ml.VPrintln("gms","commiting:",rs)
  if rs.Coordinator {
    g.coordinators[rs.ID]=rs
  } else {
    g.ml.VPrintln("gms","added new Daemon")
    g.newDaemonCallBack(rs.ID)
    g.daemons[rs.ID]=rs
  }
  delete(g.pendingCommits,rs.ID)
}

func (g * Group) LocalAbort(rs RemoteServer) {
  g.ml.VPrintln("gms","aborting:",rs)
  delete(g.pendingCommits,rs.ID)
}
