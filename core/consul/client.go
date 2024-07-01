package consul

import (
	"encoding/json"
	"fmt"
	"github.com/sniperCore/core/log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

type Consul struct {
	Conn      *api.Client
	LastIndex *api.QueryOptions
}

func NewConsul(ConConfig *ConsulConfig) (*Consul, error) {
	cfg := api.DefaultConfig()
	cfg.Address = strings.Join([]string{ConConfig.Host, ":", strconv.Itoa(ConConfig.Port)}, "")
	conn, err := api.NewClient(cfg)
	if err != nil {
		log.Error("consul new client error, err = ", err)
	}
	lastIndex := &api.QueryOptions{
		WaitIndex: 0,
		WaitTime:  10 * time.Second,
	}
	return &Consul{
		Conn:      conn,
		LastIndex: lastIndex,
	}, nil
}

func (consul *Consul) Register() {
	if ConConfig.WorkerType == "ms" {
		ConConfig.WorkerType = ConConfig.ServiceName
	}
	if ConConfig.ServiceHost == "" || ConConfig.Host == "0.0.0.0" {
		ConConfig.ServiceHost = consul.GetLocalIP()
	}

	check := &api.AgentServiceCheck{}
	if ConConfig.Check.Enabled {
		check.DeregisterCriticalServiceAfter = strings.Join([]string{strconv.Itoa(ConConfig.Check.RemoveTime), "m"}, "")
		check.Timeout = strings.Join([]string{strconv.Itoa(ConConfig.Check.Timeout), "s"}, "")
		check.Interval = strings.Join([]string{strconv.Itoa(ConConfig.Check.Interval), "s"}, "")
		check.TCP = strings.Join([]string{ConConfig.ServiceHost, ":", strconv.Itoa(ConConfig.ServicePort)}, "")

	}

	checkInfo, _ := json.Marshal(check)
	log.Info("consul register check: ", string(checkInfo))
	service := &api.AgentServiceRegistration{
		ID:      strings.Join([]string{ConConfig.WorkerType, ":", ConConfig.ServiceHost + ":" + fmt.Sprintf("%d", ConConfig.ServicePort)}, ""),
		Name:    ConConfig.WorkerType,
		Address: ConConfig.ServiceHost,
		Port:    ConConfig.ServicePort,
		Tags:    ConConfig.Tag,
		Check:   check,
		Meta:    ConConfig.WorkerMeta,
	}

	serviceInfo, _ := json.Marshal(service)
	log.Info("consul register service: ", string(serviceInfo))
	if err := consul.Conn.Agent().ServiceRegister(service); err != nil {
		time.Sleep(2 * time.Second)
		if err = consul.Conn.Agent().ServiceRegister(service); err != nil {
			log.Error("consul register service error, check: ", check, "service: ", service, "msg: ", err)
		}
	} else {
		log.Info("consul register service success ", string(serviceInfo))
	}
}

func (consul *Consul) DisRegister() {
	if ConConfig.Check.Enabled == false {
		return
	}
	id := strings.Join([]string{ConConfig.WorkerType, ":", ConConfig.ServiceHost, ":", strconv.Itoa(ConConfig.ServicePort)}, "")
	consul.Conn.KV().Delete(id, nil)
	if err := consul.Conn.Agent().ServiceDeregister(id); err != nil {
		log.Error("disregister service error, id: ", id, "msg: ", err)
	} else {
		log.Info("disregister service success: ", id)
	}
}

func (consul *Consul) HealthyCheck(serverName string) (api.HealthChecks, error) {
	healthChecks, _, err := consul.Conn.Health().Checks(serverName, nil)
	if err != nil {
		return nil, err
	}

	return healthChecks, nil
}

func (consul *Consul) GetLocalIP() string {
	self, err := consul.Conn.Agent().Self()
	if err != nil {
		log.Error("consul agent self error, msg: ", err)
	}

	ip := self["Member"]["Addr"].(string)
	if ip == "" {
		ip = self["DebugConfig"]["BindAddr"].(string)
	}

	log.Info("consul bind addr: ", ip)

	return ip
}
