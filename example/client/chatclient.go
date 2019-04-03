package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kricen/minaming"
	"github.com/kricen/minaming/example/pb"
	"google.golang.org/grpc/resolver"
)

var (
	wg sync.WaitGroup
)

func init() {
	// register resolver
	enr := minaming.NewResolverBuilder([]string{"localhost:2379"})
	resolver.Register(&enr)

}
func main() {
	fmt.Println("start", time.Now().Unix())
	startWithWroc()
	fmt.Println("end", time.Now().Unix())
}

func startWithWroc() {
	m := minaming.NewMicro("", "", 0, []string{"http://localhost:2379"})
	m.ReferServices("chat_server")
	conn, ok := m.GetConn("chat_server")
	if !ok {
		fmt.Println("notExists")
		return
	}
	fmt.Println(conn == nil)
	chatClient := pb.NewChatServerClient(conn)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := chatClient.SayHello(context.Background(), &pb.GreetReq{Name: "orican", Msg: "hello"})
			if err != nil {
				fmt.Println(err.Error)
			}
			fmt.Println(resp.GetMsg())
		}()
	}
	//
	// resp, err := chatClient.SayHello(context.Background(), &pb.GreetReq{Name: "orican", Msg: "hello1"})
	// if err != nil {
	// 	fmt.Println(err.Error)
	// }
	// fmt.Println(resp.GetMsg())
	// time.Sleep(5 * time.Second)
	// resp, err = chatClient.SayHello(context.Background(), &pb.GreetReq{Name: "orican", Msg: "hello2"})
	// if err != nil {
	// 	fmt.Println(err.Error)
	// }
	// fmt.Println(resp.GetMsg())
	// time.Sleep(5 * time.Second)
	//
	// resp, err = chatClient.SayHello(context.Background(), &pb.GreetReq{Name: "orican", Msg: "hello3"})
	// if err != nil {
	// 	fmt.Println(err.Error)
	// }
	// fmt.Println(resp.GetMsg())
	wg.Wait()

}
