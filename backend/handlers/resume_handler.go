package handlers

import "github.com/gofiber/fiber/v2"

func HandleGetResume(c *fiber.Ctx) error {
	return c.Render("resume", fiber.Map{"PageTitle": "Resume"}, "layouts/base")
}
