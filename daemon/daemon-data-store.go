package main

import (
	"log"
	"strconv"
	"errors"
  . "EDHT/common"
)

var hashtable Hashtable= Hashtable{0,make(map[string]string)}


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

/*Returns a pointer to a newly constructed Daemon object.
*/
func SpawnDaemon(address string, port string) *RemoteServer{
	remote_server := RemoteServer{Address:address, Port:port, ID:0, Coordinator:false}
	table := Hashtable{0, make(0, make(map[string]string))}
	return &Daemon{ServerDetails:remote_server, hashtable:table}
}

func RegisterDaemon(coordinator_ip string, coordinator_port string) error{

}

/*Given a DataPair struct, Put will associate the member "key" with member "value" in the daemon's data store.
 *Performing multiple puts with the same key but different values will result in the key being 
 *associated with the most recent value. The function's return value
 *is non-nil if the storage is successful. The empty string is not accepted as a valid key and
 *will result in Put failure. */
func (t *RemoteServer) Put(pair DataPair, reply *string) error {
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
func (t *RemoteServer) Get(key string, reply *string) error {
	val, err := lookup(key)
	*reply = val
	return err
}

/*GetAllKeys will return an array consisting of all keys in the daemon's data store in arbitrary order. 
 *Argument arg is ignored, and reply is not used. Error is nil on success, non-nil on failure.
 */
func (t *RemoteServer) GetAllKeys(arg string, reply *string) error{
	var keys[
}


func (t *RemoteServer) WriteToDisk(path string, reply *string) error{

}
