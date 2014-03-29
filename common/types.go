package common

type RemoteServer struct{
  Address     string
  Port        string
  ID          uint64
  Coordinator bool
}

type DaemonData struct{
  ServerDetails RemoteServer
  Hashtable     Hashtable
}

type RegisterReply struct {
  Coordinators     map[uint64]RemoteServer
  Daemons          map[uint64]RemoteServer
  ID               uint64
  Nshards          uint
  Nfailures        uint
}


type Tuple struct {
  Key   string
  Value string
}

type Hashtable struct {
	Size int
	Store map[string]string
}

type Shard struct {
  Start uint64
  End uint64
  Daemons map[uint64]bool
}
