package main
import (
  . "EDHT/common"
  "sync"
)

var pendingCommmits map[string]string

var PreCommitLock sync.Mutex

func init() {
  pendingCommmits = make(map[string]string)

}


func preCommit(key string, value string) bool {
  PreCommitLock.Lock()
  defer PreCommitLock.Unlock()
  if _, ok := pendingCommmits[key]; ok {
    return false // there is a pending commit for this key DO NOT PRECOMMIT
  } else {
    pendingCommmits[key] =value
    return true
  }
}

func commit(key string) string{


  t := Tuple{key,pendingCommmits[key]}

  ov,replace,err := lookup(key)

  if (err!=nil || !replace) { // there wasn't a key here
    addkey(1)
    addbytes(len(t.Value))
  } else {
    addbytes(len(t.Value)-len(ov))
  }

  ml.NPrintln("ov:",ov)
  // update the stats on number of keys and data

  insert(t)
  delete(pendingCommmits,key)
  return ov
}

func abort(key string) {
  delete(pendingCommmits,key)
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
