package main

import "gitlab.com/vorozhko/loadbalancer"

import "log"

func main() {
	var lb loadbalancer.Loadbalancer
	err := lb.Start("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
