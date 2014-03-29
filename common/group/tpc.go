package group
import (
  "EDHT/utils"
 . "EDHT/common"
)
var pendingCommits      map[uint64]RemoteServer

// this method does the two phase commit for a new server
func AttachRSToGroup_local(rs RemoteServer) RegisterReply {
  verboseLog("attaching:",rs)
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

  var rpc = func(v RemoteServer)(bool){
    return (propseRegisterRPC(&rs,v.Address+":"+v.Port)==1)
  }

  var rc = func(v RemoteServer){
    registerRPC(&rs,v.Address+":"+v.Port)
  }

  var ra = func(v RemoteServer){
    RollBackRegisterRPC(&rs,v.Address+":"+v.Port)
  }
  var acceptors map[uint64]RemoteServer = make(map[uint64]RemoteServer)
  for k,v := range defaultGroup.Coordinators {
    acceptors[k] = v
  }
  t := utils.InitTPC(acceptors,id,lpc,lc,la,rpc,rc,ra)


  if (t.Run()){
    if rs.Coordinator {
      return RegisterReply{defaultGroup.Coordinators,defaultGroup.Daemons,rs.ID,defaultGroup.Nshards,defaultGroup.Nfailures}
    }else {
      return RegisterReply{nil,nil,rs.ID,0,0}
    }
  }
  return RegisterReply{nil,nil,0,0,0}
}



func preCommit(rs RemoteServer)bool{
  verboseLog("precommiting:",rs)
  pendingCommits[rs.ID] = rs
  return true
}

func localCommit(rs RemoteServer){
  verboseLog("commiting:",rs)
  if rs.Coordinator {
    defaultGroup.Coordinators[rs.ID]=rs
  } else {
    verboseLog("added new Daemon")
    newDaemonCallback(rs.ID)
    defaultGroup.Daemons[rs.ID]=rs
  }
  delete(pendingCommits,rs.ID)
}

func localAbort(rs RemoteServer) {
  verboseLog("aborting:",rs)
  delete(pendingCommits,rs.ID)
}
