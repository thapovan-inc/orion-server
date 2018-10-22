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
	"fmt"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
	"math"
	"regexp"
	"time"
)

var uuidRegex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func validateSpan(span *orionproto.Span) error {
	uuidValidationError := validateUUIDs(span)
	if uuidValidationError != nil {
		return uuidValidationError
	}
	timestampValidationError := validateTimestamp(span.Timestamp)
	if timestampValidationError != nil {
		return timestampValidationError
	}
	eventIDError := validateEventID(span)
	if eventIDError != nil {
		return eventIDError
	}
	return nil
}

func validateUUIDs(spanData *orionproto.Span) error {
	isValid := isValidUUID(spanData.TraceContext.TraceId)
	var err error
	if !isValid {
		err = errors.New("invalid TraceID")
	} else {
		isValid = isValidUUID(spanData.SpanId)
		if !isValid {
			err = errors.New("invalid SpanID")
		}
	}
	return err

}

func isValidUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

func validateTimestamp(timestamp uint64) error {
	logger := util.GetLogger("server", "validateTimestamp")
	currentTimestamp := time.Now().UnixNano() / 1000
	maxDiff := util.GetConfig().SpanValidation.AllowedDrift
	if int64(math.Abs(float64(timestamp)-float64(currentTimestamp))) > maxDiff {
		logger.Sugar().Errorf("Received timestamp %v, current timestamp %v, difference %v, maxAllowedDrift %v", timestamp, currentTimestamp, +(int64(timestamp) - currentTimestamp), maxDiff)
		return fmt.Errorf("timestamp exceeds maximum allowed drift")
	}
	return nil
}

func validateEventID(spanData *orionproto.Span) error {
	var eventID uint64
	switch event := spanData.Event.(type) {
	case *orionproto.Span_StartEvent:
		eventID = event.StartEvent.EventId
	case *orionproto.Span_EndEvent:
		eventID = event.EndEvent.EventId
	case *orionproto.Span_LogEvent:
		eventID = event.LogEvent.EventId
	}
	if eventID < 1 {
		return errors.New("invalid event ID")
	}
	return nil
}
