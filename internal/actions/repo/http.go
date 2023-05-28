package repo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler() func(c *gin.Context) {
	return func(c *gin.Context) {
		result, err := ListRepos(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute starhook"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"list": result})
	}
}
