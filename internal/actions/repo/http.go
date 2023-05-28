package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler() func(c *gin.Context) {
	return func(c *gin.Context) {
		result, err := ListRepos(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute starhook: %+v", err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"list": result})
	}
}
