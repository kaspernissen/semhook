package asyncscan

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ScanHandler struct {
	service Service
}

func NewScanHandler() ScanHandler {
	return ScanHandler{
		service: NewService(),
	}
}

func (s *ScanHandler) HandlerStart(repoRoot string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get the uploaded file
		file, err := c.FormFile("rule")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "temp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary directory"})
			return
		}
		//defer os.RemoveAll(tempDir) move to Scan struct

		// Save the uploaded file to the temporary directory
		rulePath := fmt.Sprintf("%s/%s", tempDir, file.Filename)
		if err := c.SaveUploadedFile(file, rulePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		scanID, err := s.service.scan(c.Request.Context(), rulePath, repoRoot)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to start scan %v", err)})
			return
		}
		// Return a success response
		c.JSON(http.StatusOK, gin.H{"scanid": scanID})
	}
}

func (s *ScanHandler) HandlerStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		scanID := c.Query("scanid")
		if scanID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing scanid parameter"})
			return
		}
		status, err := s.service.queryProgress(scanID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get progress of scan %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"scanid": scanID, "progress": status})
	}
}

func (s *ScanHandler) HandlerGetResult() func(c *gin.Context) {
	return func(c *gin.Context) {
		scanID := c.Query("scanid")
		if scanID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing scanid parameter"})
			return
		}
		result, err := s.service.getResult(scanID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute semgrep"})
			return
		}

		// Return a success response
		c.JSON(http.StatusOK, gin.H{"scan_result": result})
	}
}
