package service

import (
	"fmt"
	"net"
	log "github.com/sirupsen/logrus"
	consul "github.com/hashicorp/consul/api"
)

func (s *Service) registerConsul() {
	c, err := consul.NewClient(&consul.Config{
		Address: "consul-consul-server.default.svc.cluster.local:8500",
	})
	if err != nil {
		log.Fatal("Failed to connect to Consul: ", err)
	}
	s.ConsulAgent = c.Agent()
	s.ConsulHealth = c.Health()

	ifaces, err := net.Interfaces()
	// handle err
	var ip net.IP
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			break; //get only the first ip?
		}
	}

	// To register multiple services with the same Name,
	// we must define directly the ID and the specific Address. 
	// To get a unique ID we use the associated IP as a suffix.
	serviceDef := &consul.AgentServiceRegistration{
		ID: fmt.Sprintf("%s-%s", s.Name, ip.String()),
		Name: s.Name,
		Address: ip.String(),
		Port: s.Port,
		Check: &consul.AgentServiceCheck{
			Interval: "10s",
			HTTP: fmt.Sprintf("http://%s:%d/health", ip.String(), s.Port),
		},
	}

	log.Printf("try to register the service %s %s", serviceDef.ID, serviceDef.Address)
	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		log.Fatal("Failed to register Service: ", err)
	}
}

func (s *Service) countConsulServices() int {
	results, _, err := s.ConsulHealth.Service(s.Name, "", true, nil) 
	if err != nil {
		log.Fatal("Consul server failed to get service: ", err)
	}

	counter := 0
	for _, result := range results {
		svc := result.Service
		log.Printf("service %s with address %s found", svc.Service, svc.Address)
		counter++
	}
	return counter
}