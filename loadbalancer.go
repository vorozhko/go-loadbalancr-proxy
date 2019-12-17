package loadbalancer

import (
	"fmt"
	"io/ioutil"
	"net/http"

	config "gitlab.com/vorozhko/loadbalancer/config"
	"gitlab.com/vorozhko/loadbalancer/roundrobin"
)

//Loadbalancer - represent an app instance
type Loadbalancer struct {
	config     *config.Config
	roundRobin *roundrobin.RoundRobin
}

//Start - main entry point
func (lb *Loadbalancer) Start(configFile string) (err error) {
	lb.config, err = config.InitConfig(configFile)
	if err != nil {
		return err
	}
	//todo: replace default Round Robin with algorithm selection
	lb.roundRobin = roundrobin.InitRoundRobin(lb.config.Targets[0])
	//todo: replace ListenAndServe with multi port listen
	http.HandleFunc("/", lb.httpProxy)
	listenPort := fmt.Sprintf(":%d", lb.config.Listeners[0])
	fmt.Printf("Listening on port: %s", listenPort)
	return http.ListenAndServe(listenPort, nil)
}

func (lb *Loadbalancer) getUpstream() string {
	//todo: replace default Round Robin with algorithm selection
	return lb.roundRobin.GetUpstream()
}

func (lb *Loadbalancer) httpProxy(w http.ResponseWriter, req *http.Request) {
	var client http.Client

	upstream := lb.getUpstream()
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
