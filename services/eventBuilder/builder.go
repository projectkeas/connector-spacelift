package eventBuilder

import (
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cee "github.com/cloudevents/sdk-go/v2/event"
)

type mapper func(map[string]interface{}) (cloudevents.Event, cee.ValidationError)

var (
	Version string            = "1.0"
	mappers map[string]mapper = map[string]mapper{
		"audit": mapAuditTrail,
	}
)

func NewCloudEventFromWebhook(payload map[string]interface{}, eventName string, eventTime time.Time) (cloudevents.Event, cee.ValidationError) {
	mapper, found := mappers[eventName]
	if found {
		event, err := mapper(payload)
		if err == nil {
			event.SetTime(eventTime)
		}
		return event, err
	}

	return cloudevents.Event{}, cee.ValidationError{
		"NotFound": fmt.Errorf("a mapper for the event '%s' was not found", eventName),
	}
}

// ============================
// 			Helpers
// ============================

func lookupValueFromPayloadPath(payload map[string]interface{}, path ...string) string {
	for _, p := range path {
		item, found := payload[p]

		if !found {
			break
		}

		switch v := item.(type) {
		case string:
			return v
		case map[string]interface{}:
			payload = v
		default:
			return fmt.Sprint(v)
		}
	}

	return ""
}
