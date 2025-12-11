package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

var pool pdfium.Pool

// InitPDFium initializes the PDFium library with WebAssembly runtime
func InitPDFium() error {
	var err error
	pool, err = webassembly.Init(webassembly.Config{
		MinIdle:  1,
		MaxIdle:  3,
		MaxTotal: 10,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize PDFium: %w", err)
	}
	return nil
}

// ClosePDFium closes the PDFium pool
func ClosePDFium() {
	if pool != nil {
		pool.Close()
	}
}

// ConvertPDFToBase64PNG converts a PDF file to PNG images and returns base64 encoded strings
func ConvertPDFToBase64PNG(pdfData []byte) ([]string, error) {
	// Get an instance from the pool
	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		return nil, fmt.Errorf("failed to get PDFium instance: %w", err)
	}
	defer instance.Close()

	// Open the PDF document
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF document: %w", err)
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Get page count
	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	base64Images := make([]string, 0, pageCount.PageCount)

	// Render each page to PNG
	for i := 0; i < pageCount.PageCount; i++ {
		// Render page to image
		renderResult, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    i,
				},
			},
			DPI: 150, // Good quality for most use cases
		})
		if err != nil {
			return nil, fmt.Errorf("failed to render page %d: %w", i+1, err)
		}

		// Encode image to PNG
		var buf bytes.Buffer
		if err := png.Encode(&buf, renderResult.Result.Image); err != nil {
			return nil, fmt.Errorf("failed to encode page %d to PNG: %w", i+1, err)
		}

		// Convert to base64
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		base64Images = append(base64Images, base64Str)
	}

	return base64Images, nil
}
