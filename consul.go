package main

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

func initConsul() (*consul.Client, error) {
	consulConfig := consul.DefaultConfig()
	return consul.NewClient(consulConfig)
}

// Register a service with consul local agent
func register(client *consul.Client, name, address string, port int) error {
	reg := &consul.AgentServiceRegistration{
		ID:      name,
		Name:    name,
		Address: address,
		Port:    port,
	}
	return client.Agent().ServiceRegister(reg)
}

// DeRegister a service with consul local agent
func unregister(client *consul.Client, id string) error {
	return client.Agent().ServiceDeregister(id)
}

// Service return a service
func service(client *consul.Client, service, tag string) (string, error) {
	passingOnly := true
	addrs, _, err := client.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return "", fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", addrs[0].Service.Address, addrs[0].Service.Port), nil
}
