package main

import (
	"os"
	"log"
	"strconv"
	"errors"
  . "EDHT/common"
)

//var hashtable Hashtable= Hashtable{0,make(map[string]string)} deprecated; each Daemon has it's own store now.


//insert value in tuple into hashtable.
func insert(pair DataPair){
	hashtable.Store[pair.Key]= pair.Value
	hashtable.Size++
}

//return value corresponding to 'key'
func lookup(key string) (string, error){
	if value, in_store := dict[key]; in_store{
    		return hashtable.Store[key], nil
	}
	return "", errors.New("daemon lookup error: nonexistent key.")
}

/*Returns a pointer to a newly constructed Daemon object, with the specified address and port number.
*/
func SpawnDaemon(address string, port string) *RemoteServer{
	remote_server := RemoteServer{Address:address, Port:port, ID:0, Coordinator:false}
	table := Hashtable{0, make(0, make(map[string]string))}
	return &Daemon{ServerDetails:remote_server, hashtable:table}
}

/*Given the ip address and port for a known coordinator, will attempt to connect
 *to the coordinator with the provided ip address and port number.
 *TODO: test error condition with AttachToGroup
*/
func (t *Daemon) RegisterDaemon(coordinator_ip string, coordinator_port string) error{
	return AttachToGroup(coordinator_ip, coordinator_port)

}

/*Given a DataPair struct, Put will associate the member "key" with member "value" in the daemon's data store.
 *Performing multiple puts with the same key but different values will result in the key being
 *associated with the most recent value. The function's return value
 *is non-nil if the storage is successful. The empty string is not accepted as a valid key and
 *will result in Put failure. */
func (t *Daemon) Put(pair DataPair, reply *string) error {
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
	keys := make([]string, len(t.Store))
	i := 0
	for key, _ := range t.Store{
		keys[i] := key
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

	for key, value := range t.Store{
		_, err := f.WriteString(key + ", " + value)
		if(err != nil){
			return errors.New("WriteToDisk: WriteString error.")
		}
	}
	f.Sync()
	return nil
}
