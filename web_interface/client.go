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
  "mlog"
)

var (
  getF func(key string) (string,error)
  putF func(key string,value string) (bool,string)
  ml mlog.MLog
)

func gethandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, getForm)
}
func getshandler(w http.ResponseWriter, r *http.Request) {

    key:=r.FormValue("key")
    value,err := getF(key)
    if err != nil {
      value = "error"
    }
    fmt.Fprintf(w,"key: %s \nvalue: %s\n",key,value)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, putForm)
}
func shandler(w http.ResponseWriter, r *http.Request) {
    key:= r.FormValue("key")
    value:= r.FormValue("value")
    succ, ov := putF(key,value)
    if succ {
      fmt.Fprintf(w,"Submitted key:%s value:%s oldvalue:%s",key,value,ov)
      w.Header().Set("suc","true")
      w.Header().Set("key",key)
      w.Header().Set("value",value)
      w.Header().Set("value",ov)
    } else {
      fmt.Fprintf(w,"Error: was not able to submit key:%s\n",key)
      w.Header().Set("suc","false")
    }
}

func StartUp(logger mlog.MLog,port string, get func(key string)(string,error), put func(key string, value string) (bool,string)) {
  ml = logger
  getF = get
  putF = put
  ml.VPrintln("web","starting web inteface")

  http.HandleFunc("/put",handler)
  http.HandleFunc("/put/submit",shandler)
  http.HandleFunc("/get",gethandler)
  http.HandleFunc("/get/submit",getshandler)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
