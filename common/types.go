package common

type RemoteServer struct{
  Address string
  Port string
  ID int64
}


type Tuple struct {
        Key, Value string
}

type Hashtable struct {
	//can add metdata? (Karl)
	Size int
	Store map[string]string
}
