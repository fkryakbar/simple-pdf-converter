# Simple PDF Converter

A lightweight, high-performance REST API service for converting PDF documents to PNG images. Built with Go and PDFium WebAssembly runtime.

## ✨ Features

- 📄 Convert PDF pages to high-quality PNG images
- 🔐 API Key authentication
- 🚀 Fast processing with PDFium WebAssembly
- 🐳 Docker ready with minimal image size (~15-25MB)
- 📦 Base64 encoded output for easy integration

## 🛠️ Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - Fast HTTP web framework
- **PDF Engine**: [go-pdfium](https://github.com/klippa-app/go-pdfium) with WebAssembly runtime
- **Language**: Go 1.25+

## 📋 Prerequisites

- Go 1.25 or higher
- Docker (optional, for containerized deployment)

## 🚀 Quick Start

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/simple-pdf-converter.git
   cd simple-pdf-converter
   ```

2. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env and set your API_KEY
   ```

3. **Run the application**
   ```bash
   go run .
   ```

### Docker Deployment

```bash
# Build the image
docker build -t simple-pdf-converter .

# Run the container
docker run -d -p 8080:8080 \
  -e PORT=8080 \
  -e API_KEY=your_secret_api_key \
  --name pdf-converter \
  simple-pdf-converter
```

## 📡 API Reference

### Convert PDF to PNG

**Endpoint**: `POST /api/convert`

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| `x-api-key` | string | Yes | Your API key |
| `Content-Type` | multipart/form-data | Yes | Form data content type |

**Body** (form-data):
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `file` | file | Yes | PDF file to convert |

**Response**:

```json
{
  "message": "success",
  "data": [
    "iVBORw0KGgoAAAANSUhEUgAA...", // Page 1 (base64 PNG)
    "iVBORw0KGgoAAAANSUhEUgAA..." // Page 2 (base64 PNG)
    // ... more pages
  ]
}
```

**Error Response**:

```json
{
  "message": "Error description",
  "data": null
}
```

### Example Usage

**cURL**:

```bash
curl -X POST http://localhost:8080/api/convert \
  -H "x-api-key: your_api_key" \
  -F "file=@document.pdf"
```

**JavaScript (fetch)**:

```javascript
const formData = new FormData();
formData.append("file", pdfFile);

const response = await fetch("http://localhost:8080/api/convert", {
  method: "POST",
  headers: {
    "x-api-key": "your_api_key",
  },
  body: formData,
});

const result = await response.json();
// result.data contains array of base64 PNG strings
```

## ⚙️ Environment Variables

| Variable  | Default | Description                |
| --------- | ------- | -------------------------- |
| `PORT`    | `8080`  | Server port                |
| `API_KEY` | -       | API key for authentication |

## 📁 Project Structure

```
simple-pdf-converter/
├── main.go              # Application entry point
├── handlers/
│   └── convert.go       # PDF conversion handler
├── middleware/
│   └── auth.go          # API key authentication
├── utils/
│   └── pdf.go           # PDFium utilities
├── Dockerfile           # Multi-stage Docker build
├── .env.example         # Environment variables template
└── README.md
```

## 📄 License

MIT License
