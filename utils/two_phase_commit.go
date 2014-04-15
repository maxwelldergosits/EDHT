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
  remoteCommit    func(rs RemoteServer)map[string]string
  remoteAbort     func(rs RemoteServer)map[string]string

  failure         func(rs RemoteServer)

}

func InitTPC(acceptors map[uint64]RemoteServer, id uint64,
              localPreCommit func(), localCommit func(), localAbort func(),
              remotePreCommit func(rs RemoteServer)(bool,error), remoteCommit func(rs RemoteServer)map[string]string,
              remoteAbort func(rs RemoteServer)map[string]string, failure func(rs RemoteServer)) (TwoPhaseCommit) {


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

func (t * TwoPhaseCommit) Run() (bool,map[string]string){

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


  var action func(rs RemoteServer)map[string]string
  if (doCommit) {
    action = t.remoteCommit
  } else {
    action = t.remoteAbort
  }

  var ret map[string]string

  for k,v := range t.acceptors {
    if (k == t.id) {continue}
    n++
    go func(v RemoteServer){
      ret = action(v)
      done <- false//value doesn't matter
    }(v)
  }
  if doCommit {
    t.localCommit()
  } else {
    t.localAbort()
  }
  for i:=0; i<n;i++ {
    <-done
  }
  return doCommit,ret
}
