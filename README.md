# Loadbalancer in Golang

## Usage
Describe load balancing ports and targets in config.yaml. See [example config file](example/config.yaml).

```
package main

import (
	"log"

	loadbalancer "gitlab.com/vorozhko/loadbalancer/loadbalancer"
)

func main() {
	server := loadbalancer.LoadBalancer{}
	err := server.Start("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}

```
## Roadmap

### Milestone 0 ([Release 0.0.1](https://gitlab.com/vorozhko/loadbalancer/-/tags/v0.0.1))
* Prototype first load balancer version

### Milestone 1 ([Release 0.0.2](https://gitlab.com/vorozhko/loadbalancer/-/tags/v0.0.2))
* Design code structure
* Impleement round roubin load balancing

### Milestone 2 ([Release 0.0.3](https://gitlab.com/vorozhko/loadbalancer/-/tags/v0.0.3))
* YAML: Load multi listeners and backend servers endpoints
* Store servers status: up/down
* Exlcude down servers
* Least connections algorithm
* Content based routing
* Sticky sessions support

### Milestone 3 ([Release 0.0.4](https://gitlab.com/vorozhko/loadbalancer/-/tags/v0.0.4))
* Refactor upstream selection
* Refactor health checker

### Milestone 4
* Retry on error
* Request parsing: client ip, host, port, path
* HTTP compression
* Error handling - show custom error message to user
* Content filtering - modify request content by some rules


## Features list

### Requests handling
* Request parsing: client ip, host, port, path
* TLS termination
* DDOS protection

### Backend server selection
* Exclude down servers
* Check if persistent session enabled
* Content based routing - /video and /images go to different servers
* Assymetric load - if manual weight for servers is set
* Select best of 2 random servers
* Least connections
* Backend server load - get reported by application response headers if present
* Response time

### Upstream requests
* HTTP compression
* Error handling - show custom error message to user
* Content filtering - modify request content by some rules
* Retry on error
* 2nd backend server on error - use second best backend server from server selection step

### Persistent data
* Servers status: up/down
* Servers load
* Servers connections
* Servers resposne time: median, 95pt, 99pt for past 5 minutes
* Monitoring: response time, network IO, status codes, errors rate
* Storage VaultDB

### Loadbalancer manager
* YAML: config file to manage LB settings
* YAML: Load multi listeners and backend servers endpoints
* YAML: Load persistent settings
* YAML: Load manual weights for backend servers
* Re-read config file by SIGNAL

### Infrastructure
* Add Docker support
* Run in distiributed HA pairs
* Share state between HA pairs
* Add Kubernetes support

## Releases

### Release 0.0.1
* Added: YAML config file to manage LB settings
* Added: Load single listener and multiple backend servers
* Added: Least connections algorithm for servers selection

### Release 0.0.2
* New code structure
* Impleemented round roubin load balancing

### Release 0.0.3
* Added support for multiple target groups on the same port split by request Path
* Added content based routing. See ```Path``` in YAML file.
* Added health check monitor 
* Exlcude down servers from upstream selection
* Implemented Least connections algorithm
* Added support for sticky sessions

### Release 0.0.4
* Refactoring of upstream selection and health checker

### Links
* [Load balancing wiki page](https://en.wikipedia.org/wiki/Load_balancing_(computing))