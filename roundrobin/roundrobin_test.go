package roundrobin

import (
	"fmt"
	"testing"

	config "gitlab.com/vorozhko/loadbalancer/config"
)

func TestRoundRobin(t *testing.T) {
	var targets config.TargetGroup
	targets.Instances = []string{"localhost0", "localhost1", "localhost2"}
	targets.FromPort = 80
	targets.ToPort = 8080
	rr := InitRoundRobin(targets)
	if len(rr.Targets.Instances) != 3 {
		t.Errorf("Got %d instances, but expected 3", len(rr.Targets.Instances))
	}

	for i := 0; i < len(rr.Targets.Instances); i++ {
		expectedUpstream := fmt.Sprintf("localhost%d:%d", i, rr.Targets.ToPort)
		upstream := rr.GetUpstream()
		if expectedUpstream != upstream {
			t.Errorf("Got %s upstream, but expected %s upstream", upstream, expectedUpstream)
		}
	}
}
