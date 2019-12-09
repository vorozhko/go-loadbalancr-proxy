# Loadbalancer in Golang

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

### Request to backend server
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

## Roadmap

### Milestone 1
* Design go code modular structure
* Refactor implemented features into modular structure

### Milestone 2
* YAML: Load multi listeners and backend servers endpoints
* Store servers status: up/down
* Exlcude down servers
* Retry on error
* Request parsing: client ip, host, port, path

## Releases

### Release 0.0.1
* Added: YAML config file to manage LB settings
* Added: Load single listener and multiple backend servers
* Added: Least connections algorithm for servers selection

### Links
[Load balancing wiki page](https://en.wikipedia.org/wiki/Load_balancing_(computing))