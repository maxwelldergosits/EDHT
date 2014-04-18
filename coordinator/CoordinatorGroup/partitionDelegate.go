package CoordinatorGroup
import (
  "EDHT/coordinator/CoordinatorGroup/group"
  . "EDHT/common"
)
type PD struct {
  gms group.Group
}

func (pd * PD) GetDaemon(id uint64) RemoteServer {
  return pd.gms.GetDaemon(id)
}

func (pd * PD) GetLocalID() uint64 {
  return pd.gms.GetID()
}

func (pd * PD) DeleteDaemon(id uint64) {
  pd.gms.Delete(pd.gms.GetDaemon(id))
}
