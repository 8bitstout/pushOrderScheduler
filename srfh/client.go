package srfh

import (
	"fmt"
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
		}
	}()
}

func (c *Client) WriteMessage(m string) {
	fmt.Println("Writing message to:", c.conn.RemoteAddr().String())
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		log.Println(err)
	}
}
