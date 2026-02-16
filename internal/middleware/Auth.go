package middleware

import (
	"Kopsis-Spensa/internal/config"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)
	if err != nil {
		return c.Redirect("/login")
	}

	uid := sess.Get("users_uid")
	if uid == nil {
		return c.Redirect("/login")
	}

	c.Locals("users_uid", uid)
	c.Locals("Role", sess.Get("Role"))
	c.Locals("NamaLengkap", sess.Get("NamaLengkap"))
	c.Locals("Username", sess.Get("Username"))

	return c.Next()
}

func RedirectIfLoggedIn(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)

	if err != nil && sess.Get("user_uid") != nil {
		return c.Redirect("/login")
	}
	return c.Next()
}
