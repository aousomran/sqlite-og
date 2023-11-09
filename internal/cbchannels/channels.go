package cbchannels

import pb "gitlab.com/aous.omr/sqlite-og/gen/proto"

type CallbackChannels struct {
	ChanSend    chan *pb.Invoke
	ChanReceive chan *pb.InvocationResult
}

func New() *CallbackChannels {
	return &CallbackChannels{
		ChanSend:    make(chan *pb.Invoke),
		ChanReceive: make(chan *pb.InvocationResult),
	}
}
