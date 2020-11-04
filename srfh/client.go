package srfh

import (
	"fmt"
	"github.com/8bitstout/orderPushScheduler"
	pb "github.com/8bitstout/orderPushScheduler/proto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Client struct {
	URL  *url.URL
	conn *websocket.Conn
}

func MakeClient(u string) *Client {
	parsedURL, _ := url.Parse(u)
	c := &Client{
		URL: parsedURL,
	}
	fmt.Println("Client connecting to:", parsedURL.String())
	conn, _, err := websocket.DefaultDialer.Dial(parsedURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	return c
}

func (c *Client) ListenForMessages() {
	go func() {
		for {
			msgType, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}
			if msgType == websocket.TextMessage {
				fmt.Println(string(msg))
			}
			if msgType == websocket.BinaryMessage {
				order := &pb.Order{}
				proto.Unmarshal(msg, order)
				fmt.Println("Received order ", order.Id, "created at", order.CreatedAt)
			}
		}
	}()
}

func (c *Client) WriteMessage(m []byte) {
	fmt.Println("Writing message to:", c.conn.RemoteAddr().String())
	err := c.conn.WriteMessage(websocket.BinaryMessage, m)
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) SendNewOrder() {
	order := orderPushScheduler.MakeOrder(10)
	fmt.Println("New Order: ", order.String())
	message, err := proto.Marshal(order)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(message)
	c.WriteMessage(message)
}
