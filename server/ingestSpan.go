// Copyright 2018-Present Thapovan Info Systems Pvt. Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"errors"
	"github.com/thapovan-inc/orion-server/publisher"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func ingestSpan(spanData *orionproto.Span, namespace string) error {
	logger := util.GetLogger("server", "ingestSpan")
	spanValidationError := validateSpan(spanData)
	var message string
	if spanValidationError != nil {
		message = spanValidationError.Error()
	} else {
		pub, err := publisher.GetSpanPublisher()
		if err == nil && pub != nil {
			targetTopic := "incoming-spans"
			key := getKeyForSpanData(spanData, namespace)
			logger.Debug("Publishing data", zap.String("key", key))
			err = pub.PublishSpan(targetTopic, []byte(key), spanData)
			if err != nil {
				logger.Error("Error occurred when publishing to topic ", zap.String("targeTopic", targetTopic), zap.Error(err))
			} else {
				return nil
			}
		} else {
			if err == nil {
				err = errors.New("publisher is nil")
			}
			logger.Error("Unable to grab a publisher", zap.Error(err))
		}
		message = "Internal server error has occured"
	}
	return errors.New(message)
}

func getKeyForSpanData(spanData *orionproto.Span, namespace string) string {
	var key strings.Builder
	key.WriteString(namespace)
	key.WriteString("_")
	key.WriteString(strings.ToLower(spanData.TraceContext.TraceId))
	key.WriteString("_")
	key.WriteString(strings.ToLower(spanData.SpanId))
	key.WriteString("_")
	var timestamp uint64 = 0
	switch event := spanData.Event.(type) {
	case *orionproto.Span_StartEvent:
		timestamp = event.StartEvent.EventId
	case *orionproto.Span_EndEvent:
		timestamp = event.EndEvent.EventId
	case *orionproto.Span_LogEvent:
		timestamp = event.LogEvent.EventId
	}
	key.WriteString(strconv.FormatUint(timestamp, 10))
	return key.String()
}
