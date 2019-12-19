package main

import (
	"log"

	loadbalancer "gitlab.com/vorozhko/loadbalancer/loadbalancer"
)

func main() {
	server := loadbalancer.LoadBalancer{}
	err := server.Start("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
