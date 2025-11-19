package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/gofiber/fiber/v2"
)

// UploadHandler sederhana
type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func UploadRoutes(route fiber.Router) {
	h := NewUploadHandler()
	route.Post("/upload", UploadFileHandler(h))
}

// UploadFile
// @Summary Upload file stream
// @Description Upload audio file (mp3) via multipart/form-data
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /v1/upload [post]
func UploadFileHandler(h *UploadHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil file dari multipart form (field name: "file")
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "file is required")
		}

		// Pastikan folder upload ada
		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return errors.InternalError(fmt.Sprintf("failed to create upload dir: %v", err))
		}

		// Buat nama file unik
		dstFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		dstPath := filepath.Join(uploadDir, dstFilename)

		// Buka source file stream
		src, err := fileHeader.Open()
		if err != nil {
			return errors.InternalError(fmt.Sprintf("failed to open uploaded file: %v", err))
		}
		defer src.Close()

		// Buat file tujuan dan stream copy (tidak memuat semua ke memory)
		dst, err := os.Create(dstPath)
		if err != nil {
			return errors.InternalError(fmt.Sprintf("failed to create destination file: %v", err))
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return errors.InternalError(fmt.Sprintf("failed to save uploaded file: %v", err))
		}

		res := map[string]interface{}{
			"filename": fileHeader.Filename,
			"stored":   dstFilename,
			"path":     dstPath,
			"size":     fileHeader.Size,
		}

		return pkg.ResponseApiOK(c, "File uploaded successfully", res)
	}
}
