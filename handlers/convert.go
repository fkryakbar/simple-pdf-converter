package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"simple-pdf-converter/utils"

	"github.com/gin-gonic/gin"
)

// ConvertPDFRequest handles the PDF to PNG conversion request
func ConvertPDF(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File is required. Please upload a file with key 'file'",
			"data":    nil,
		})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".pdf":  true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}

	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file type. Only PDF and image files (PNG, JPG, JPEG) are allowed",
			"data":    nil,
		})
		return
	}

	// Open the file
	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to open uploaded file",
			"data":    nil,
		})
		return
	}
	defer fileHeader.Close()

	// Read file content
	fileData, err := io.ReadAll(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read uploaded file",
			"data":    nil,
		})
		return
	}

	// Handle image files - convert directly to base64
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		base64Image := utils.ConvertImageToBase64(fileData)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    []string{base64Image},
		})
		return
	}

	// For PDF files, check magic bytes
	if len(fileData) < 4 || string(fileData[:4]) != "%PDF" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file content. The file is not a valid PDF",
			"data":    nil,
		})
		return
	}

	// Convert PDF to PNG base64
	base64Images, err := utils.ConvertPDFToBase64PNG(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to convert PDF: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    base64Images,
	})
}
