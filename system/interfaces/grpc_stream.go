package interfaces

import("github.com/PunyGoood/DCS/grpc/proto")

typr GetStream interface {

	Send(msg *grpc.StreamMessage) (err error)
}

typr GetStreamBidirectional interface {

}