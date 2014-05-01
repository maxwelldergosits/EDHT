package utils

import (
 // "errors"
  . "EDHT/common"
	"github.com/peterbourgon/diskv"
  "strconv"
  "encoding/gob"
  "bytes"
)

type StringStringMap struct {
  d *diskv.Diskv
}

type Uint64ServerMap struct {
  d * diskv.Diskv
}



/* StringString Map Functions */

func NewStringStringMap (dirname string) StringStringMap{
  d := diskv.New(diskv.Options{
    BasePath:     dirname,
    Transform:    func(s string) []string { return []string{} },
    CacheSizeMax: 1024 * 1024, // 1MB
  })
  return StringStringMap{d}
}

func NewStringStringFromMap (dirname string,old map[string]string) StringStringMap {
  d := diskv.New(diskv.Options{
    BasePath:     dirname,
    Transform:    func(s string) []string { return []string{} },
    CacheSizeMax: 1024 * 1024, // 1MB
  })

  m := StringStringMap{d}
  for k,v := range old {
    m.Put(k,v)
  }
  return m
}

func (m * StringStringMap) Put(key, value string) error {
  bs := stringToBytes(value)
  return m.d.Write(key,bs)
}

func (m* StringStringMap) Get(key string) (string,error) {
  bs, err := m.d.Read(key)
  if err != nil {return "",err}
  return bytesToString(bs)
}

func (m* StringStringMap) Delete(key string) error {
  return m.d.Erase(key)
}

func (m * StringStringMap) Map() map[string]string {
  keyChan := m.d.Keys()
  out := make(map[string]string)
  for key := range keyChan {
    val, err := m.d.Read(key)
    if (err != nil) {continue}
    rs, err := bytesToString(val)
    if (err != nil) {continue}
    out[key] = rs
  }
  return out
}
/* Uint64Server Map functions */
func (m * Uint64ServerMap) Put(key uint64, rs RemoteServer) error {
  bs := remoteServerToBytes(rs)
  k := strconv.FormatUint(key,10)
  return m.d.Write(k,bs)
}


func NewUint64ServerMap (dirname string) Uint64ServerMap{
  d := diskv.New(diskv.Options{
    BasePath:     dirname,
    Transform:    func(s string) []string { return []string{} },
    CacheSizeMax: 1024 * 1024, // 1MB
  })
  return Uint64ServerMap{d}
}

func NewUint64ServerFromMap (dirname string,old map[uint64]RemoteServer) Uint64ServerMap{
  d := diskv.New(diskv.Options{
    BasePath:     dirname,
    Transform:    func(s string) []string { return []string{} },
    CacheSizeMax: 1024 * 1024, // 1MB
  })

  m := Uint64ServerMap{d}
  for k,v := range old {
    m.Put(k,v)
  }
  return m
}


func (m* Uint64ServerMap ) Get(key uint64) (RemoteServer,error) {
  k := strconv.FormatUint(key,10)
  bs, err := m.d.Read(k)
  if err != nil {return RemoteServer{},err}
  return bytesToRemoteServer(bs)
}

func (m* Uint64ServerMap) Delete(key uint64) error {
  k := strconv.FormatUint(key,10)
  return m.d.Erase(k)
}

func (m * Uint64ServerMap) Map() map[uint64]RemoteServer {
  keyChan := m.d.Keys()
  out := make(map[uint64]RemoteServer)
  for key := range keyChan {
    val, err := m.d.Read(key)
    if (err != nil) {continue}
    rs, err := bytesToRemoteServer(val)
    if (err != nil) {continue}
    k, err := strconv.ParseUint(key,10,64)
    if (err != nil) {continue}
    out[k] = rs
  }
  return out
}

/* Utility Functions */

func remoteServerToBytes(rs RemoteServer) []byte {
  w:= new(bytes.Buffer)
  encoder := gob.NewEncoder(w)
  encoder.Encode(rs)
  return w.Bytes()
}

func bytesToRemoteServer(bs []byte) (RemoteServer, error) {
  rs := new(RemoteServer)
  w := bytes.NewBuffer(bs)
  decoder := gob.NewDecoder(w)
  err := decoder.Decode(rs)
  return *rs,err
}

func stringToBytes(s string) []byte {
  w:= new(bytes.Buffer)
  encoder := gob.NewEncoder(w)
  encoder.Encode(s)
  return w.Bytes()
}

func bytesToString(bs []byte) (string,error) {
  w := bytes.NewBuffer(bs)
  decoder := gob.NewDecoder(w)
  var s string
  err := decoder.Decode(&s)
  return s,err
}
