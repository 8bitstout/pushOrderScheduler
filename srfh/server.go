package srfh

import "github.com/8bitstout/orderPushScheduler/tcp"

type Srfh struct {
	Server *tcp.Server
}

func MakeSrfh(port string) *Srfh {
	return &Srfh{
		Server: tcp.MakeServer(port),
	}
}

func (s *Srfh) Run() {
	s.Server.Run()
}
