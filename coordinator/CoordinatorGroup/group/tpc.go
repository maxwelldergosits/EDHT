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
    return rpc_stubs.PropseRegisterRPC(&rs,v)
  }

  var rc = func(v RemoteServer)map[string]string{
    rpc_stubs.RegisterRPC(&rs,v)
    return nil
  }

  var ra = func(v RemoteServer)map[string]string {
    rpc_stubs.RollBackRegisterRPC(&rs,v)
    return nil
  }
  var acceptors map[uint64]RemoteServer = make(map[uint64]RemoteServer)
  for k,v := range g.coordinators.Map() {
    acceptors[k] = v
  }

  var failure = func(rs RemoteServer, e error) {
    g.ml.NPrintln(e.Error())
    g.Delete(rs)
  }
  t := utils.InitTPC(acceptors,g.id,lpc,lc,la,rpc,rc,ra,failure,true)

  ok,_ := t.Run()
  if (ok==nil){
    if rs.Coordinator {
      return RegisterReply{g.coordinators.Map(),g.daemons.Map(),rs.ID,g.nshards,g.nfailures}
    }else {
      return RegisterReply{nil,nil,rs.ID,0,0}
    }
  }
  g.ml.VPrintln("gms","didn't add server",ok.Error())
  return RegisterReply{}
}



func (g * Group) PreCommit(rs RemoteServer)bool{
  g.ml.VPrintln("gms","precommiting:",rs)
  g.pendingCommits.Put(rs.ID,rs)
  return true
}

func (g * Group) LocalCommit(rs RemoteServer){
  g.ml.VPrintln("gms","commiting:",rs)
  if rs.Coordinator {
    g.coordinators.Put(rs.ID,rs)
  } else {
    g.newDaemonCallBack(rs.ID)
    g.daemons.Put(rs.ID,rs)
  }
  g.pendingCommits.Delete(rs.ID)
}

func (g * Group) LocalAbort(rs RemoteServer) {
  g.ml.VPrintln("gms","aborting:",rs)
  g.pendingCommits.Delete(rs.ID)
}
