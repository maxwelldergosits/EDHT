package main

import (
	"log"
	"strconv"
	"errors"
  . "EDHT/common"
)

var hashtable Hashtable= Hashtable{0,make(map[string]string)}

type Daemon int
//insert value in tuple into hashtable.
func insert(pair Tuple){
	hashtable.Store[pair.Key]= pair.Value
	hashtable.Size++

}

//return value corresponding to 'key'
func lookup(key string) string{
	return hashtable.Store[key]
}

//returns element number in hashtable (e.g first element store returns 1)
func (t *Daemon) Put(document *Tuple, reply *string) error {
	log.Printf("in Hashtable.Put")

	if (*document).Value == "" {
		log.Printf("value is empty string")
		return  errors.New("nothing to store")
	}

	insert(*document)
	log.Printf("\nstored to hashtable")
	log.Printf("\nTable contents: ", hashtable.Store)

  // what are we returning with here? -Maxwell
	*reply = strconv.Itoa(hashtable.Size-1)
	return nil
}


func (t *Daemon) Get(key *string, reply *string) error {
	log.Println("\nin Hashtable.Get\n")
	//for now, key is title of document
	*reply= lookup(*key)
	return nil

}
