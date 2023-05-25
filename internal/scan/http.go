package scan

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Handler(repoRoot string) func(c *gin.Context) {
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

		// Execute the semgrep binary with the uploaded file
		result, err := scan(rulePath, repoRoot)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute semgrep"})
			return
		}

		// Process the output as needed
		// For example, you can return it in the response

		// Return a success response
		c.JSON(http.StatusOK, gin.H{"scan_result": result})
	}
}
