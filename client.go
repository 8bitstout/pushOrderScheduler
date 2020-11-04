package orderPushScheduler

import (
	"fmt"
	pb "github.com/8bitstout/orderPushScheduler/proto"
	"github.com/8bitstout/orderPushScheduler/push2device"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	logInfo     *log.Logger
	logError    *log.Logger
	push2device *push2device.Push2Device
}

func MakeClient(h *Hub, c *websocket.Conn) *Client {
	return &Client{
		Hub:         h,
		conn:        c,
		send:        make(chan []byte, 256),
		logInfo:     log.New(os.Stdout, "INFO:Client:\t", log.Ldate|log.Ltime),
		logError:    log.New(os.Stdout, "ERROR:Client:\t", log.Ldate|log.Ltime),
		push2device: &push2device.Push2Device{},
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logError.Println(err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Println(message)
		order := &pb.Order{}
		err = proto.Unmarshal(message, order)
		if err != nil {
			c.logError.Println("ReadPump:Unmarshal:", err)
		}
		fmt.Println("Received order ", order.Id, "created at", order.CreatedAt)
		c.logInfo.Println(order.String())
		c.Hub.Broadcast(message)
	}
}

// Write To Push2Device
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message) // write if ws existed on other end

			// expect Order protobuf
			order := pb.Order{}
			err = proto.Unmarshal(message, &order)
			if err != nil {
				c.logError.Println("Error unmarshalling protobuf:", err)
			}
			c.push2device.CreateOrderPush(&order)

			// Add queued message to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
