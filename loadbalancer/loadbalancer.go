package loadbalancer

import (
	"fmt"
	"net/http"

	config "gitlab.com/vorozhko/loadbalancer/config"
	"gitlab.com/vorozhko/loadbalancer/roundrobin"
	"gitlab.com/vorozhko/loadbalancer/targetgroup"
)

//LoadBalancer - loadbalancer app instance
type LoadBalancer struct {
	config *config.Config
}

//Start - main entry point
func (app *LoadBalancer) Start(configFile string) (err error) {
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

func (app *LoadBalancer) startListeners() error {
	if len(app.config.Listeners) == 0 {
		return fmt.Errorf("No open ports defined for listeners")
	}

	targetGroupsList := make(map[int][]*targetgroup.TargetGroup, len(app.config.Targets))
	for _, target := range app.config.Targets {
		//Multiple target groups support split by URI.Path
		tg := targetgroup.InitTargetGroup(target)
		roundRobin := roundrobin.RoundRobin{}
		tg.SetUpstreamSelection(&roundRobin)
		targetGroupsList[target.FromPort] = append(targetGroupsList[target.FromPort], tg)
	}

	//open listener ports
	for _, port := range app.config.Listeners {
		listenPort := fmt.Sprintf(":%d", port)
		//Support of multiple target groups per inbound port
		targets, ok := targetGroupsList[port]

		//each listen port has it's own go routine
		go func() {
			mux := http.NewServeMux()
			if ok == true {
				//Split target groups handlers by path.
				//Path should be uniq per target group
				for _, tg := range targets {
					path := tg.GetPath()
					if path == "" {
						path = "/"
					}
					//one handler per port and path
					mux.Handle(path, tg)
				}
			}
			fmt.Printf("Listening on port: %s\n", listenPort)
			http.ListenAndServe(listenPort, mux)
		}()

	}
	return nil
}
