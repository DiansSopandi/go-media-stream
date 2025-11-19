package handler

import (
	service "github.com/DiansSopandi/media_stream/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// func GalleryRoutes(route fiber.Router) {
// 	handler := NewGalleryHandler()

// 	route.Get("/galleries", middleware.LimitGlobalMiddleware(), getAllGalleriesHandler(handler))
// 	route.Get(galleryByIDRoute, middleware.LimitGlobalMiddleware(), getGalleryByIDHandler(handler))
// 	route.Get("/galleries/title/:title", middleware.LimitGlobalMiddleware(), getGalleriesByTitleHandler(handler))
// 	route.Get("/galleries/:slug/title", middleware.LimitGlobalMiddleware(), getGalleriesBySlugHandler(handler))
// 	route.Post("/galleries", middleware.LimitGlobalMiddleware(), createGalleryHandler(handler))
// 	route.Put(galleryByIDRoute, middleware.LimitGlobalMiddleware(), updateGalleryHandler(handler))
// 	route.Delete(galleryByIDRoute, middleware.LimitGlobalMiddleware(), deleteGalleryHandler(handler))
// }

// func getAllGalleriesHandler(handler *GalleryHandler) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		page := c.QueryInt("page", 1)
// 		limit := c.QueryInt("limit", 10)

// 		galleries, err := handler.GetAllGalleries(page, limit)

// 		if err != (pkg.ResponseApiError{}) {
// 			message, logError := pkg.ExtractErrorMessage(err)
// 			pkg.ResponseApiWrapper(c, logError, message, int(err.HttpStatusCode), nil, err)

// 			return c.Status(int(err.HttpStatusCode)).JSON(fiber.Map{
// 				"success": false,
// 				"status":  err.Status,
// 				"message": message,
// 				"data":    []models.Gallery{},
// 			})
// 		}

// 		// if len(articles) == 0 {
// 		// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 		// 		"success": false,
// 		// 		"status":  pkg.ApiErrorNotFound,
// 		// 		"message": "No articles found",
// 		// 		"data":    []models.Article{},
// 		// 	})
// 		// }

// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"success": true,
// 			"status":  pkg.ApiStatusSuccessOk,
// 			"message": "Articles fetch successfully...",
// 			"data":    galleries,
// 		})
// 	}
// }

// func (h *GalleryHandler) GetAllGalleries(page, limit int) (pkg.Paginator, pkg.ResponseApiError) {
// 	galleries, err := h.service.GetAllGalleries(page, limit)
// 	if err != (pkg.ResponseApiError{}) {
// 		return pkg.Paginator{}, err
// 	}
// 	return galleries, pkg.ResponseApiError{}
// }

func userRoutex(route fiber.Router, h *UserHandler) {
	handler := NewUserHandler(h.UserService)
	// Define your user-related routes here
	// For example:
	// app.Get("/api/users/:id", h.GetUserByID)
	route.Get("/users", GetAllUserHandler(handler))
	// route.Get("/users/:id", func(c *fiber.Ctx) error {
	// 	id, err := c.ParamsInt("id")
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "Invalid user ID",
	// 		})
	// 	}

	// 	user, err := h.GetUserByID(id)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Failed to retrieve user",
	// 		})
	// 	}

	// 	return c.JSON(user)
	// })
	// You can add more user-related routes here
	// route.Post("/users", func(c *fiber.Ctx) error {
	// 	var userRequest service.UserRequest
	// 	if err := c.BodyParser(&userRequest); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "Invalid request body",
	// 		})
	// 	}

	// 	user, err := h.UserService.CreateUser(&userRequest)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Failed to create user",
	// 		})
	// 	}

	// 	return c.Status(fiber.StatusCreated).JSON(user)
	// })
}

func GetAllUserHandler(h *UserHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := h.UserService.GetAllUsers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve users",
			})
		}

		return c.JSON(users)
	}
}

// func (h *UserHandler) GetUserByID(id int) (*service.UserResponse, error) {
// 	user, err := h.UserService.GetUserByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &service.UserResponse{
// 		ID:        user.ID,
// 		Username:  user.Username,
// 		Email:     user.Email,
// 		Role:      user.Role,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	}, nil
// }

// GalleryHandler methods to handle gallery-related requests
// @Summary GetAllGalleries
// @Description Retrieve all galleries
// @Tags Gallery
// @Accept json
// @Produce json
// @Param page query int false "Page number (optional)"
// @Param limit query int false "Limit per page (optional)"
// @Success 200 {array} dto.GalleryResponse
// @Failure 500 {object} pkg.ResponseApiError
// @Router /api/galleries [get]
// func (h *GalleryHandler) GetAllGalleries(page, limit int) (pkg.Paginator, pkg.ResponseApiError) {
// 	galleries, err := h.service.GetAllGalleries(page, limit)
// 	if err != (pkg.ResponseApiError{}) {
// 		return pkg.Paginator{}, err
// 	}
// 	return galleries, pkg.ResponseApiError{}
// }

// Add methods for handling user-related requests here
// For example:
// func (h *UserHandler) GetUser(c *fiber.Ctx) error {
//     // Implement user retrieval logic here
//     return nil
// }
// func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
//     // Implement user creation logic here
//     return nil
// }
// func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
//     // Implement user update logic here
//     return nil
// }
// func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
//     // Implement user deletion logic here
//     return nil
// }
// You can also add methods for handling user authentication, registration, etc.
