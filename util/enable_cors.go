package util

import (
	"github.com/gofiber/fiber/v2"
)

// This should only be used in develop
func EnableCors(c *fiber.Ctx) {
	c.Append("Access-Control-Allow-Origin", "*")
	c.Append("Access-Control-Allow-Methods", "*")
	c.Append("Access-Control-Allow-Headers", "*")
}
