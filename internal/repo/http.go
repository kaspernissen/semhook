package sync

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func Handler() func(c *gin.Context) {
	return func(c *gin.Context) {
		cmd := exec.Command("starhook", "list")
		result, err := cmd.CombinedOutput()
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
