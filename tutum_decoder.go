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
  // "fmt"
)

type TutumDecoderConfig struct {
  // Map of message field names to message string values. Note that all
  // values *must* be strings. Any specified Pid and Severity field values
  // must be parseable as int32. Any specified Timestamp field value will be
  // parsed against the specified TimestampLayout. All specified user fields
  // will be created as strings.
  MessageFields MessageTemplate `toml:"message_fields"`
}

type TutumDecoder struct {
  messageFields MessageTemplate
}

func (td *TutumDecoder) ConfigStruct() interface{} {
  return new(TutumDecoderConfig)
}

func (td *TutumDecoder) Init(config interface{}) (err error) {
  conf := config.(*TutumDecoderConfig)
  td.messageFields = conf.MessageFields
  return
}

func (td *TutumDecoder) Decode(pack *PipelinePack) (packs []*PipelinePack, err error) {
  // fmt.Printf("Message: %v\n", pack.Message)

  // get ContainerName
  container_name := pack.Message.FindFirstField("ContainerName").GetValueString()[0]

  // parse names from container_name
  stack, service, container, uuid := ParseTutumNames(container_name)

  // add fields to message with tutum service data
  if stack != "" {
    // add fields to message with tutum stack data
    field := message.NewFieldInit("TutumStackName", message.Field_STRING, "")
    field.AddValue(stack)
    pack.Message.AddField(field)    
  }

  if service != "" {
    field := message.NewFieldInit("TutumServiceName", message.Field_STRING, "")
    field.AddValue(service)
    pack.Message.AddField(field)
  }

  if container != "" {
    // add fields to message with tutum stack data
    field := message.NewFieldInit("TutumContainerName", message.Field_STRING, "")
    field.AddValue(container)
    pack.Message.AddField(field)    
  }

  if uuid != "" {
    // add fields to message with tutum stack data
    field := message.NewFieldInit("TutumUuid", message.Field_STRING, "")
    field.AddValue(uuid)
    pack.Message.AddField(field)    
  }
  
  // fmt.Printf("Message: %v\n", pack.Message)
  if err = td.messageFields.PopulateMessage(pack.Message, nil); err != nil {
    return
  }
  
  // fmt.Printf("Message: %v\n", pack.Message)
  return []*PipelinePack{pack}, nil
}

func init() {
  RegisterPlugin("TutumDecoder", func() interface{} {
    return new(TutumDecoder)
  })
}
