package sync

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Execute the starhook binary with the "sync" argument
		result, err := sync(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute starhook"})
			return
		}

		// Process the output as needed
		// For example, you can return it in the response

		// Return a success response
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Sync completed: %s", result)})
	}
}
