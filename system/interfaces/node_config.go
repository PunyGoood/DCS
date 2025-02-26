package interfaces

type NodeConfig interface {
	Init() error
	Reload() error
	GetANode() entity
	GetAuthKey() string
	Address string
}
