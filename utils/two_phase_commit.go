/*

This file presents a generic two phase commit type with associated method calls

this allows for less code duplicaton and inter module dependencies


*/
package utils

import . "EDHT/common"


type TwoPhaseCommit struct {
  //private fields

  acceptors map[uint64]RemoteServer
  id uint64

  localPreCommit func()
  localCommit func()
  localAbort func()

  remotePreCommit func(rs RemoteServer) (bool,error)
  remoteCommit    func(rs RemoteServer)
  remoteAbort     func(rs RemoteServer)

  failure         func(rs RemoteServer)

}

func InitTPC(acceptors map[uint64]RemoteServer, id uint64,
              localPreCommit func(), localCommit func(), localAbort func(),
              remotePreCommit func(rs RemoteServer)(bool,error), remoteCommit func(rs RemoteServer),
              remoteAbort func(rs RemoteServer), failure func(rs RemoteServer)) (TwoPhaseCommit) {


  return TwoPhaseCommit{

      acceptors,
      id,
      localPreCommit,
      localCommit,
      localAbort,
      remotePreCommit,
      remoteCommit,
      remoteAbort,
      failure}


}

func (t * TwoPhaseCommit) Run() (bool){

  t.localPreCommit()

  var n int
  done := make(chan(bool))

  // send the preCommit to each remote server
  for k,v := range t.acceptors {
    if k == t.id {continue}
    n++
    go func(v RemoteServer) {
      succ, err := t.remotePreCommit(v)
      if err != nil {
        t.failure(v)
        delete(t.acceptors,v.ID)
        done <- false
      } else {
        done <- succ
      }
    }(v)
  }

  doCommit := true
  //wait for every remote server to respond
  for i:=0; i<n;i++ {
    doCommit = (doCommit&&<-done)
  }
  n = 0 // reset n for the next round

  var action func(rs RemoteServer)
  if (doCommit) {
    t.localCommit()
    action = t.remoteCommit
  } else {
    t.localAbort()
    action = t.remoteAbort
  }

  for k,v := range t.acceptors {
    if (k == t.id) {continue}
    n++
    go func(v RemoteServer){
      action(v)
      done <- false//value doesn't matter
    }(v)
  }

  for i:=0; i<n;i++ {
    <-done
  }
  return doCommit
}
