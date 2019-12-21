package upstream

import "fmt"

//RoundRobin - Round robin load balancer target selection
type RoundRobin struct {
	currentInstance int
	maxInstances    int
}

type LeastConnect struct {
	instanceConnections map[string]int
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
	lc.instanceConnections = make(map[string]int, len(instances))
	return lc
}

func (lc *LeastConnect) GetNextUpstream(instances []string) string {
	leastConnectInstance := instances[0]
	leastConnectRequests := lc.instanceConnections[leastConnectInstance]
	for _, inst := range instances {
		if lc.instanceConnections[inst] < leastConnectRequests {
			leastConnectRequests = lc.instanceConnections[inst]
			leastConnectInstance = inst
		}
	}
	//todo: remove debug information
	fmt.Printf("Least connect host %s with %d connections\n", leastConnectInstance, leastConnectRequests)
	return leastConnectInstance
}

func (lc *LeastConnect) Connect(host string) {
	lc.instanceConnections[host]++
}

func (lc *LeastConnect) Disconnect(host string) {
	lc.instanceConnections[host]--
}
