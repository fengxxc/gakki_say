package bot

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Controller(r *gin.Engine) {
	arcRouter := r.Group("/api")

	{
		arcRouter.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"ok":      true,
				"message": "元气~",
			})
		})
	}
}
