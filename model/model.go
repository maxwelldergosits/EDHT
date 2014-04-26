package main

import "EDHT/coordinator/CoordinatorGroup/partition"
import  . "EDHT/common"
import "fmt"
import "os"
import "math/rand"
import "strconv"
import "github.com/mad293/mlog"
import "encoding/json"

type PD int

var ml mlog.MLog
var n int
var iterations int
var numkeys int


func (pd * PD) GetDaemon(id uint64) RemoteServer {
  return RemoteServer{}
}
func (pd * PD) DeleteDaemon(id uint64) {
}
func (pd * PD) GetLocalID() uint64{
  return 0
}

func (pd * PD) Logger() *mlog.MLog {
  return &ml
}

func main() {



  ml = mlog.Create([]string{},"",true,true)
  n, _ = strconv.Atoi(os.Args[1])
  iterations, _ = strconv.Atoi(os.Args[2])
  numkeys, _ = strconv.Atoi(os.Args[3])
  pts := partition.MakeKeySpace(n,new(PD))


  initDiff()

  data := make([]map[string]uint64,n)
  for i:= 0; i < n; i++ {
    data[i] = make(map[string]uint64)
  }

  for t := 0; t < iterations; t++ {

    ranges := pts.Ranges()
    sizes := make([]uint64,len(ranges)/2)
    str := strconv.FormatInt(int64(t),10)
    keys := GetKeys(ranges)

    for i:= 0; i<len(ranges); i+=2 {
      sizes[i/2] = ranges[i+1] - ranges[i]
      //data[i/2]["size"+str] = sizes[i/2]
      data[i/2]["keys"+str] = uint64(keys[i/2])
    }

    pts = doDiff(pts)
  }

    b, _ := json.MarshalIndent(data,"","  ")

    fmt.Println("data:")
    os.Stdout.Write(b)

}


var mockedData map[uint64]bool

func initDiff() {


  mockedData = make(map[uint64]bool)

  for i := 0; i < numkeys; i++ {

    //ri := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
    ri := uint64(rand.Uint32())
    fmt.Println(ri)
    mockedData[ri] = true

  }

}

func GetKeys(ranges []uint64) []uint {

  keys := make([]uint,n)

  for k,_ := range mockedData {

    for i := 0; i < n; i++ {
      start := ranges[2*i]
      end := ranges[(2*i)+1]
      if start <= k && k <= end {
        keys[i] +=1
        break
      }
    }

  }

  return keys
}

func doDiff(pts *partition.PartitionSet) * partition.PartitionSet {

  ranges := pts.Ranges()


  keys := GetKeys(ranges)

  newpts := pts.Recalc(keys)
  return &newpts


}
