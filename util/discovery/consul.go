package discovery

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"sync"
	"errors"
)

type Consul struct {
	stdServers map[string]*Server
	stdCli     *api.Client
	lock       sync.Mutex
}

func (p *Consul) GetServers() (servers map[string]*Server, err error) {
	if p.stdServers == nil {
		err = errors.New("not init")
		return
	}

	return p.stdServers, nil
}

func (p *Consul) WatchServer(ServersChanged) {

}
func (p *Consul) init() (err error) {
	consulConf := api.DefaultConfig()
	consulConf.Address = consulAddr
	cli, err := api.NewClient(consulConf)
	if err != nil {
		panic(err)
	}
	p.stdCli = cli
	p.watchService("services", "checks")
	return
}

func (p *Consul) UpdateServerTTL(server *Server, status string) (err error) {
	err = p.stdCli.Agent().UpdateTTL("service:"+server.Id, "", status)
	return
}

func (p *Consul) RegisterService(server *Server, ttl string) (err error) {
	s := &api.AgentServiceRegistration{
		ID:      server.Id,
		Address: server.Address,
		Name:    server.Name,
		Port:    server.Port,
		Tags:    []string{},
		Check: &api.AgentServiceCheck{
			TTL: ttl,
		},
	}

	err = p.stdCli.Agent().ServiceRegister(s)
	return
}

func (p *Consul) UnRegisterService(serverId string) (err error) {
	err = p.stdCli.Agent().ServiceDeregister(serverId)
	return
}

func (p *Consul) updateServicesById(name []string, passOnly bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.stdServers = map[string]*Server{}
	for _, name := range name {
		serverData, _, e := p.stdCli.Health().Service(name, "", passOnly,
			&api.QueryOptions{})
		if e != nil {
			return
		}

		for _, entry := range serverData {
			p.stdServers[entry.Service.ID] = &Server{
				Address: entry.Service.Address,
				Port:    entry.Service.Port,
				Name:    entry.Service.Service,
				Id:      entry.Service.ID,
			}
		}
	}

	return
}

func (p *Consul) watchService(types ...string) {
	isReady := false
	readyChan := make(chan int)
	for _, t := range types {
		go func(t string) {
			plan, e := watch.Parse(map[string]interface{}{
				"type": t,
			})
			if e != nil {
				return
			}
			plan.Handler = func(id uint64, raw interface{}) {
				switch data := raw.(type) {
				case map[string][]string:
					names := []string{}
					for name := range data {
						names = append(names, name)
					}
					p.updateServicesById(names, false)
				case []*api.HealthCheck:
					names := []string{}
					for _, name := range data {
						names = append(names, name.ServiceID)
					}
					p.updateServicesById(names, false)
				}
				if !isReady {
					isReady = true
					close(readyChan)
				}
			}
			plan.Run(consulAddr)
		}(t)
	}
	select {
	case <-readyChan:
	}
	return
}

var consulAddr = "127.0.0.1:8500"
