package targetgroup

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/vorozhko/loadbalancer/config"
	healthcheck "gitlab.com/vorozhko/loadbalancer/healthcheck"
	"gitlab.com/vorozhko/loadbalancer/upstream"
)

//TargetGroup - is a loadbalancer target group instance
type TargetGroup struct {
	toPort          int
	fromPort        int
	path            string
	instances       []string
	isStickySession bool
	instanceHealth  healthcheck.InstanceHealth
	leastConnect    upstream.LeastConnect
	roundrobin      upstream.RoundRobin
}

func InitTargetGroup(target config.ConfigTargetGroup) *TargetGroup {
	tg := TargetGroup{}
	tg.fromPort = target.GetFromPort()
	tg.toPort = target.GetToPort()
	tg.path = target.GetPath()
	tg.instances = target.GetInstances()
	tg.isStickySession = target.GetStickySession()
	tg.leastConnect = upstream.InitLeastConnect(tg.instances)
	tg.instanceHealth = healthcheck.InitInstanceHealth(tg.instances, tg.toPort)
	return &tg
}

func (tg *TargetGroup) getUpstream(w http.ResponseWriter, req *http.Request) (string, error) {
	//todo: here several algorithms of upstream selection could be implemented
	//the best one should win
	//Algorithms example: round robin, least connect, sticky session could determine host selection
	//Priority must go to sticky session if set, then least connect and round robin.

	//least connect
	upstreamHost, err := tg.getUpstreamLeastConnect()
	if err != nil {
		return "", err
	}

	if upstreamHost == "" {
		//round robin
		upstreamHost, err = tg.getUpstreamRoundRobin()
		if err != nil {
			return "", err
		}
		fmt.Printf("Round robin host %s\n", upstreamHost)
	}

	//sticky session
	//todo: rename sticky with better name
	if tg.isStickySession == true {
		stickyCookie, err := req.Cookie("sticky")
		if err == nil {
			fmt.Printf("Sticky host %s\n", stickyCookie.Value)
			return stickyCookie.Value, nil
		}

		//if stickyness enabled then set new cookie
		newCookie := http.Cookie{Name: "sticky", Value: upstreamHost}
		http.SetCookie(w, &newCookie)
	}
	return upstreamHost, nil
}

//getUpstreamRoundRobin - return next upstream by round robin algorithm
func (tg *TargetGroup) getUpstreamRoundRobin() (string, error) {
	instances := tg.instanceHealth.GetHealthyInstances()
	if len(instances) == 0 {
		return "", fmt.Errorf("No healthy hosts found\n")
	}
	return tg.roundrobin.GetNextUpstream(instances), nil
}

//getUpstreamLeastConnect - return upstream host with small number of concurrent connections
func (tg *TargetGroup) getUpstreamLeastConnect() (string, error) {
	instances := tg.instanceHealth.GetHealthyInstances()
	if len(instances) == 0 {
		return "", fmt.Errorf("No healthy hosts found\n")
	}
	return tg.leastConnect.GetNextUpstream(instances), nil
}

func (tg *TargetGroup) GetPath() string {
	return tg.path
}

func (tg *TargetGroup) GetInstances() []string {
	return tg.instances
}

func (tg *TargetGroup) GetToPort() int {
	return tg.toPort
}

//ServeHTTP - implement http.Handler.ServeHTTP for server mux
func (tg *TargetGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var client http.Client

	//get upstream
	upstreamHost, err := tg.getUpstream(w, req)
	if err != nil {
		fmt.Print(err)
		return
	}
	tg.leastConnect.Connect(upstreamHost)
	upstream := fmt.Sprintf("%s:%d", upstreamHost, tg.toPort)
	upstreamRes, err := client.Get(upstream + req.RequestURI)
	tg.leastConnect.Disconnect(upstreamHost)

	//process response
	if err != nil {
		//todo: replace with Logger middleware
		fmt.Printf("%s", err)
		tg.instanceHealth.SetHealth(false, upstreamHost, tg.toPort)
		fmt.Printf("%s makred unhealty\n", upstreamHost)
		return
	}
	if upstreamRes == nil {
		//todo: replace with Logger middleware
		fmt.Printf("Empty response from server")
		tg.instanceHealth.SetHealth(false, upstreamHost, tg.toPort)
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
