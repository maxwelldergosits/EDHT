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
