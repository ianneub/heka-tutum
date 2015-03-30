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
  // "time"
)

type TutumDecoderConfig struct {
  // Map of message field names to message string values. Note that all
  // values *must* be strings. Any specified Pid and Severity field values
  // must be parseable as int32. Any specified Timestamp field value will be
  // parsed against the specified TimestampLayout. All specified user fields
  // will be created as strings.
  MessageFields MessageTemplate `toml:"message_fields"`
  ApiKey []string `toml:"api_key"`
}

type TutumDecoder struct {
  messageFields MessageTemplate
  apikey []string
}

func (td *TutumDecoder) ConfigStruct() interface{} {
  return new(TutumDecoderConfig)
}

func (td *TutumDecoder) Init(config interface{}) (err error) {
  conf := config.(*TutumDecoderConfig)
  td.messageFields = conf.MessageFields
  td.apikey = conf.ApiKey
  return
}

func (td *TutumDecoder) Decode(pack *PipelinePack) (packs []*PipelinePack, err error) {
  // fmt.Printf("Message: %v\n", pack.Message)
  // fmt.Printf("Tags: %v\n", td.tags)

  field := message.NewFieldInit("Tags", message.Field_STRING, "")
  for _, value := range td.tags {
    field.AddValue(value)
  }
  pack.Message.AddField(field)
  // fmt.Printf("Message: %v\n", pack.Message)
  if err = td.messageFields.PopulateMessage(pack.Message, nil); err != nil {
    return
  }
  // time.Sleep(1 * time.Second)
  // fmt.Printf("Message: %v\n", pack.Message)
  return []*PipelinePack{pack}, nil
}

func init() {
  RegisterPlugin("TutumDecoder", func() interface{} {
    return new(TutumDecoder)
  })
}
