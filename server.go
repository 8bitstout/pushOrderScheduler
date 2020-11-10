package orderPushScheduler

import (
	"context"
	pb "github.com/8bitstout/orderPushScheduler/order"
	"github.com/8bitstout/orderPushScheduler/push2device"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Scheduler struct {
	pb.UnimplementedScheduleOrderPushServer
	orderMap map[string]*pb.Order
	logInfo  *log.Logger
	logError *log.Logger
}

func (s *Scheduler) SchedulePushNotification(ctx context.Context, in *pb.Order) (*pb.Result, error) {
	s.logInfo.Println("Scheduling New Order:", in.GetId())
	c := push2device.MakePush2DeviceClient()
	err := c.CreatePushNotification(in.GetId())
	if err != nil {
		return &pb.Result{Success: false, Response: "Failed to register push notification"}, err
	}

	return &pb.Result{Success: true, Response: "Push notification scheduled for Order " + in.Id}, nil
}

func (s *Scheduler) Run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("|failed to listen: %v", err)
	}

	g := grpc.NewServer()
	pb.RegisterScheduleOrderPushServer(g, s)
	s.logInfo.Println("Schedule Order Push server registered")
	if err = g.Serve(lis); err != nil {
		log.Fatalf("Failed to server %v", err)
	}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		orderMap: make(map[string]*pb.Order),
		logInfo:  log.New(os.Stdout, "INFO:Scheduler:", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR:Scheduler:", log.Ldate|log.Ltime),
	}
}
