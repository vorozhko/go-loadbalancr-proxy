package upstream

import "fmt"

//RoundRobin - Round robin load balancer target selection
type RoundRobin struct {
	currentInstance int
	maxInstances    int
}

type LeastConnect struct {
	connections map[string]int
}

//GetUpstream - return next upstream host based on round robin algorithm
func (rr *RoundRobin) GetNextUpstream(instances []string) string {
	max := len(instances)
	if rr.currentInstance >= max {
		rr.currentInstance = 0
	}
	instance := rr.currentInstance
	rr.currentInstance++
	return instances[instance]
}

func InitLeastConnect(instances []string) LeastConnect {
	lc := LeastConnect{}
	lc.connections = make(map[string]int, len(instances))
	return lc
}

func (lc *LeastConnect) GetNextUpstream(instances []string) string {
	host := instances[0]
	connectedTimes := lc.connections[host]
	for _, inst := range instances {
		if lc.connections[inst] < connectedTimes {
			connectedTimes = lc.connections[inst]
			host = inst
		}
	}
	//todo: remove debug information
	fmt.Printf("Least connect host %s with %d connections\n", host, connectedTimes)
	return host
}

func (lc *LeastConnect) Connect(host string) {
	lc.connections[host]++
}

func (lc *LeastConnect) Disconnect(host string) {
	lc.connections[host]--
}
