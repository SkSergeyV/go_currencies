package middleware

import (
	//"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func IsDemoRequest(path string) bool {
	return strings.HasPrefix(path, "/api/v1/demo/")
}

func TypeServiceMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		isDemo := IsDemoRequest(c.Path())
		c.Locals("isDemo", isDemo)
		//log.Printf("path=%q, isDemo=%v", c.Path(), isDemo)
		return c.Next()
	}
}
