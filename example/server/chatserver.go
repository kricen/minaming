package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kricen/minaming"
	"github.com/kricen/minaming/example/pb"
	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

var (
	serviceName = "chat_server"
	num         int64
)

type ChatServer struct {
}

func (s *ChatServer) SayHello(ctx context.Context, req *pb.GreetReq) (*pb.GreetResq, error) {
	// fmt.Printf("from user:%s,he said:%s\n", req.GetName(), req.GetMsg())
	atomic.AddInt64(&num, 1)
	return &pb.GreetResq{Msg: fmt.Sprintf("Hi %s thanks for your greet", req.GetName())}, nil
}

func main() {
	port := flag.Int("port", 10010, "http listen port")
	flag.Parse()
	fmt.Println(*port)
	wrocStart(*port)
}

func wrocStart(port int) {
	mc := minaming.NewMicro(serviceName, "127.0.0.1", port, []string{"http://localhost:2379"})
	s := grpc.NewServer()
	pb.RegisterChatServerServer(s, &ChatServer{})
	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Println(num)
		}
	}()
	println("start chatServer")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		mc.CloseService()
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()
	if err := s.Serve(mc.CreateListener()); err != nil {
		log.Fatal(err)
		panic(err)
	}

}
