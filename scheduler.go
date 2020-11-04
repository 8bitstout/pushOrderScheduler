package orderPushScheduler

import (
	pb "github.com/8bitstout/orderPushScheduler/proto"
	"github.com/8bitstout/orderPushScheduler/push2device"
)

type Scheduler struct {
	Server      *Server
	Push2Device *push2device.Push2Device
}

func MakeScheduler(port string) *Scheduler {
	return &Scheduler{
		Server:      MakeServer(port),
		Push2Device: &push2device.Push2Device{},
	}
}

func MakeOrder(id int32) *pb.Order {
	return &pb.Order{Id: id}
}

func (s *Scheduler) Run() {
	s.Server.Run()
}
