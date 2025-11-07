package middleware

import (
	"net/http"

	"github.com/blog/internal/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, map[string]any{
				"code": utils.FAIL.Code,
				"msg":  err.Error(),
			})
		}
	}
}
