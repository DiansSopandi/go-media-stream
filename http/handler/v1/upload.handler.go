package handler

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/DiansSopandi/media_stream/dto"
	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/middlewares"
	"github.com/DiansSopandi/media_stream/models"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/DiansSopandi/media_stream/repository"
	service "github.com/DiansSopandi/media_stream/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// UploadHandler sederhana
type UploadHandler struct {
	service *service.TrackService
}

func NewUploadHandler() *UploadHandler {
	var tx *sql.Tx

	trackRepo, _ := repository.NewTrackRepository(tx)
	trackService := service.NewTrackService(trackRepo)
	return &UploadHandler{
		service: trackService,
	}
}

func UploadRoutes(route fiber.Router) {
	h := NewUploadHandler()
	route.Post("/upload", middlewares.WithTransaction(UploadFileHandler(h)))
}

// UploadFile
// @Summary Upload file stream
// @Description Upload audio file (mp3) via multipart/form-data
// @Tags Upload
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Param file formData file true "File to upload"
// @Param trackDto body dto.TrackCreateRequest true "Create Track Request"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /v1/upload [post]
func UploadFileHandler(h *UploadHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var createTrackDto dto.TrackCreateRequest

		if err := c.BodyParser(&createTrackDto); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		// Ambil file dari multipart form (field name: "file")
		// fileHeader, err := c.FormFile("file")
		// if err != nil {
		// 	return fiber.NewError(fiber.StatusBadRequest, "file is required")
		// }

		// Pastikan folder upload ada
		// uploadDir := "uploads"
		// if err := os.MkdirAll(uploadDir, 0755); err != nil {
		// 	return errors.InternalError(fmt.Sprintf("failed to create upload dir: %v", err))
		// }

		// Buat nama file unik
		// dstFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		// dstPath := filepath.Join(uploadDir, dstFilename)

		// Buka source file stream
		// src, err := fileHeader.Open()
		// if err != nil {
		// 	return errors.InternalError(fmt.Sprintf("failed to open uploaded file: %v", err))
		// }
		// defer src.Close()

		// Buat file tujuan dan stream copy (tidak memuat semua ke memory)
		// dst, err := os.Create(dstPath)
		// if err != nil {
		// 	return errors.InternalError(fmt.Sprintf("failed to create destination file: %v", err))
		// }
		// defer dst.Close()

		// if _, err := io.Copy(dst, src); err != nil {
		// 	return errors.InternalError(fmt.Sprintf("failed to save uploaded file: %v", err))
		// }

		// res := map[string]interface{}{
		// 	"filename": fileHeader.Filename,
		// 	"stored":   dstFilename,
		// 	"path":     dstPath,
		// 	"size":     fileHeader.Size,
		// }

		res, err := h.CreateUpload(c, &createTrackDto)

		if err != nil {
			// return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			return errors.InternalError(fmt.Sprintf("Failed to create track: %v", err))

		}

		return pkg.ResponseApiOK(c, "File uploaded successfully", res)
		// return pkg.ResponseApiOK(c, "File uploaded successfully", nil)
	}
}

func (h *UploadHandler) CreateUpload(c *fiber.Ctx, uploadDto *dto.TrackCreateRequest) (models.Track, error) {
	// Implementasi logika pembuatan upload di sini
	tx := c.Locals(middlewares.TxContextKey).(*sql.Tx)

	trackRepo, _ := repository.NewTrackRepository(tx)

	// get user id from jwt claims in context
	userIDInt64, err := extractUserIDFromClaims(c)
	fmt.Println("Extracted User ID:", userIDInt64)
	if err != nil {
		return models.Track{}, err
	}
	// get user id from jwt claims in context
	// claims := c.Locals("user").(jwt.MapClaims)
	// userID := fmt.Sprintf("%v", claims["sub"])
	// fmt.Println("User ID from claims:", userID)

	track := models.Track{
		Title:       uploadDto.Title,
		Description: uploadDto.Description,
		Artist:      uploadDto.Artist,
		Album:       uploadDto.Album,
		Genre:       uploadDto.Genre,
		Duration:    uploadDto.Duration,
		Filename:    uploadDto.Filename,
		IsPublic:    uploadDto.IsPublic,
		UserID:      int(userIDInt64),
	}

	return service.NewTrackService(trackRepo).CreateTrack(tx, &track)
}

func extractUserIDFromClaims(c *fiber.Ctx) (int64, error) {
	userClaimsAny := c.Locals("user")
	if userClaimsAny == nil {
		return 0, errors.Unauthorized("user claims not found")
	}

	var claimsMap map[string]interface{}
	switch v := userClaimsAny.(type) {
	case jwt.MapClaims:
		claimsMap = map[string]interface{}(v)
	case map[string]interface{}:
		claimsMap = v
	default:
		return 0, errors.Unauthorized("invalid user claims type")
	}

	// cek beberapa key umum untuk user id
	keys := []string{"user_id", "sub", "id", "uid"}
	for _, k := range keys {
		if val, ok := claimsMap[k]; ok {
			if id := toInt64(val); id > 0 {
				return id, nil
			}
		}
	}

	return 0, errors.Unauthorized("user id not found in token claims")
}

func toInt64(v interface{}) int64 {
	switch vv := v.(type) {
	case float64:
		return int64(vv)
	case float32:
		return int64(vv)
	case int:
		return int64(vv)
	case int64:
		return vv
	case int32:
		return int64(vv)
	case string:
		if parsed, err := strconv.ParseInt(vv, 10, 64); err == nil {
			return parsed
		}
	}
	return 0
}
