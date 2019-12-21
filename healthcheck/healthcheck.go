package healthcheck

import (
	"fmt"
	"net/http"
	"time"
)

type InstanceHealth struct {
	instances map[string]Instance
}

type Instance struct {
	host    string
	port    int
	healthy bool
}

func InitInstanceHealth(instances []string, toPort int) InstanceHealth {
	ih := InstanceHealth{}
	ih.SetInstances(instances, toPort)
	//run background health checker to recover status of unhealthy hosts
	go ih.HealthChecker(25)
	return ih
}

func (ih *InstanceHealth) SetInstances(instances []string, toPort int) {
	ih.instances = make(map[string]Instance)
	for _, host := range instances {
		instance := Instance{}
		instance.host = host
		instance.healthy = true
		instance.port = toPort
		ih.instances[instance.host] = instance
	}
}

//healthChecker - periodically check bad healthy hosts for recover
func (ih *InstanceHealth) HealthChecker(checkInterval time.Duration) {
	for {
		if len(ih.instances) > 0 {
			for host, instance := range ih.instances {
				if instance.healthy == true {
					continue
				}
				upstream := fmt.Sprintf("%s:%d", instance.host, instance.port)
				res, err := http.Get(upstream)
				if err == nil && res != nil {
					instance.healthy = true
					ih.instances[host] = instance
					fmt.Printf("%s makred as healty\n", instance.host)
				}
			}
		}
		time.Sleep(checkInterval)
	}
}

func (ih *InstanceHealth) GetHealthyInstances() []string {
	instances := make([]string, 0)
	for _, instance := range ih.instances {
		if instance.healthy == true {
			instances = append(instances, instance.host)
		}
	}
	return instances
}

func (ih *InstanceHealth) SetHealth(healthy bool, host string) {
	inst := ih.instances[host]
	inst.healthy = healthy
	ih.instances[host] = inst
}
