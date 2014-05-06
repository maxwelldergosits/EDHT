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
  "EDHT/utils"
  "strings"
  "encoding/json"
  "strconv"
  "github.com/mad293/mlog"
)

var (
  delegate WebDelegate
  ml mlog.MLog
)

func gethandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprintf(w, getForm)
}
func getshandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    key:=r.FormValue("key")
    value,err := delegate.GetF(key)
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
    w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprintf(w, putForm)
}
func shandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    key:= r.FormValue("key")
    if (key == "horses") {
      fmt.Fprintf(w,easter_egg)
      return
    }
    w.Header().Set("Content-Type","text/json")
    value:= r.FormValue("value")
    getOV:= r.FormValue("ov")
    getOVbool,_ := strconv.ParseBool(getOV)
    unsafe, _ := strconv.ParseBool(r.FormValue("unsafe"))
    option:= map[string]bool{"ov":getOVbool}
    option["unsafe"] = unsafe
    err,values := delegate.PutF(key,value,option)
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
    }else if (unsafe) {
      b, _ := json.MarshalIndent(map[string]string{"unsafe":"true"},"","  ")
      w.Write(b)
    }else {
      b, _ := json.MarshalIndent(map[string]string{"succ":"false"},"","  ")
      w.Write(b)
    }
    w.Header().Set("Content-Type","text/json")
}


func infoHandler(w http.ResponseWriter, r * http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type","text/json")
  uri := r.RequestURI
  if strings.Contains(uri,"keys") {
    keys := delegate.Info(1)
    b, _ := json.MarshalIndent(keys,"","  ")
    w.Write(b)
  } else if strings.Contains(uri,"heat") {
    keys := delegate.Info(1)
    sizes := delegate.Info(2)
    out := make([]map[string]uint64,len(keys))
    for i := range keys {
      out[i] = make(map[string]uint64)
      out[i]["keys"] = keys[i]
      out[i]["size"] = sizes[i]
    }
    b, _ := json.MarshalIndent(out,"","  ")
    w.Write(b)
  } else if  strings.Contains(uri,"topology") {
    d,c := delegate.Topology()
    db, _ := json.Marshal(d)
    dj :=string(db)
    cb, _ := json.Marshal(c)
    cj :=string(cb)
    out := map[string]string{"daemons":dj,"coordinators":cj}
    b, _ := json.MarshalIndent(out,"","  ")
    w.Write(b)
  } else {
    f := func(in float64)string {return strconv.FormatFloat(in,'f',5,64)}
    keys := delegate.Info(1)
    out := map[string]string{}
    out["Standard Deviation"] = f(utils.StdDev(keys))
    out["std"] = f(utils.StdDev(keys))
    out["Mean"] = f(utils.Mean(keys))
    out["Coefficient of Variation"] = f(utils.StdDev(keys)/utils.Mean(keys))
    out["cv"] = f(utils.StdDev(keys)/utils.Mean(keys))
    for i,v := range utils.Devs(keys) {
      out[strconv.Itoa(i)+" dev"] = f(float64(v))
    }
    b, _ := json.MarshalIndent(out,"","  ")
    w.Write(b)
  }
}

func genHandler(w http.ResponseWriter, r * http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprintln(w,requests)
}


func StartUp(logger mlog.MLog,port string,del WebDelegate) {
  ml = logger
  delegate = del
  ml.VPrintln("web","starting web inteface")

  fs := http.FileServer(http.Dir(del.Dir()))


  http.HandleFunc("/put",handler)
  http.HandleFunc("/put/submit",shandler)
  http.HandleFunc("/get",gethandler)
  http.HandleFunc("/get/submit",getshandler)
  http.HandleFunc("/stats/balance/",infoHandler)
  http.HandleFunc("/gen/",genHandler)
  nfs := http.StripPrefix("/view/", fs)
  http.Handle("/view/",nfs)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
