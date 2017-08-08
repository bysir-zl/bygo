package discovery

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"context"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"strings"
	"github.com/bysir-zl/bygo/log"
	"strconv"
)

var root = "/service/"

type Etcd struct {
	cli *clientv3.Client
	ttl int64

	etcdEndpoints []string
}

func (p *Etcd) GetServers() (servers map[string]*Server, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rsp, err := p.cli.Get(ctx, root, clientv3.WithPrefix())
	if err != nil {
		return
	}
	servers = make(map[string]*Server, len(rsp.Kvs))
	for _, kv := range rsp.Kvs {
		server := Server{}
		err = json.Unmarshal(kv.Value, &server)
		if err != nil {
			return
		}
		servers[server.Id] = &server
	}

	return
}

func (p *Etcd) WatchServer(fun ServersChanged) {
	ctx := context.TODO()
	ch := p.cli.Watch(ctx, root, clientv3.WithPrefix())
	go func() {
		for {
			select {
			case c := <-ch:
				for _, e := range c.Events {
					change := SC_Online
					server := Server{}
					if e.Type == mvccpb.PUT {
						json.Unmarshal(e.Kv.Value, &server)
					} else {
						change = SC_Offline
						server.Id = strings.Split(string(e.Kv.Key), root)[1]
					}

					fun(&server, change)
				}
			}
		}
	}()

	// 通知所有服务上线
	servers, err := p.GetServers()
	if err != nil {
		log.ErrorT("etcd", err)
		return
	}
	for _, s := range servers {
		fun(s, SC_Online)
	}
}

func (p *Etcd) UpdateServerTTL(leaseId string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	id, err := strconv.ParseInt(leaseId, 10, 64)
	if err != nil {
		return
	}
	r, err := p.cli.KeepAlive(ctx, clientv3.LeaseID(id))
	if err != nil {
		return
	}
	// 这里必须接受并取出才能KeepAlive, 我也不知道为什么
	<-r

	return
}

func (p *Etcd) RegisterService(server *Server) (leaseId string, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	bs, _ := json.Marshal(server)
	leaseOption, _ := p.cli.Grant(context.TODO(), p.ttl)
	_, err = p.cli.Put(ctx, root+server.Id, string(bs), clientv3.WithLease(leaseOption.ID))
	if err != nil {
		return
	}
	leaseId = strconv.FormatInt(int64(leaseOption.ID), 10)
	return
}

func (p *Etcd) UnRegisterService(serverId string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = p.cli.Delete(ctx, root+serverId)
	return
}

func (p *Etcd) init() (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   p.etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	p.cli = cli
	return
}

func NewEtcd(etcdEndpoints []string, ttl int64) Discoverer {
	etcd := &Etcd{
		ttl:           ttl,
		etcdEndpoints: etcdEndpoints,
	}
	err := etcd.init()
	if err != nil {
		panic(err)
	}
	return etcd
}
