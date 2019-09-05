package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vearne/p2p-sharer/models"
	"net/http"
)

func ConcurrentLimit(n int) gin.HandlerFunc {
	limter := models.NewConcurentLimiter(n)
	return func(c *gin.Context) {
		if limter.TryEnter() {
			defer limter.Exit()
			c.Next()
		} else {
			c.Abort()
			c.JSON(http.StatusTooManyRequests,
				models.ErrResponse{"E002", "too many requests"})
		}
	}
}
