package main

import "gitlab.com/vorozhko/loadbalancer"

func main() {
	var lb loadbalancer.Loadbalancer
	lb.Start("config.yaml")
}
