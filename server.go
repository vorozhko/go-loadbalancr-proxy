package loadbalancer

import (
	"fmt"
	"net/http"

	config "gitlab.com/vorozhko/loadbalancer/config"
	"gitlab.com/vorozhko/loadbalancer/loadbalancer"
	"gitlab.com/vorozhko/loadbalancer/roundrobin"
)

//Server - is server instance
type Server struct {
	config *config.Config
}

//Start - main entry point
func (app *Server) Start(configFile string) (err error) {
	app.config, err = config.InitConfig(configFile)
	if err != nil {
		return err
	}
	if len(app.config.Listeners) == 0 {
		return fmt.Errorf("No open ports defined for listeners")
	}
	finish := make(chan bool)
	err = app.startListeners()
	if err != nil {
		return err
	}
	<-finish
	return nil
}

func (app *Server) startListeners() error {
	for index, port := range app.config.Listeners {
		//todo: separate app from Load balancer object
		lb := loadbalancer.Loadbalancer{}
		if len(app.config.Targets) > index {
			//todo: replace default Round Robin with algorithm selection
			roundRobin := roundrobin.InitRoundRobin(app.config.Targets[index])
			lb.SetUpstreamSelection(roundRobin)
		}
		//todo: replace ListenAndServe with multi port listen
		listenPort := fmt.Sprintf(":%d", port)

		go func() {
			mux := http.NewServeMux()
			//http.HandleFunc("/", lb.httpProxy)
			reqHandler := &lb
			mux.Handle("/", reqHandler)
			fmt.Printf("Listening on port: %s\n", listenPort)
			http.ListenAndServe(listenPort, mux)
		}()
	}
	return nil
}
