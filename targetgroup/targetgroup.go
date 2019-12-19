package targetgroup

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/vorozhko/loadbalancer/config"
	"gitlab.com/vorozhko/loadbalancer/roundrobin"
)

//TargetGroup - is a loadbalancer target group instance
type TargetGroup struct {
	toPort             int
	fromPort           int
	path               string
	instances          []string
	selection          *roundrobin.RoundRobin
	instanceNotHealthy map[string]bool
}

func InitTargetGroup(target config.ConfigTargetGroup) *TargetGroup {
	tg := TargetGroup{}
	tg.fromPort = target.GetFromPort()
	tg.toPort = target.GetToPort()
	tg.path = target.GetPath()
	tg.instances = target.GetInstances()
	//todo: create custom type for instances or instance health
	tg.instanceNotHealthy = make(map[string]bool, len(target.GetInstances()))
	go func() {
		//health checker
		for {
			if len(tg.instanceNotHealthy) > 0 {
				for instance, status := range tg.instanceNotHealthy {
					if status == false {
						continue
					}
					upstream := fmt.Sprintf("%s:%d", instance, target.GetToPort())
					res, err := http.Get(upstream)
					if err == nil && res != nil {
						tg.instanceNotHealthy[instance] = false
						fmt.Printf("%s makred as healty\n", instance)
					}
				}
			}
			time.Sleep(25 * time.Second)
		}
	}()
	return &tg
}

func (tg *TargetGroup) getUpstream() (string, error) {
	//todo: replace default Round Robin with algorithm selection
	instances := make([]string, 0)
	for _, inst := range tg.instances {
		if tg.instanceNotHealthy[inst] == false {
			instances = append(instances, inst)
		}
	}
	if len(instances) == 0 {
		return "", fmt.Errorf("No healthy hosts found\n")
	}
	return tg.selection.GetNextUpstreamIndex(instances), nil
}

//SetUpstreamSelection - set an algorithm for select of upstream server
func (tg *TargetGroup) SetUpstreamSelection(roundRobin *roundrobin.RoundRobin) {
	tg.selection = roundRobin
}

func (tg *TargetGroup) GetPath() string {
	return tg.path
}

//ServeHTTP - implement http.Handler.ServeHTTP for server mux
func (tg *TargetGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var client http.Client

	upstreamHost, err := tg.getUpstream()
	if err != nil {
		fmt.Print(err)
		return
	}
	upstream := fmt.Sprintf("%s:%d", upstreamHost, tg.toPort)
	upstreamRes, err := client.Get(upstream + req.RequestURI)

	if err != nil {
		//todo: replace with Logger middleware
		fmt.Printf("%s", err)
		tg.instanceNotHealthy[upstreamHost] = true
		fmt.Printf("%s makred unhealty\n", upstreamHost)
		return
	}
	if upstreamRes == nil {
		//todo: replace with Logger middleware
		fmt.Printf("Empty response from server")
		tg.instanceNotHealthy[upstreamHost] = true
		fmt.Printf("%s makred unhealty\n", upstreamHost)
		return
	}
	defer upstreamRes.Body.Close()

	for hk := range upstreamRes.Header {
		w.Header().Add(hk, upstreamRes.Header.Get(hk))
	}

	body, err := ioutil.ReadAll(upstreamRes.Body)
	if err != nil {
		//todo: replace with Logger middleware
		fmt.Printf("%s", err)
		//todo: replace with Error middleware which will print standard error message to the user
		fmt.Fprintf(w, "Internal Server Error")
		//todo: display debug information only when debug is enabled
		fmt.Fprintf(w, "%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

	}
	//if everything OK, use the status from upstream
	w.WriteHeader(upstreamRes.StatusCode)
	fmt.Fprintf(w, "%s", body)
}
