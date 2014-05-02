package main

import (
	"errors"
  . "EDHT/common"
  "sync"
)
type Daemon int
//var hashtable Hashtable= Hashtable{0,make(map[string]string)} deprecated; each Daemon has it's own store now.

var (
  nbytes int
  nkeys int
  keyMutex sync.Mutex
  byteMutex sync.Mutex
)

//insert value in tuple into hashtable.
func insert(pair Tuple){
  ml.VPrintln("kv","Adding key:",pair.Key)
	data.Put(pair.Key,pair.Value)
}

//return value corresponding to 'key'
func lookup(key string) (string, bool,error){
  val,err := data.Get(key)
  if err == nil {
    ml.VPrintf("got %s for %s",val,key)
    return val,true,nil
  }
	return "", false, errors.New("daemon lookup error: nonexistent key.")
}


func iterateKeys(iterFunc func(key string)) {

  for key := range data.KeyChan() {
    iterFunc(key)
  }
}

func deleteKey(key string) {

  err := data.Delete(key)
  if err == nil {
    addkey(-1)
  }

}

// function to keep track of the number of keys in the system
// arg = number of keys to add
func addkey(arg int) {
  keyMutex.Lock()
  nkeys += arg
  keyMutex.Unlock()
}

func NKeys() int {
  return nkeys
}

func NBytes() uint {
  return uint(nbytes)
}
// function to keep track of the number of bytes in the system
// arg = number of bytes to add (or subtract if it is negative)
func addbytes(arg int)  {
  byteMutex.Lock()
  nbytes += arg
  byteMutex.Unlock()
}


func (t * Daemon) GetAllKVsInRange(Keyrange Range, reply *map[string]string) error{

  newmap := make(map[string]string)

  for k,v := range data.Map() {
    if (Keyrange.Start < k && k < Keyrange.End) {
      newmap[k] = v
    }
  }

  *reply = newmap
  return nil
}

