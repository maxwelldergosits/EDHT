package main

import (
  . "EDHT/common"
)

type DaemonGroup struct {
}

func (dg * DaemonGroup) Put(kvs []Tuple, options map[string]string) map[string]string {
  return nil
}

func (dg * DaemonGroup) Get(ks []string) (error, []Tuple) {
  return nil,nil
}

func (dg * DaemonGroup) Delete(ks []string){
  
}

func (dg * DaemonGroup) Add(d RemoteServer) error{
  return nil
}

