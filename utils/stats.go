package utils


func Mean(in []uint) float64 {

  s := 0.0
  for _,v := range in {
    s += float64(v)
  }
  return (s/float64(len(in)))
}


func StdDev(in []uint) float64 {
  mean := Mean(in)
  s := 0.0
  for _,iv := range in {
    v := float64(iv)
    s += ((mean - v)*(mean -v))
  }
  return s/float64(len(in))
}

func CV(in []uint) float64 {

  mean := Mean(in)
  sd   := StdDev(in)

  return (sd/mean)

}
