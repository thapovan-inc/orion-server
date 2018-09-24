package server

import (
	"errors"
	"fmt"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
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
	currentTimestamp := time.Now().UnixNano() / 1000
	maxDiff := util.GetConfig().SpanValidation.AllowedDrift
	if +(int64(timestamp) - currentTimestamp) > maxDiff {
		logger := util.GetLogger("server", "validateTimestamp")
		logger.Errorf("Received timestamp %v, current timestamp %v, difference %v, maxAllowedDrift %v", timestamp, currentTimestamp, +(int64(timestamp) - currentTimestamp), maxDiff)
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
