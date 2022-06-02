package auditTrail

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/projectkeas/connector-spacelift/services/httpRequest"
	"github.com/projectkeas/sdks-service/configuration"
	"github.com/projectkeas/sdks-service/server"
)

var (
	apiKey    string
	targetUri string
)

func New(server *server.Server) func(c *fiber.Ctx) error {
	server.GetConfiguration().RegisterChangeNotificationHandler(func(config configuration.ConfigurationRoot) {
		apiKey = config.GetStringValueOrDefault("ingestion.auth.token", "")
		targetUri = config.GetStringValueOrDefault("ingestion.uri", "http://keas-ingestion.keas.svc.cluster.local/ingest")
	})

	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		payload := map[string]interface{}{}
		c.BodyParser(&payload)

		unixTimestamp, found := payload["timestamp"]
		var eventTime string
		if found {
			// Format of the timestamp is in milliseconds
			// The json value is actually a float64 so we need to convert through that first
			ut := int64(unixTimestamp.(float64))
			utt := time.UnixMilli(ut)
			eventTime = utt.UTC().Format(time.RFC3339)
		} else {
			eventTime = time.Now().UTC().Format(time.RFC3339)
		}

		payload["timestamp"] = eventTime

		envelope := map[string]interface{}{
			"metadata": map[string]string{
				"source":    "Spacelift",
				"version":   "0.1.0",
				"type":      "AuditEntry",
				"eventTime": eventTime,
			},
			"payload": &payload,
		}

		statusCode, responseBody, err := httpRequest.PostJson(targetUri, envelope, map[string]string{
			"Authorization": "ApiKey " + apiKey,
		})

		c.Response().Header.Set("Content-Type", "application/json")
		c.Send([]byte(responseBody))
		c.SendStatus(statusCode)
		return err
	}
}
