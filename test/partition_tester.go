package main

import (
  "EDHT/coordinator/partition"
  "math"
)

var keysMock []int

var n = 3000

func f(i int) int {

  x := float64((i-1500)/n)

  return int(math.Exp(-math.Pow(x,2)))

}

func main() {



  keysMock = make([]int,n,n)

  partitions := partition.MakeKeySpace(n,gatherInfo,updateShard)


  for i in range

  for i:=0; i< n; i++ {

    keysMock[i] = f(i)

  }

}


func updateShard(s *partition.Shard) {


}

func gatherInfo(s *partition.Shard) {


}



