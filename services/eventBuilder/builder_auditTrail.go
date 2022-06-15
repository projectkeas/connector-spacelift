package eventBuilder

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	cee "github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

func mapAuditTrail(payload map[string]interface{}) (cloudevents.Event, cee.ValidationError) {
	event := cloudevents.NewEvent(Version)

	event.SetID(uuid.New().String()) // TODO : Raise with Spacelift
	event.SetSource(lookupValueFromPayloadPath(payload, "account"))
	event.SetType("io.spacelift." + lookupValueFromPayloadPath(payload, "action"))
	event.SetData(*cloudevents.StringOfApplicationJSON(), payload)
	event.SetDataSchema("https://schemas.keas.io/spacelift/audit/0.1.0")

	err := event.Validate()
	if err != nil {
		return event, err.(cee.ValidationError)
	}

	return event, nil
}
