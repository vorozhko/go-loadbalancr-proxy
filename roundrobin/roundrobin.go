package roundrobin

import (
	"fmt"

	"gitlab.com/vorozhko/loadbalancer/config"
)

//RoundRobin - Round robin load balancer target selection
type RoundRobin struct {
	Targets            config.TargetGroup
	roundRobinInstance int
}

//InitRoundRobin - construct RoundRobin with data
func InitRoundRobin(tg config.TargetGroup) *RoundRobin {
	var rr RoundRobin
	rr.roundRobinInstance = 0
	rr.Targets = tg
	return &rr
}

//GetUpstream - return next upstream host based on round robin algorithm
func (rr *RoundRobin) GetUpstream() string {
	upstream := ""
	for index, upstreamInstance := range rr.Targets.Instances {
		if rr.roundRobinInstance == index {
			upstream = fmt.Sprintf("%s:%d", upstreamInstance, rr.Targets.ToPort)
			rr.roundRobinInstance++
			if rr.roundRobinInstance == len(rr.Targets.Instances) {
				rr.roundRobinInstance = 0
			}
			break
		}
	}
	return upstream
}
