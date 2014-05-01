package main
import (
  . "EDHT/common"
  "sync"
)


var PreCommitLock sync.Mutex

func preCommit(key string, value string) bool {
  ml.VPrintln("tpc","precommiting")
  PreCommitLock.Lock()
  defer PreCommitLock.Unlock()

  _, err := preCommits.Get(key)
  if err == nil {
    return false // there is a pending commit for this key DO NOT PRECOMMIT
  } else {
    ml.VPrintln("tpc","precommiting")
    preCommits.Put(key,value)
    return true
  }
}

func commit(key string) string{


  str,_ := preCommits.Get(key)
  t := Tuple{key,str}

  ov,replace,err := lookup(key)

  if (err!=nil || !replace) { // there wasn't a key here
    addkey(1)
    addbytes(len(t.Value))
  } else {
    addbytes(len(t.Value)-len(ov))
  }

  ml.VPrintln("data","ov:",ov)
  // update the stats on number of keys and data

  insert(t)
  preCommits.Delete(key)
  return ov
}


func abort(key string) {
  preCommits.Delete(key)
}

func (t* Daemon) PreCommit(pair Tuple, reply * bool) error{

  key   := pair.Key
  value := pair.Value


  *reply = preCommit(key,value)
  return nil
}

func (t * Daemon) Commit(key string, reply *map[string]string) (error){
  ov := commit(key)
  ret := make(map[string]string)
  ret["ov"] = ov
  *reply = ret
  return nil
}

func (t * Daemon) Abort(key string, reply *map[string]string) (error){
  *reply = nil
  abort(key)
  return nil
}
