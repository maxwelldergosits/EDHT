package utils

import (
  "errors"
)

type Pmap struct {
  data map[string]string
  directory string
}

func (pm * Pmap) Put(key string, value string, replace bool) bool {
  if _, exists := pm.data[key];exists {
    if ! replace {
      return false
    } else {
      pm.data[key]= value
      return true
    }
  } else {
      pm.data[key]= value
      return true
  }
}

func (pm * Pmap) Get(key string) (string,error) {

  if val,exists := pm.data[key]; exists {
    return val,nil
  } else {
    return "",errors.New("Key doesn't exist")
  }

}


func NewPmap(dirname string) Pmap {
  return Pmap{make(map[string]string),dirname}
}


