package tutum

import (
  "strings"
)

func ParseTutumNames(container_name string) (stack, service, container, uuid string) {
  names := strings.Split(container_name, ".")
  switch len(names) {
  case 1:
    // only has a UUID
    uuid = names[0]
  case 2:
    // only has a service
    service = getServiceNameFromContainerName(names[0])
    container = names[0]
    uuid = names[1]
  case 3:
    // contains a service and stack
    stack = names[1]
    service = getServiceNameFromContainerName(names[0])
    container = names[0]
    uuid = names[2]
  }

  return
}

func getServiceNameFromContainerName(service string) (string) {
  parts := strings.Split(service, "-")
  parts = parts[:len(parts)-1]
  return strings.Join(parts, "-")
}
