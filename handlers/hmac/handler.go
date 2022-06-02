package hmac

import (
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/projectkeas/sdks-service/configuration"
	"github.com/projectkeas/sdks-service/server"

	chmac "crypto/hmac"
	csha "crypto/sha256"
)

var (
	secret string
)

func NewSha256(headerName string, server *server.Server, configName string) func(c *fiber.Ctx) error {

	server.GetConfiguration().RegisterChangeNotificationHandler(func(config configuration.ConfigurationRoot) {
		secret = config.GetStringValueOrDefault(configName, "")
	})

	return func(c *fiber.Ctx) error {

		if secret == "" {
			c.SendStatus(fiber.StatusServiceUnavailable)
			return nil
		}

		sig := c.Get(headerName)

		if sig == "" {
			c.SendStatus(fiber.StatusBadRequest)
			return nil
		} else {
			sig = strings.TrimPrefix(sig, "sha256=")
		}

		payload := c.Body()
		mac := chmac.New(csha.New, []byte(secret))
		mac.Write(payload)

		left, _ := hex.DecodeString(sig)
		right := mac.Sum(nil)

		if chmac.Equal(left, right) {
			return c.Next()
		}

		c.SendStatus(fiber.StatusBadRequest)
		return nil
	}
}
