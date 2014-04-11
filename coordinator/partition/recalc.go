package partition

func init() {

  f = .2
}

func pd(e,a uint) float32 {

  af := float32(a)
  ef := float32(e)

  return (ef - af)*2 / (ef + af)

}

var f float32

func (self *Shard) Copy() *Shard {
    r := new(Shard)
    *r = *self
    return r
}

func (t* PartitionSet) Recalc() {


  t.GatherInfo()
  for i:=1; i < len(t.shards); i++ {

    shard_1 := t.shards[i-1].Copy()
    shard_2 := t.shards[i].Copy()

    s1 := shard_1.Start
    e1 := shard_1.End

    s2 := shard_2.Start
    e2 := shard_2.End

    k1 := shard_1.Keys
    k2 := shard_2.Keys

    pdiff := pd(k1,k2)

    if pdiff < -.05 {
      e_1n := uint64(-pdiff * f * float32(e2-s2) + float32(e1))
      s_2n := e_1n+1

      t.shards[i].Start = s_2n
      t.shards[i-1].End = e_1n

    }else if pdiff > .05 {
      e_1n := uint64(-pdiff * f * float32(e1-s1) + float32(e1))
      s_2n := e_1n+1

      t.shards[i].Start = s_2n
      t.shards[i-1].End = e_1n
    }

  }

  t.UpdateInfo()

}
