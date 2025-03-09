package interfaces

type NodeConfig interface {
	Init() error
	Reload() error
	GetANode() Entity
	GetAuthKey() string

	GetNodeName() string
	IsMaster() bool
}
