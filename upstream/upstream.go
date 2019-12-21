package upstream

//RoundRobin - Round robin load balancer target selection
type RoundRobin struct {
	currentInstance int
	maxInstances    int
}

//GetUpstream - return next upstream host based on round robin algorithm
func (rr *RoundRobin) GetNextUpstreamIndex(instances []string) string {
	max := len(instances)
	if rr.currentInstance >= max {
		rr.currentInstance = 0
	}
	instance := rr.currentInstance
	rr.currentInstance++
	return instances[instance]
}
