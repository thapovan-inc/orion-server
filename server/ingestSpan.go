package server

import (
	"errors"
	"github.com/thapovan-inc/orion-server/publisher"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
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
			logger.Debugln("Publishing data for key ", key)
			err = pub.PublishSpan(targetTopic, []byte(key), spanData)
			if err != nil {
				logger.WithError(err).Errorln("Error occurred when publishing to topic ", targetTopic)
			} else {
				return nil
			}
		} else {
			if err == nil {
				err = errors.New("publisher is nil")
			}
			logger.WithError(err).Errorln("Unable to grab a publisher")
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
