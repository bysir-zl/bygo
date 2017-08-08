package discovery

import (
	"testing"
	"github.com/bysir-zl/bygo/log"
	"time"
)

var etcd = NewEtcd([]string{"127.0.0.1:2380"}, 30)

func TestGetServer(t *testing.T) {

	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			ser, err := etcd.GetServers()
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range ser {
				log.InfoT("test", s)
			}
		}
	}
}

func TestUpdateTTl(t *testing.T) {
	err := etcd.UpdateServerTTL("7587823890455770676")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterService(t *testing.T) {
	leaseId, err := etcd.RegisterService(&Server{
		Id:      "game-1",
		Name:    "game-1",
		Port:    8090,
		Address: "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)
	err = etcd.UpdateServerTTL(leaseId)
	if err != nil {
		t.Fatal(err)
	}

}

func TestUnRegisterService(t *testing.T) {
	err := etcd.UnRegisterService("game-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWatch(t *testing.T) {

	etcd.WatchServer(func(server *Server, change ServerChange) {
		log.InfoT("test", server, change)
	})
	select {}
}
