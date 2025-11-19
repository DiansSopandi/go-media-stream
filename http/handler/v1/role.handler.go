package handler

import (
	"database/sql"
	"fmt"

	"github.com/DiansSopandi/media_stream/dto"
	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/middlewares"
	"github.com/DiansSopandi/media_stream/models"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/DiansSopandi/media_stream/repository"
	service "github.com/DiansSopandi/media_stream/services"
	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	var tx *sql.Tx
	// roleRepo, _ := repository.NewRoleRepository(false)
	roleRepo, _ := repository.NewRoleRepository(tx)
	roleService := service.NewRoleService(roleRepo)
	return &RoleHandler{
		service: roleService,
	}
}

func RolesRoutes(route fiber.Router) {
	handler := NewRoleHandler()
	route.Get("/roles", GetAllRolesHandler(handler))
	route.Post("/roles", middlewares.WithTransaction(CreateRoleHandler(handler)))
}

func GetAllRolesHandler(handler *RoleHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, err := handler.GetAllRoles()
		if err != nil {
			return errors.InternalError(fmt.Sprintf("Failed to fetch roles: %v", err))
		}

		// return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// 	"message": "Role fetch successfully...",
		// 	"status":  fiber.StatusOK,
		// 	"success": true,
		// 	"data":    roles,
		// })
		return pkg.ResponseApiOK(c, "Role fetch successfully...", roles)
	}
}

func CreateRoleHandler(handler *RoleHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var createRoledto dto.RoleCreateRequest

		if err := c.BodyParser(&createRoledto); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := handler.CreateRole(c, &createRoledto)

		if err != nil {
			// return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			return errors.InternalError(fmt.Sprintf("Failed to create role: %v", err))

		}

		return pkg.ResponseApiOK(c, "Role created successfully", res)
	}
}

// GetAllRoles
// @Summary GetAllRoles
// @Description Get all roles
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Role
// @Router /v1/roles [get]
func (h *RoleHandler) GetAllRoles() ([]models.Role, error) {
	return h.service.GetAllRoles()
}

// CreateRole
// @Summary CreateRole
// @Description Create a new role
// @Tags Role
// @Accept json
// @Produce json
// @Param roleDto body dto.RoleCreateRequest true "Create Role Request"
// @Security BearerAuth
// @Success 200 {object} models.Role
// @Router /v1/roles [post]
func (h *RoleHandler) CreateRole(c *fiber.Ctx, roleDto *dto.RoleCreateRequest) (models.Role, error) {
	// Ambil TX dari context
	tx := c.Locals(middlewares.TxContextKey).(*sql.Tx)
	// isTx := tx != nil
	// roleRepo, err := repository.NewRoleRepository(isTx)
	roleRepo, _ := repository.NewRoleRepository(tx)

	role := models.Role{
		Name:        roleDto.Name,
		Description: roleDto.Description,
	}

	roleServiceWithTx := service.NewRoleService(roleRepo)

	return roleServiceWithTx.CreateRoles(tx, &role)
}
