# HTTP proxy todo

## Features list
[Load balancing](https://en.wikipedia.org/wiki/Load_balancing_(computing))

## Milestone 0 [completed]

* Add target hosts and listen ports configuration thorugh YAML parameter
* Calculate number of connection per target
* Add connections based ballancing between instances

## Milestone 1
### Retry
* Add retry function
* Calculate target response time

### Persistence
* Persistent connection

### HA Pairs
* Run in pairs
* Replicate sessions between pairs

## Milestone 2
### API
* Introduce REST API framework

### Listeners manager
* API to manage listeners

### Targets
* API to manage targets

## Milestone 3

### Stats
* API to display statistics(Prometheus format?)



