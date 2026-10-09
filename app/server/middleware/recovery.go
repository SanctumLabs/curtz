//go:build legacy

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sanctumlabs/fupi/app/tools/monitoring"
)

func NewRecoveryMiddleware() Middleware {
	return func(c *gin.Context) {
		defer monitoring.RecoverWithContext(c)
		gin.Recovery()
	}
}
