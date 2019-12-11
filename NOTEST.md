# Notes
* Use single APP object to load the app and configuration
* Use single LB object per each Lister port
* Use multi Target objects per LB object

[Listening multiple ports on golang http servers](https://gist.github.com/filewalkwithme/0199060b2cb5bbc478c5)

Example of multi listen config file:
```
- listen: 8080
  targets: [http://localhost:8081, http://localhost:8082]
- listen: 8081
  targets: [http://localhost:8081, http://localhost:8082]
```