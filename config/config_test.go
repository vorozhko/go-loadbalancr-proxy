package config

import (
	"testing"
)

func TestYamlLoad(t *testing.T) {
	yaml := []byte(`
listeners: [80,443]
targetGroups: 
- name: web-http
  fromPort: 80	
  toPort: 8080
  instances:
  - http://localhost
  - http://localhost2
  - http://localhost3
- name: web-tls
  fromPort: 443
  toPort: 8080
  instances:
  - http://localhost
  - http://localhost2
  - http://localhost3
`)
	c, err := parseYaml(yaml)
	if err != nil {
		t.Error(err)
	}
	if c == nil {
		t.Error("Config object should not be nil")
	}
	if len(c.Targets) == 0 {
		t.Error("Got 0 Targets, but expect more than 0")
	}
	if len(c.Listeners) == 0 {
		t.Error("Got 0 Listeners, but expected more than 0")
	}
}
