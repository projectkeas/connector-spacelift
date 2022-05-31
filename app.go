package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/projectkeas/sdks-service/server"
)

func main() {
	app := server.New("connector-spacelift")

	app.WithEnvironmentVariableConfiguration("KEAS_")

	app.WithConfigMap("connector-spacelift-cm")
	app.WithSecret("connector-spacelift-secret")

	app.ConfigureHandlers(func(f *fiber.App, server *server.Server) {
		f.Get("/", func(c *fiber.Ctx) error {
			value := server.GetConfiguration().GetStringValueOrDefault("log.level", "not set")
			return c.SendString(fmt.Sprintf("Hello, World 👋! Log Level is: %s", value))
		})
	})

	app.Build().Run()
}
