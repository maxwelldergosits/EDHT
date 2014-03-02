/**
* TCPArithClient
 */

package main

import (
	"net/rpc"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Tuple struct {
        Key, Value string
}

type Hashtable struct {
	//can add metdata? (Karl)
	size int
	store map[string]string
}

func main() {
	c:= make(chan int,3)
	numClients:= 3
	for clientID:=1; clientID <= numClients; clientID++ {
		//create client threads
		go func(){
			if len(os.Args) != 2 {
				fmt.Println("Usage: ", os.Args[0], "server:port")
				os.Exit(1)
			}
			service := os.Args[1]

			client, err := rpc.Dial("tcp", service)
			if err != nil {
				log.Fatal("dialing:", err)
			}
	// Synchronous call
			document := Tuple{"test" + strconv.Itoa(clientID),"testing..."}
			fmt.Printf("document to store: ",document)
			reply:= ""
			err= client.Call("Hashtable.Put", document, &reply)
			fmt.Printf("\nreturned from Call to server\n")
			if err != nil {
				//log.Fatal("hashtable error:",err)
			}
			fmt.Printf("\nyour document is in position: ", reply)
	
		//retrieve document just stored.
			err= client.Call("Hashtable.Get", "test"+ strconv.Itoa(clientID), &reply)
			if err != nil {
				log.Fatal("hashtable error: ", err)
			}
			fmt.Printf("\nyour document is: ", reply)
			c<-clientID
			client.Close()
		}()
		//for some reason doesn't work if I don't use channels?(Karl)
		chanVal:= <-c
		fmt.Printf("\n Client %d is done.\n",chanVal)
	}
}
