package loadbalancer

import (
	"fmt"
	"net/http"

	config "gitlab.com/vorozhko/loadbalancer/config"
	"gitlab.com/vorozhko/loadbalancer/roundrobin"
	"gitlab.com/vorozhko/loadbalancer/targetgroup"
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

	finish := make(chan bool)
	err = app.startListeners()
	if err != nil {
		return err
	}
	<-finish
	return nil
}

func (app *Server) startListeners() error {
	if len(app.config.Listeners) == 0 {
		return fmt.Errorf("No open ports defined for listeners")
	}

	targetGroupsList := make(map[int][]*targetgroup.TargetGroup, len(app.config.Targets))
	for _, target := range app.config.Targets {
		//Multiple target groups support split by URI.Path
		tg := targetgroup.InitTargetGroup(target.FromPort, target.ToPort, target.Instances, target.Path)
		roundRobin := roundrobin.RoundRobin{}
		tg.SetUpstreamSelection(&roundRobin)
		targetGroupsList[target.FromPort] = append(targetGroupsList[target.FromPort], tg)
	}

	//open listener ports
	for _, port := range app.config.Listeners {
		listenPort := fmt.Sprintf(":%d", port)
		targets, ok := targetGroupsList[port]
		go func() {
			mux := http.NewServeMux()
			if ok == true {
				for _, tg := range targets {
					path := tg.GetPath()
					if path == "" {
						path = "/"
					}
					mux.Handle(path, tg)
				}
			}
			fmt.Printf("Listening on port: %s\n", listenPort)
			http.ListenAndServe(listenPort, mux)
		}()
	}
	return nil
}
