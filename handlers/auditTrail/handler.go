package auditTrail

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/projectkeas/connector-spacelift/services/eventBuilder"
	"github.com/projectkeas/sdks-service/eventPublisher"
	"github.com/projectkeas/sdks-service/server"
)

func New(server *server.Server) func(c *fiber.Ctx) error {
	// publisher setup
	nc, err := server.GetService(eventPublisher.SERVICE_NAME)
	if err != nil {
		panic(err)
	}
	client := (*nc).(eventPublisher.EventPublisherService)

	return func(context *fiber.Ctx) error {
		context.Accepts("application/json")
		errorResult := map[string]interface{}{
			"message": "An error occurred whilst processing your request",
		}

		payload := map[string]interface{}{}
		context.BodyParser(&payload)

		unixTimestamp, found := payload["timestamp"]
		var eventTime time.Time
		if found {
			// Format of the timestamp is in milliseconds
			// The json value is actually a float64 so we need to convert through that first
			ut := int64(unixTimestamp.(float64))
			utt := time.UnixMilli(ut)
			eventTime = utt.UTC()
		} else {
			eventTime = time.Now().UTC()
		}

		cloudEvent, validationErr := eventBuilder.NewCloudEventFromWebhook(payload, "audit", eventTime)
		if validationErr != nil {
			errorResult["reason"] = "Unable to validate as a cloud event"
			errors := []map[string]string{}
			for key, value := range validationErr {
				errors = append(errors, map[string]string{
					"attribute": key,
					"error":     value.Error(),
				})
			}
			errorResult["errors"] = errors
			return context.Status(fiber.StatusBadRequest).JSON(errorResult)
		}

		if !client.Publish(cloudEvent) {
			errorResult["reason"] = "publish"
			return context.Status(fiber.StatusInternalServerError).JSON(errorResult)
		}

		context.Status(fiber.StatusAccepted).JSON(cloudEvent)
		return nil
	}
}
