package middleware

import "github.com/gin-gonic/gin"

func NewRecoveryMiddleware() Middleware {
	// for custom recover, change can be made here
	return func(c *gin.Context) {
		gin.Recovery()
	}
}
