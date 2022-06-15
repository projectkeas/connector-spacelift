package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectkeas/connector-spacelift/handlers/auditTrail"
	"github.com/projectkeas/connector-spacelift/handlers/hmac"
	"github.com/projectkeas/sdks-service/eventPublisher"
	"github.com/projectkeas/sdks-service/server"
)

func main() {
	app := server.New("connector-spacelift")

	app.WithEnvironmentVariableConfiguration("KEAS_")

	app.WithConfigMap("connector-spacelift-cm")
	app.WithRequiredSecret("connector-spacelift-secret")

	// The ingestion secret is required for auth with the ingestion API
	app.WithRequiredSecret("ingestion-secret")

	app.ConfigureHandlers(func(f *fiber.App, server *server.Server) {
		f.Route("integrations/spacelift", func(router fiber.Router) {
			router.Post("/audit", hmac.NewSha256("X-Signature-256", server, "spacelift.webhook.token"), auditTrail.New(server))
		})
	})

	server := app.Build()

	server.RegisterService(eventPublisher.SERVICE_NAME, eventPublisher.New(server.GetConfiguration()))

	server.Run()
}
