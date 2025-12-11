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
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file type. Only PDF files are allowed",
			"data":    nil,
		})
		return
	}

	// Validate MIME type
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
	pdfData, err := io.ReadAll(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read uploaded file",
			"data":    nil,
		})
		return
	}

	// Check PDF magic bytes
	if len(pdfData) < 4 || string(pdfData[:4]) != "%PDF" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file content. The file is not a valid PDF",
			"data":    nil,
		})
		return
	}

	// Convert PDF to PNG base64
	base64Images, err := utils.ConvertPDFToBase64PNG(pdfData)
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
