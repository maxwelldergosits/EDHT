/*

This files is all rpc wrappers for each daemon method.

*/
package main

import (
  . "EDHT/common"
  "errors"
)

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

// arg = 1 if you want the number of keys
// arg = 2 if you want the number of bytes
func (t * Daemon) GetInfo(arg uint, reply * uint) error {
  switch arg {
    case 1:
      *reply = NKeys()
      return nil
    case 2:
      *reply = NBytes()
      return nil
  }
  return errors.New("Invalid Parameter")
}
