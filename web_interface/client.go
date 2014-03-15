/*
Simple web interface for adding and retrieving key value pairs



Useage:

  There are two different web pages that are delivered.

  <coordinator-ip>:<port>/put

  This presents a web page for storing k-v pairs in the data base

  <coordinator-ip>:<port>/get

  This presents a web page for retrieving k-v pairs in the data base

  These pages are just simple HTML forms around a very simple restful API

  You can access these functions directly using


  <coordinator-ip>:<port>/put/submit?key=<key>&value=<value>


  <coordinator-ip>:<port>/get/submit?key=<key>

*/
package web_interface

import (
  "net/http"
  "log"
  "fmt"
)

func gethandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, getForm)
}
func getshandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf(r.FormValue("key"))
    fmt.Fprintf(w,"<put value here>")
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, putForm)
}
func shandler(w http.ResponseWriter, r *http.Request) {
    key:= r.FormValue("key")
    value:= r.FormValue("value")
    fmt.Fprintf(w,"Submitted: key:",key,"\n","value:",value)
}

func StartUp() {
  log.Println("starting")

  http.HandleFunc("/put",handler)
  http.HandleFunc("/put/submit",shandler)
  http.HandleFunc("/get",gethandler)
  http.HandleFunc("/get/submit",getshandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
