package orderPushScheduler

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/8bitstout/orderPushScheduler/order"
	"github.com/8bitstout/orderPushScheduler/push2device"
	"github.com/dgraph-io/badger/v2"

	"google.golang.org/grpc"
)

const (
	port          = ":50051"
	STATE_PENDING = iota
	STATE_SCHEDULED
	STATE_FAILED
)

type Scheduler struct {
	pb.UnimplementedScheduleOrderPushServer
	orderMap map[string]*pb.Order
	db       *badger.DB
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

	r := &pb.Result{Success: true, Response: "Push notification sheduled for order" + in.Id}
	fmt.Println(r)
	return r, nil
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
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return &Scheduler{
		orderMap: make(map[string]*pb.Order),
		logInfo:  log.New(os.Stdout, "INFO:Scheduler:", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR:Scheduler:", log.Ldate|log.Ltime),
		db:       db,
	}
}
