package types

import "github.com/PretendoNetwork/nex-go"

type MetaBinary struct {
	DataID            uint32
	OwnerPID          uint32
	Name              string
	DataType          uint16
	Buffer            []byte
	Permission        uint8
	DeletePermission  uint8
	Flag              uint32
	Period            uint16
	Tags              []string
	PersistenceSlotID uint16
	ExtraData         []string
	CreationTime      *nex.DateTime
	UpdatedTime       *nex.DateTime
	ReferredTime      *nex.DateTime
	ExpireTime        *nex.DateTime
}

func NewMetaBinary() *MetaBinary {
	return &MetaBinary{}
}
