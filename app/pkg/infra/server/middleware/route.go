package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// settle runs the app's error handler for an error a handler returned, so the response, and with it the status code, is
// final when the caller reads it. It reports the error as handled. Fiber's own logger middleware does the same; without
// it a middleware that wraps the handlers would see the status 200 of a request that is about to become a 500.
func settle(c *fiber.Ctx, err error) {
	if err == nil {
		return
	}
	if handlerErr := c.App().ErrorHandler(c, err); handlerErr != nil {
		_ = c.SendStatus(fiber.StatusInternalServerError)
	}
}

// routeTemplate returns the route pattern that served the request (/users/:id), which is the only form of the path that
// may become a metric label or part of a span name. The second result is false when no route matched. Fiber then still
// reports the catch-all "/" of the first middleware, so a "/" route for any other request path means "not found" or
// "method not allowed", and the raw path must not be used. requestPath is the path as it was before the handlers ran.
func routeTemplate(c *fiber.Ctx, requestPath string) (string, bool) {
	route := c.Route()
	if route == nil || (route.Path == "/" && requestPath != "/") {
		return "", false
	}
	return route.Path, true
}

// pathSet builds a lookup of paths for inSet.
func pathSet(paths []string) map[string]struct{} {
	set := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		set[strings.ToLower(path)] = struct{}{}
	}
	return set
}

// inSet reports whether path is in set the way Fiber's router would match it: ignoring case and a trailing slash.
func inSet(set map[string]struct{}, path string) bool {
	path = strings.ToLower(path)
	if _, ok := set[path]; ok {
		return true
	}
	if len(path) > 1 {
		_, ok := set[strings.TrimSuffix(path, "/")]
		return ok
	}
	return false
}
