package main

import "gitlab.com/vorozhko/loadbalancer"

import "log"

func main() {
	server := loadbalancer.Server{}
	err := server.Start("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
