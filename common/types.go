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

type ShardCopy struct {
  Start uint64
  End uint64
  Daemons map[uint64]bool
}

type ConnectReply struct {
  RegReply RegisterReply
  Partitions []ShardCopy
}

type Tuple struct {
  Key   string
  Value string
}

type Hashtable struct {
	Size int
	Store map[string]string
}

type Range struct {
  Start string
  End string
}

type ServerRange struct {
  Server RemoteServer
  Range Range
}

type PutRequest struct{
  Key string
  Value string
  Options map[string]bool
}
