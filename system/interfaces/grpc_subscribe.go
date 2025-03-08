package interfaces

typr GrpcSubscribe interface {
	GetStream() GrpcStream
	GetStreamBidirectional() GrpcStreamBidirectional
	IsFinished() chan bool
}