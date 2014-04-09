package utils

import (
  "errors"
)

type Pmap struct {
  data map[string]string
  directory string
}

func (pm * Pmap) Put(key string, value string, replace bool) bool {
  var flag bool
  if _, exists := pm.data[key];exists {
    if ! replace {
      flag =  false
    } else {
      pm.data[key]= value
      flag =  true
    }
  } else {
      pm.data[key]= value
      flag =  true
  }
  return flag
}

func (pm * Pmap) Get(key string) (string,error) {
  var v string
  var e error
  if val,exists := pm.data[key]; exists {
    v,e = val,nil
  } else {
    v,e = "",errors.New("Key doesn't exist")
  }
  return v,e
}


func NewPmap(dirname string) Pmap {
  return Pmap{make(map[string]string),dirname}
}


