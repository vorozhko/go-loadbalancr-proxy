package loadbalancer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Loadbalancer struct {
	config     Config
	connection int
	targets    []TargetInstance
}

type Config struct {
	Listen  string
	Targets []string
}

type TargetInstance struct {
	host        string
	connections int
}

func (t *TargetInstance) StartConnection() {
	t.connections++
}
func (t *TargetInstance) EndConnection() {
	t.connections--
}
func (t *TargetInstance) GetHost() string {
	return t.host
}
func (t *TargetInstance) GetConnections() int {
	return t.connections
}

func (lb *Loadbalancer) Start(configFile string) {
	err := lb.loadConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	lb.initTargets()
	http.HandleFunc("/", lb.httpProxy)
	log.Fatal(http.ListenAndServe(":"+lb.config.Listen, nil))
}

func (lb *Loadbalancer) loadConfig(filename string) error {
	yml, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var c Config
	err = yaml.Unmarshal(yml, &c)
	if err != nil {
		return err
	}
	lb.config = c
	return nil
}

func (lb *Loadbalancer) initTargets() {
	if len(lb.config.Targets) == 0 {
		return
	}
	lb.targets = make([]TargetInstance, len(lb.config.Targets))
	for index, target := range lb.config.Targets {
		lb.targets[index].host = target
	}
}

func (lb *Loadbalancer) httpProxy(w http.ResponseWriter, req *http.Request) {
	lb.connection++
	fmt.Printf("Total connections: %d\n", lb.connection)
	var client http.Client
	target := lb.getTarget()
	target.StartConnection()
	fmt.Printf("server %s, connections %d\n", target.GetHost(), target.GetConnections())
	response, err := client.Get(target.GetHost() + req.RequestURI)
	target.EndConnection()
	fmt.Printf("server %s, connections %d\n", target.GetHost(), target.GetConnections())
	if err != nil {
		log.Println(err)
		return
	}
	resheader := w.Header()
	for hk := range response.Header {
		resheader.Add(hk, response.Header.Get(hk))
	}
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
	defer response.Body.Close()

	w.WriteHeader(response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
	fmt.Fprintf(w, "%s", body)
}

func (lb *Loadbalancer) getTarget() *TargetInstance {
	minIndex := 0
	for index, t := range lb.targets {
		if t.GetConnections() < lb.targets[minIndex].GetConnections() {
			minIndex = index
		}
	}

	return &lb.targets[minIndex]
}
