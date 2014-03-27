package group
import (
  "EDHT/utils"
 . "EDHT/common"
)
var pendingCommits      map[int64]RemoteServer

// this method does the two phase commit for a new server
func AttachRSToGroup_local(rs RemoteServer) RegisterReply {
  verboseLog("attaching:",rs)
  //update id
  machineid := rs.ID
  rs.ID = utils.GenId(machineid,rs.Coordinator)

  done := make(chan bool)
  var num = 0

  //precommit to all nodes
  for _,v := range defaultGroup.Coordinators {

    if v.ID == id {continue} // we shouldn't wait for a response from ourselves
    num++
    // propose precommit
    go func(v RemoteServer) {
      addr := v.Address+":"+v.Port
      done <- (propseRegisterRPC(&rs,addr) == 1)
    }(v)
  }

  commit := true
  //wait for respsonses
  for i := 0; i < num; i++ {
    commit = (commit && <-done)
  }

  if commit {
    localCommit(rs)
    for _,v := range defaultGroup.Coordinators {
      if v == rs || v.ID == id {num--; continue} // we shouldn't wait for a response from ourselves or the new server
      go func(v RemoteServer) {
        addr := v.Address+":"+v.Port
        var commit = (registerRPC(&rs,addr) == 1)
        done <- commit
      }(v)
    }
  } else {
    localAbort(rs)
    for _,v := range defaultGroup.Coordinators {
      if v == rs {continue}
      go func(v RemoteServer) {
        addr := v.Address+":"+v.Port
        var commit = (registerRPC(&rs,addr) == 1)
        done <- commit
      }(v)
    }
  }
  for i := 0; i < num; i++ {
    <-done
  }
  if rs.Coordinator {
    return RegisterReply{defaultGroup.Coordinators,defaultGroup.Daemons,rs.ID}
  }else {
    return RegisterReply{nil,nil,rs.ID}
  }
}



func preCommit(rs RemoteServer) int {

  verboseLog("precommiting:",rs)
  pendingCommits[rs.ID] = rs

  return 1
}

func localCommit(rs RemoteServer) int {

  verboseLog("commiting:",rs)

  if rs.Coordinator {
    defaultGroup.Coordinators[rs.ID]=rs
  } else {
    verboseLog("added new Daemon")
    defaultGroup.Daemons[rs.ID]=rs
  }
  delete(pendingCommits,rs.ID)

  return 1
}

func localAbort(rs RemoteServer) int {
  verboseLog("aborting:",rs)
  delete(pendingCommits,rs.ID)
  return 1
}
