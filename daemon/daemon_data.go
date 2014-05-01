package main

import (
	"os"
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
    return val,true,nil
  }
	return "", false, errors.New("daemon lookup error: nonexistent key.")
}


func iterateKeys(iterFunc func(key,value string)) {

  for key,value := range data.Map() {
    iterFunc(key,value)
  }
}

func deleteKey(key string) {

  data.Delete(key)

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




/*WriteToDisk will attempt to dump the entire contents of the daemon's data store into a file
 *named by "path." As discussed on Friday, "path" will temporarily take the name of a local file
 *in the same directory, and will eventually be changed to provide support for absolute paths.
 *Error is nil on success, non-nil on failure; Reply is not used.
*/

func (t *Daemon) WriteToDisk(path string, reply *string) error{
	f, err := os.Create("path")
	if(err != nil){
		return errors.New("WriteToDisk: OS file error.")
	}

	defer f.Close()

	for key, value := range data.Map() {
		_, err := f.WriteString(key + ", " + value)
		if(err != nil){
			return errors.New("WriteToDisk: WriteString error.")
		}
	}
	f.Sync()
	return nil
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

