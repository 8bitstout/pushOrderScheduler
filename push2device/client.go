package push2device

import (
	"fmt"
	pb "github.com/8bitstout/orderPushScheduler/proto"
	"github.com/golang/protobuf/jsonpb"
	"net/http"
	"sync"
)

type Push2Device struct {
	clientMu sync.Mutex
	client   *http.Client
}

func (p *Push2Device) CreateOrderPush(pushOrder *pb.Order) {
	m := jsonpb.Marshaler{EmitDefaults: true}
	j, _ := m.MarshalToString(pushOrder)
	//r, _ := http.NewRequest("POST", "endpoint", bytes.NewBuffer([]byte(j)))
	//p.client.Do(r)
	fmt.Println("Creating push2device request: ", j)
}
