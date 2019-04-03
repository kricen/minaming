package minaming

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	_DEFAULT_SCHEME = "minaming"
	_AUTHOR         = "kricen"
)

var (
	cli *clientv3.Client
)

type MicroNaming struct {
	Endpoints  []string
	Scheme     string
	serverName string
	addr       string
}

// NewHi create a Hi instance
func NewMicroNaming(endpoints []string) MicroNaming {
	return MicroNaming{Endpoints: endpoints, Scheme: _DEFAULT_SCHEME}
}

// Unregiste delete name from etcd
func (m *MicroNaming) Unregister() error {
	key := fmt.Sprintf("%s/%s/%s", m.Scheme, m.serverName, m.addr)
	fmt.Println(key)
	_, err := cli.Delete(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

// Register register microserver address under the scheme/name
func (m *MicroNaming) Register(name, addr string) error {
	var err error

	// create client if not define
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   m.Endpoints,
			DialTimeout: 10 * time.Second,
		})
		if err != nil {
			return err
		}
	}

	// create lease with 12s TTL
	leaseResp, err := cli.Grant(context.Background(), 12)
	if err != nil {
		return err
	}

	// put key & value
	key := fmt.Sprintf("%s/%s/%s", m.Scheme, name, addr)
	_, err = cli.Put(context.Background(), key, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	m.serverName = name
	m.addr = addr

	// generate a goroutine to keep lease alive per 10s
	go func() {
		for range time.Tick(time.Second * 10) {
			cli.KeepAliveOnce(context.Background(), leaseResp.ID)
		}
	}()

	return nil
}
