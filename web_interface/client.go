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
  "encoding/json"
  "strconv"
  "github.com/mad293/mlog"
)

var (
  getF func(key string) (string,error)
  putF func(string,string,map[string]bool) (error,map[string]string)
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
    group := make(map[string]string)
    group["key"] =key
    group["value"] = value
    b, _ := json.MarshalIndent(group,"","  ")
    w.Write(b)
    w.Header().Set("Content-Type","text/json")
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, putForm)
}
func shandler(w http.ResponseWriter, r *http.Request) {
    key:= r.FormValue("key")
    value:= r.FormValue("value")
    getOV:= r.FormValue("ov")
    getOVbool,_ := strconv.ParseBool(getOV)
    option:= map[string]bool{"ov":getOVbool}
    err,values := putF(key,value,option)
    if (err == nil) {

      succ, _:= strconv.ParseBool(values["succ"])
      group := make(map[string]string)
      if succ {
        if (getOVbool) {
          ov := values["ov"]
          group["ov"] = ov
        }
        group["key"] =key
        group["value"] = value
        group["succ"] = "true"
      } else {
        group["succ"] = "false"
      }
      b, _ := json.MarshalIndent(group,"","  ")
      w.Write(b)
    } else {
      b, _ := json.MarshalIndent(map[string]string{"succ":"false"},"","  ")
      w.Write(b)
    }
    w.Header().Set("Content-Type","text/json")
}

func StartUp(logger mlog.MLog,port string, get func(key string)(string,error), put func(string, string, map[string]bool) (error,map[string]string)) {
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
