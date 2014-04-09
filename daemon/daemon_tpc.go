package main
import (
  . "EDHT/common"
)

var pendingCommmits map[string]string

func InitTPC() {
  pendingCommmits = make(map[string]string)

}

func preCommit(key string, value string) bool {
  if _, ok := pendingCommmits[key]; ok {
    return false // there is a pending commit for this key DO NOT PRECOMMIT
  } else {
    pendingCommmits[key] =value
    return true
  }
}

func commit(key string) string{


  t := Tuple{key,pendingCommmits[key]}

  nbs := 0
  ov, err := lookup(key)


// update the stats on number of keys and data
  if ov != "" || err != nil {
    addkey(1)
    nbs += len(key)
  } else {
    nbs -= (len(ov))
  }
  nbs += len(t.Value)
  addbytes(nbs)

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
  ml.NPrintln("ov:",ov)
  *reply = ret
  return nil
}

func (t * Daemon) Abort(key string, reply *map[string]string) (error){
  *reply = nil
  abort(key)
  return nil
}
