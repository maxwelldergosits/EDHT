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
