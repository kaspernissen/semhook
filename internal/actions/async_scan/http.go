package asyncscan

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandlerStart(repoRoot string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get the uploaded file
		file, err := c.FormFile("rule")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Create a temporary directory
		tempDir, err := ioutil.TempDir("", "temp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary directory"})
			return
		}
		defer os.RemoveAll(tempDir)

		// Save the uploaded file to the temporary directory
		rulePath := fmt.Sprintf("%s/%s", tempDir, file.Filename)
		if err := c.SaveUploadedFile(file, rulePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		result := ""
		// Return a success response
		c.JSON(http.StatusOK, gin.H{"scan_result": result})
	}
}

func HandlerStatus() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func HandlerGetResult() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
