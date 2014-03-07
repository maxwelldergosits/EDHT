package utils

import (
  "testing"
)

func TestGen_id(t * testing.T) {

  var x = gen_id(0,0,0)

  if (x != 0) {
    t.Error("isn't 0")
  }



}
