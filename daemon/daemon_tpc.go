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

func commit(key string) {
  t := Tuple{key,pendingCommmits[key]}
  insert(t)
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

func (t * Daemon) Commit(key string, reply *bool) (error){
  commit(key)
  return nil
}

func (t * Daemon) Abort(key string, reply *bool) (error){
  abort(key)
  return nil
}
