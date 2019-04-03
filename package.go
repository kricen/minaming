package minaming

import (
	"fmt"
	"net"
	"sync"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

type micro struct {
	serviceName string
	port        int
	// closeFuncChan chan closeFunc
	endpoints    []string
	serviceConns sync.Map
	microNaming  MicroNaming
	sync.WaitGroup
}

// NewMicro : initial a
func NewMicro(serviceName, ip string, port int, endpoints []string) *micro {
	var mn MicroNaming
	if serviceName != "" && port != 0 {
		mn = NewMicroNaming(endpoints)
		err := mn.Register(serviceName, fmt.Sprintf("%s:%d", ip, port))
		if err != nil {
			panic(err.Error)
		}
	}
	m := &micro{serviceName: serviceName, port: port, endpoints: endpoints, serviceConns: sync.Map{}, microNaming: mn}

	return m
}

//ReferServices
//package all  grpc conntions to maps ，so that you can use it conveniently next time
func (m *micro) ReferServices(svcNames ...string) {
	m.startServiceConns(svcNames, m.endpoints)
}

func (m *micro) CreateListener() net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", m.port))
	if err != nil {
		panic(err.Error())
	}

	return lis
}

// startServiceConns:
func (m *micro) startServiceConns(serverList []string, endpoints []string) {
	for _, serviceName := range serverList {

		// register resolver
		enr := NewResolverBuilder(endpoints)
		resolver.Register(&enr)
		address := fmt.Sprintf("%s://%s/%s", _DEFAULT_SCHEME, _AUTHOR, serviceName)
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))

		if err != nil {
			fmt.Printf(`connect to '%s' service failed: %v`, serviceName, err)
		}
		fmt.Printf(`connected to '%s' `, serviceName)
		m.serviceConns.Store(serviceName, conn)
	}
}

// CloseService:
// 1: remove server from balancer
// 2: close all established conns
func (m *micro) CloseService() {
	m.microNaming.Unregister()
	m.serviceConns.Range(func(key interface{}, value interface{}) bool {
		conn := value.(*grpc.ClientConn)
		conn.Close()
		return true
	})
	log.Println("GoodBye,See you next time")
}

func (m *micro) GetConn(serverName string) (*grpc.ClientConn, bool) {
	connInterface, ok := m.serviceConns.Load(serverName)
	if !ok {
		return nil, ok
	}

	return connInterface.(*grpc.ClientConn), true
}
