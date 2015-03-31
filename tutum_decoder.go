/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2014
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Ian Neubert (ian@ianneubert.com)
#
# ***** END LICENSE BLOCK *****/

package tutum

import (
  . "github.com/mozilla-services/heka/pipeline"
  "github.com/mozilla-services/heka/message"
  "net/http"
  "fmt"
  "encoding/json"
  "io/ioutil"
)

type TutumDecoderConfig struct {
  // Map of message field names to message string values. Note that all
  // values *must* be strings. Any specified Pid and Severity field values
  // must be parseable as int32. Any specified Timestamp field value will be
  // parsed against the specified TimestampLayout. All specified user fields
  // will be created as strings.
  MessageFields MessageTemplate `toml:"message_fields"`
  Auth string `toml:"auth"`
}

type TutumDecoder struct {
  messageFields MessageTemplate
  auth string
  services map[string]tutumService
}

type tutumService struct {
  Id string `json:"uuid"`
  Name string `json:"name"`
}

type tutumContainer struct {
  Id string `json:"uuid"`
  Name string `json:"name"`
  ServicePath string `json:"service"`
}

type tutumStack struct {
  Id string `json:"uuid"`
  Name string `json:"name"`
}

func (td *TutumDecoder) ConfigStruct() interface{} {
  return new(TutumDecoderConfig)
}

func (td *TutumDecoder) Init(config interface{}) (err error) {
  conf := config.(*TutumDecoderConfig)
  td.messageFields = conf.MessageFields
  td.auth = conf.Auth
  td.services = make(map[string]tutumService)
  return
}

func (td *TutumDecoder) Decode(pack *PipelinePack) (packs []*PipelinePack, err error) {
  // fmt.Printf("Message: %v\n", pack.Message)

  // get ContainerName
  container_name := pack.Message.FindFirstField("ContainerName").GetValueString()[0]

  // if Service is NOT in cache; then get it
  service, ok := td.services[container_name]
  if !ok {
    // find the value
    service, err = td.get_tutum_service(container_name)
    if err != nil {
      return
    }

    // save the value to the cache
    td.services[container_name] = service
  }

  // add fields to message with tutum data
  field := message.NewFieldInit("TutumServiceName", message.Field_STRING, "")
  field.AddValue(service.Name)
  pack.Message.AddField(field)

  field = message.NewFieldInit("TutumServiceId", message.Field_STRING, "")
  field.AddValue(service.Id)
  pack.Message.AddField(field)
  
  // fmt.Printf("Message: %v\n", pack.Message)
  if err = td.messageFields.PopulateMessage(pack.Message, nil); err != nil {
    return
  }
  
  // fmt.Printf("Message: %v\n", pack.Message)
  return []*PipelinePack{pack}, nil
}

func (td *TutumDecoder) get_tutum_service(container_name string) (service tutumService, err error) {
  fmt.Printf("Finding Tutum service from container named: %s\n", container_name)
  
  // get container information
  url := fmt.Sprintf("/api/v1/container/%s/", container_name)
  data, err := td.get_tutum_uri(url)
  if err != nil {
    return
  }
  // if we get back data, then let's parse it. otherwise we can return an empty service
  if len(data) > 0 {
    // parse response
    var resp tutumContainer
    err = json.Unmarshal(data, &resp)
    if err != nil {
      return
    }
    // fmt.Printf("json: %v\n\n", resp)

    // get service information
    data, err = td.get_tutum_uri(resp.ServicePath)
    if err != nil {
      return
    }

    // parse response
    err = json.Unmarshal(data, &service)
    if err != nil {
      return
    }
  }

  return
}

func (td *TutumDecoder) get_tutum_uri(uri string) (json []byte, err error) {
  base_url := "https://dashboard.tutum.co"
  url := fmt.Sprintf("%s%s", base_url, uri)

  client := &http.Client{}

  req, err := http.NewRequest("GET", url, nil)
  req.Header.Add("Authorization", td.auth)
  req.Header.Add("Accept", "application/json")

  resp, err := client.Do(req)
  if err != nil {
    return
  }

  // read body
  defer resp.Body.Close()
  json, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    return
  }

  return
}

func init() {
  RegisterPlugin("TutumDecoder", func() interface{} {
    return new(TutumDecoder)
  })
}
