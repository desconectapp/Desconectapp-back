package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            status := http.StatusInternalServerError 

            if meta, ok := c.Errors.Last().Meta.(map[string]any); ok {
                if s, ok := meta["status"].(int); ok {
                    status = s
                }
            }

            c.JSON(status, gin.H{
                "error": c.Errors.Last().Error(),
            })
            c.Abort()
            return
        }
    }
}