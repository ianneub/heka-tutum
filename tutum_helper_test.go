package tutum

import "testing"

type Names struct {
  Stack string
  Service string
  Container string
  Uuid string
}

type testpair struct {
  container_name string
  names Names
}

var tests = []testpair {
  {"dash-web-1.mystack.8767ab65", Names{"mystack", "dash-web", "dash-web-1", "8767ab65"}},
  {"web-1.mystack.8767ab65", Names{"mystack", "web", "web-1", "8767ab65"}},
  {"web-1.8767ab65", Names{"", "web", "web-1", "8767ab65"}},
  {"dash-web-1.8767ab65", Names{"", "dash-web", "dash-web-1", "8767ab65"}},
  {"afdcc06e-114c-4425-b2a9-455c4d32cadf", Names{"", "", "", "afdcc06e-114c-4425-b2a9-455c4d32cadf"}},
  {"weave", Names{"", "", "", "weave"}},
}

func TestWithStackName(t *testing.T) {
  for _, pair := range tests {    
    stack, service, container, uuid := ParseTutumNames(pair.container_name)

    if stack != pair.names.Stack {
      t.Errorf("Expected '%s', got: '%s'", pair.names.Stack, stack)
    }

    if service != pair.names.Service {
      t.Errorf("Expected '%s', got: '%s'", pair.names.Service, service)
    }

    if container != pair.names.Container {
      t.Errorf("Expected '%s', got: '%s'", pair.names.Container, container)
    }

    if uuid != pair.names.Uuid {
      t.Errorf("Expected '%s', got: '%s'", pair.names.Uuid, uuid)
    }
  }
}
