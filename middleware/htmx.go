package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func HtmxMiddleware(c *fiber.Ctx, redirect ...string) error {
	var redirectPath string
	if len(redirect) > 0 {
		redirectPath = redirect[0]
	}

	if c.Get("HX-Request") == "" {
		if redirectPath != "" {
			return c.Redirect(redirectPath)
		}
		return c.Render("dashboard", fiber.Map{"title": "Dashboard", "hx_link": c.Path(), "IsAdmin": c.Locals("is_admin")})
	}
	c.Locals("hx_link", c.Path())
	return c.Next()
}
