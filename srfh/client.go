package srfh

import (
	"context"
	"fmt"
	pb "github.com/8bitstout/orderPushScheduler/order"
	"google.golang.org/grpc"
	"log"
	"net/url"
	"time"
)

type Client struct {
	pb.ScheduleOrderPushClient
	URL  *url.URL
	conn *grpc.ClientConn
}

func MakeClient(u string) *Client {
	parsedURL, _ := url.Parse(u)
	c := &Client{
		URL: parsedURL,
	}
	fmt.Println("Client connecting to:", parsedURL.String())
	conn, err := grpc.Dial(parsedURL.String(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	//defer conn.Close()
	return c
}

func (c *Client) Schedule() {
	order := &pb.Order{
		Id: "10",
	}
	fmt.Println("New Order: ", order.String())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cli := pb.NewScheduleOrderPushClient(c.conn)
	r, err := cli.SchedulePushNotification(ctx, order)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}
