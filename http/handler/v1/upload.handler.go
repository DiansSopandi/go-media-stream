package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/DiansSopandi/media_stream/dto"
	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/middlewares"
	"github.com/DiansSopandi/media_stream/models"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/DiansSopandi/media_stream/repository"
	service "github.com/DiansSopandi/media_stream/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tcolgate/mp3"
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
// Param trackDto body dto.TrackCreateRequest true "Create Track Request"
// Param trackDto formData string false "Create Track Request as JSON (e.g. {\"title\":\"...\",\"artist\":\"...\",\"duration\":240,\"is_public\":true})"
// @Param title formData string false "Track title"
// @Param artist formData string false "Artist name"
// @Param album formData string false "Album name"
// @Param description formData string false "Track description"
// @Param genre formData string false "Genre"
// Param duration formData integer false "Duration in seconds"
// Param filename formData string false "Original filename or custom filename"
// @Param is_public formData boolean false "Is public (true/false)"
// Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/upload [post]
func UploadFileHandler(h *UploadHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// var createTrackDto dto.TrackCreateRequest

		// if err := c.BodyParser(&createTrackDto); err != nil {
		// 	return fiber.NewError(fiber.StatusBadRequest, err.Error())
		// }

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "file is required")
		}

		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return errors.InternalError(fmt.Sprintf("failed to create upload dir: %v", err))
		}

		dstFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		dstPath := filepath.Join(uploadDir, dstFilename)

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

		createTrackDto, err := parseTrackCreateRequest(c)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		metadata := map[string]interface{}{
			"filename": fileHeader.Filename,
			// "stored":   dstFilename,
			"path": dstPath,
			"size": fileHeader.Size,
		}

		jsonMetadata, _ := json.Marshal(metadata)
		if err != nil {
			return err
		}

		// duration, err := GetAudioDuration(dstPath)
		duration, err := GetMP3Duration(dstPath)
		if err != nil {
			return errors.InternalError(fmt.Sprintf("failed to get duration: %v", err))
		}

		createTrackDto.Metadata = jsonMetadata
		createTrackDto.Filename = fileHeader.Filename
		createTrackDto.Duration = int(duration)

		fmt.Println("Parsed Track Create DTO:", createTrackDto)
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
	// fmt.Println("Extracted User ID:", userIDInt64)
	if err != nil {
		return models.Track{}, err
	}
	// get user id from jwt claims in context
	// claims := c.Locals("user").(jwt.MapClaims)
	// userID := fmt.Sprintf("%v", claims["sub"])
	// fmt.Println("User ID from claims:", userID)

	// metadata := map[string]interface{
	// 	"filename": fileHeader.Filename,
	// 	"stored":   dstFilename,
	// 	"path":     dstPath,
	// 	"size":     fileHeader.Size,
	// }
	// if err := json.Unmarshal([]byte(metadata), &map[string]interface{}{}); err != nil {
	// 	metadata = map[string]interface{}{}
	// }
	// createTrackDto.Metadata = metadata

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
		Metadata:    uploadDto.Metadata,
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

func parseTrackCreateRequest(c *fiber.Ctx) (dto.TrackCreateRequest, error) {
	var d dto.TrackCreateRequest

	// 1) Coba parse JSON/body standar (akan no-op jika multipart)
	_ = c.BodyParser(&d)

	// 2) Jika multipart/form-data, override dari form fields jika ada
	ct := c.Get("Content-Type")
	if !strings.Contains(ct, "multipart/form-data") {
		return d, nil
	}

	// Jika ada single form field "trackDto" berisi JSON, unmarshal itu dulu (mendukung Swagger textarea JSON)
	if td := c.FormValue("trackDto"); td != "" {
		if err := json.Unmarshal([]byte(td), &d); err != nil {
			return d, fmt.Errorf("invalid trackDto json: %w", err)
		}
		// lanjut untuk override individual form fields jika ada
	}

	// helper untuk set string field dari form value jika ada
	setStr := func(field string, setter func(string)) {
		if v := c.FormValue(field); v != "" {
			setter(v)
		}
	}

	setStr("title", func(v string) { d.Title = v })
	setStr("description", func(v string) { d.Description = v })
	setStr("artist", func(v string) { d.Artist = v })
	setStr("album", func(v string) { d.Album = v })
	setStr("genre", func(v string) { d.Genre = v })
	setStr("filename", func(v string) { d.Filename = v })

	if v := c.FormValue("duration"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			d.Duration = i
		}
	}
	if v := c.FormValue("is_public"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			d.IsPublic = b
		}
	}

	// Jika Anda punya metadata di DTO sebagai string atau map, bisa di-handle di sini juga.

	return d, nil
}

func GetAudioDuration(path string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

// func GetMP3Duration(path string) (int, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer f.Close()

// 	metadata, err := tag.ReadFrom(f)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return metadata.Length(), nil // duration dalam detik
// }

func GetMP3Duration(path string) (float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	d := mp3.NewDecoder(f)
	var frame mp3.Frame
	var skipped int
	var duration float64

	for {
		if err := d.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		duration += frame.Duration().Seconds()
	}

	return duration, nil
}
