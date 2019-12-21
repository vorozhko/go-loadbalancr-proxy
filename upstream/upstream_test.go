package upstream

import (
	"testing"
)

func TestRoundRobin(t *testing.T) {
	instances := []string{"localhost0", "localhost1", "localhost2"}
	rr := RoundRobin{}
	upstream := rr.GetNextUpstream(instances)
	if upstream != "localhost0" {
		t.Errorf("Got %s upstream, but expected %s", upstream, "localhost0")
	}
}
