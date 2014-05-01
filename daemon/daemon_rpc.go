/*


This files is all rpc wrappers for each daemon method.

*/
package main

import (
  . "EDHT/common"
  "EDHT/common/rpc_stubs"
  "errors"
)

/*Given a Tuple struct, Put will associate the member "key" with member "value" in the daemon's data store.
 *Performing multiple puts with the same key but different values will result in the key being
 *associated with the most recent value. The function's return value
 *is non-nil if the storage is successful. The empty string is not accepted as a valid key and
 *will result in Put failure. */
func (t *Daemon) Put(pair Tuple, reply *bool) error {
	if(pair.Value == "") {
		return errors.New("daemon Put error: empty key")
	}
	insert(pair)
  *reply = true
	return nil
}

/*Get will attempt to retrieve the value associated with the given key in a daemon's data store.
 *Key is the provided lookup key, while reply will contain a string representing
 *the result of the lookup operation. The return value of the function will indicated whether the lookup
 *succeeded; as stated in Go RPC Semantics, a non-nil error value (i.e. key is not in store) will not
 *return a result in the reply parameter.*/
func (t *Daemon) Get(key string, reply *string) error {
  ml.VPrintln("getting:",key)
	val, _, err := lookup(key)
	*reply = val
	return err
}


// arg = 1 if you want the number of keys
// arg = 2 if you want the number of bytes
func (t * Daemon) GetInfo(arg uint, reply * int) error {
  switch arg {
    case 1:
      *reply = int(NKeys())
      ml.VPrintln("ps","NKeys =",*reply)
      return nil
    case 2:
      *reply = int(NBytes())
      return nil
  }
  return errors.New("Invalid Parameter")
}

func (t * Daemon) RetrieveKeysInRange(srange ServerRange, keys* []string) error {


  rs := srange.Server
  ks := srange.Range

  newKVs,err  := rpc_stubs.GetKVsInRangeDaemonRPC(ks.Start,ks.End,rs)
  if (err != nil) {
    *keys = nil
    return err
  }
  newKeys := make([]string,0,len(newKVs))
  for k,v := range newKVs {
    insert(Tuple{k,v})
  }
  *keys = newKeys
  return nil
}

func (t * Daemon) DeleteKeysInRange(ra Range, succ * bool) error {


  iterateKeys(func(key,value string) {
    if ra.Start <= key && key <= ra.End {
      deleteKey(key)
    }
  })
  return nil

}

func (t * Daemon) CommitKeys(keys []string, reply *bool) error {

  for i := range keys {
    commit(keys[i])
  }
  *reply = true
  return nil

}



func (t * Daemon) AbortKeys(keys []string, reply *bool) error {

  for i := range keys {
    abort(keys[i])
  }
  *reply = true
  return nil

}

func (t * Daemon) DeleteKeys(keys[] string,reply *bool) error {
  for i := range keys {
    deleteKey(keys[i])
  }
  return nil
}
