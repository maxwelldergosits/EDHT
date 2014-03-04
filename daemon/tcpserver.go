/**
* TCPArithServer
 */

package main

import (
	"fmt"
	"strconv"
	"net/rpc"
	"errors"
	"net"
	"os"
)

type Tuple struct {
	Key, Value string
}

type Hashtable struct {
	//can add metadata? (Karl)
	size int
	store map[string]string
}

var hashtable Hashtable= Hashtable{0,make(map[string]string)}

func (t *Tuple) GetKey(tuple *Tuple, reply *string) error{
	*reply= tuple.Key
	return nil
}

func (t *Tuple) GetValue(tuple *Tuple, reply *string) error{
	*reply= tuple.Value
	return nil
}

//insert value in tuple into hashtable.
func insert(pair Tuple){
	hashtable.store[pair.Key]= pair.Value
	hashtable.size= hashtable.size+1 

}

//return value corresponding to 'key'
func lookup(key string) string{
	return hashtable.store[key]
}

//returns element number in hashtable (e.g first element store returns 1)
func (t *Hashtable) Put(document *Tuple, reply *string) error {
	fmt.Printf("in Hashtable.Put")
	if (*document).Value == "" {
		fmt.Printf("value is empty string")
		return  errors.New("nothing to store")
	}
	//hashtable.store[0]= *value
	insert(*document)
	fmt.Printf("\nstored to hashtable")
	fmt.Printf("\nTable contents: ", hashtable.store)
	*reply = strconv.Itoa(hashtable.size-1)
	return nil
}

func (t *Hashtable) Get(key *string, reply *string) error {
	fmt.Println("\nin Hashtable.Get\n")
	//for now, key is title of document
	*reply= lookup(*key)
	return nil
}

func main() {

	tuple:= new(Tuple)
	hashtable := new(Hashtable)
	rpc.Register(tuple)
	rpc.Register(hashtable)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)
	fmt.Println(tcpAddr)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	/* This works:
	rpc.Accept(listener)
	*/
	/* and so does this:
	 */
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
	fmt.Println("END OF MAIN")

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
