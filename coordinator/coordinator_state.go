/*
This file serves as a catch all for the state that needs to be shared by the coordinators


Any change to this state will be done via a two phase commit


*/

package main
import (
  . "EDHT/common"
)

type CoordinatorState struct {

  Coordinators map[int64]RemoteServer

  Daemons map[int64]RemoteServer

  Shards map[int]Shard

  
}
