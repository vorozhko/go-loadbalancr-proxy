# HTTP proxy todo

* Add target hosts and listen ports configuration thorugh YAML parameter [done]
* Calculate number of connection per target [done]
* Add connections based ballancing between instances [done]

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


## More features
[Load balancing](https://en.wikipedia.org/wiki/Load_balancing_(computing))
