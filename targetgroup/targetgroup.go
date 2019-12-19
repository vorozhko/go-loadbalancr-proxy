package targetgroup

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/vorozhko/loadbalancer/config"
	"gitlab.com/vorozhko/loadbalancer/roundrobin"
)

//TargetGroup - is a loadbalancer target group instance
type TargetGroup struct {
	toPort    int
	fromPort  int
	path      string
	instances []string
	selection *roundrobin.RoundRobin
}

func InitTargetGroup(target config.ConfigTargetGroup) *TargetGroup {
	tg := TargetGroup{}
	tg.fromPort = target.GetFromPort()
	tg.toPort = target.GetToPort()
	tg.path = target.GetPath()
	tg.instances = target.GetInstances()
	return &tg
}

func (tg *TargetGroup) getUpstream() string {
	//todo: replace default Round Robin with algorithm selection

	nextInstance := tg.selection.GetNextUpstreamIndex(len(tg.instances))
	return tg.instances[nextInstance]
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

	upstream := fmt.Sprintf("%s:%d", tg.getUpstream(), tg.toPort)
	upstreamRes, err := client.Get(upstream + req.RequestURI)

	if err != nil {
		//todo: replace with Logger middleware
		fmt.Printf("%s", err)
	}
	if upstreamRes == nil {
		//todo: replace with Logger middleware
		fmt.Printf("Empty response from server")
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
