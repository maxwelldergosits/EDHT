package web_interface

type WebDelegate interface {
  GetF(key string) (string,error)
  PutF(string,string,map[string]bool) (error,map[string]string)
  Info(int) []uint64
  Topology() ([][]uint64,[]uint64)
}
