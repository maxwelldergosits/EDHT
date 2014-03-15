package main

import (
	"os"
  "EDHT/utils"
	"errors"
  . "EDHT/common"
)
type Daemon int
//var hashtable Hashtable= Hashtable{0,make(map[string]string)} deprecated; each Daemon has it's own store now.


//insert value in tuple into hashtable.
func insert(pair Tuple){
	local_state.Hashtable.Store[pair.Key]= pair.Value
	local_state.Hashtable.Size++
}

//return value corresponding to 'key'
func lookup(key string) (string, error){
  if val, ok := local_state.Hashtable.Store[key];ok {
    return val,nil
  }
	return "", errors.New("daemon lookup error: nonexistent key.")
}

/*Returns a pointer to a newly constructed Daemon object, with the specified address and port number.
*/
func SpawnDaemon(address string, port string) DaemonData{
	remote_server := RemoteServer{Address:address, Port:port, ID:utils.GenMachineId(), Coordinator:false}
	table := Hashtable{0, make(map[string]string)}
	d := DaemonData{remote_server, table}
  return d
}


/*Given a Tuple struct, Put will associate the member "key" with member "value" in the daemon's data store.
 *Performing multiple puts with the same key but different values will result in the key being
 *associated with the most recent value. The function's return value
 *is non-nil if the storage is successful. The empty string is not accepted as a valid key and
 *will result in Put failure. */
func (t *Daemon) Put(pair Tuple, reply *string) error {
	if(pair.Value == "") {
		return errors.New("daemon Put error: empty key")
	}
	insert(pair)
	return nil
}


/*Get will attempt to retrieve the value associated with the given key in a daemon's data store.
 *Key is the provided lookup key, while reply will contain a string representing
 *the result of the lookup operation. The return value of the function will indicated whether the lookup
 *succeeded; as stated in Go RPC Semantics, a non-nil error value (i.e. key is not in store) will not
 *return a result in the reply parameter.*/
func (t *Daemon) Get(key string, reply *string) error {
	val, err := lookup(key)
	*reply = val
	return err
}

/*GetAllKeys will return an array consisting of all keys in the daemon's data store in arbitrary order.
 *Argument arg is ignored, and reply is not used. Error is nil on success, non-nil on failure.
 */
func (t *Daemon) GetAllKeys(arg string, reply *[]string) error{
	keys := make([]string, local_state.Hashtable.Size)
	i := 0
	for key, _ := range local_state.Hashtable.Store{
		keys[i] = key
	}
	*reply = keys
	return nil
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

	for key, value := range local_state.Hashtable.Store{
		_, err := f.WriteString(key + ", " + value)
		if(err != nil){
			return errors.New("WriteToDisk: WriteString error.")
		}
	}
	f.Sync()
	return nil
}
