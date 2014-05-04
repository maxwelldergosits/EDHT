package main

import "EDHT/coordinator/CoordinatorGroup/partition"
import  . "EDHT/common"
import "fmt"
import "os"
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
  n, _ = strconv.Atoi(os.Args[1]) // Number of partitions
  iterations, _ = strconv.Atoi(os.Args[2]) // iterations
  numkeys, _ = strconv.Atoi(os.Args[3]) // keys
  pts := partition.MakeKeySpace(n,new(PD))


  initDiff()

  data := make([]map[string]uint64,n)
  for i:= 0; i < n; i++ {
    data[i] = make(map[string]uint64)
  }

  for t := 0; t < iterations; t++ {

    ranges := pts.Ranges()
    iranges := pts.IntRanges()
    sizes := make([]uint64,len(iranges)/2)
    str := strconv.FormatInt(int64(t),10)
    keys := GetKeys(ranges)

    for i:= 0; i<len(iranges); i+=2 {
      if (t == iterations-1) {
      sizes[i/2] = iranges[i+1] - iranges[i]
      data[i/2]["keys"+str] = uint64(keys[i/2])
      }
    }

    pts = doDiff(pts)
  }

    b, _ := json.MarshalIndent(data,"","  ")

    fmt.Println("data:")
    os.Stdout.Write(b)

}


var mockedData map[string]bool

func initDiff() {


  mockedData = make(map[string]bool)

  for i := 0; i < numkeys; i++ {

    mockedData[strconv.Itoa(i)] = true

  }

}

func GetKeys(ranges []Range) []uint {

  keys := make([]uint,n)

  for k,_ := range mockedData {

    for i := 0; i < n; i++ {
      end := ranges[i].End
      start := ranges[i].Start
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
