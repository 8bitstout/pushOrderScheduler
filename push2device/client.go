package push2device

import (
	"fmt"
	"net/http"
	"sync"
)

type Push2Device struct {
	clientMu sync.Mutex
	client   *http.Client
}

func (p *Push2Device) CreatePushNotification(orderId string) error {
	fmt.Println("Order", orderId, "push notification registered")
	return nil
}

func MakePush2DeviceClient() *Push2Device {
	return &Push2Device{}
}
