package main

func GetKey(key string) (string,error) {
  pts := gc.GetPartitions()
  return pts.Get(key)
}

func PutKey(key string, value string, options map[string]bool) (error,map[string]string) {
  pts := gc.GetPartitions()
  return pts.Put(key,value,options)
}
