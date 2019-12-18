package roundrobin

//RoundRobin - Round robin load balancer target selection
type RoundRobin struct {
	currentInstance int
	maxInstances    int
}

//GetUpstream - return next upstream host based on round robin algorithm
func (rr *RoundRobin) GetNextUpstreamIndex(max int) int {
	if rr.currentInstance >= max {
		rr.currentInstance = 0
	}
	instance := rr.currentInstance
	rr.currentInstance++
	return instance
}
