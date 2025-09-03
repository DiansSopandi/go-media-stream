package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/DiansSopandi/goride_be/dto"
	"github.com/DiansSopandi/goride_be/errors"
	"github.com/DiansSopandi/goride_be/middlewares"
	"github.com/DiansSopandi/goride_be/pkg"
	helper "github.com/DiansSopandi/goride_be/pkg/helper"
	"github.com/DiansSopandi/goride_be/pkg/utils"
	"github.com/DiansSopandi/goride_be/repository"
	service "github.com/DiansSopandi/goride_be/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type AuthHandler struct {
	AuthService *service.UserService
}

var (
	oauthConfig     *oauth2.Config
	appJWTSecret    []byte
	frontendURL     string
	appJWTExpiresIn time.Duration
)

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func randState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func OAuthConfig() *oauth2.Config {
	clientID := pkg.Cfg.Application.GoogleClientID
	clientSecret := pkg.Cfg.Application.GoogleClientSecret
	redirectURI := pkg.Cfg.Application.GoogleRedirectURI

	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
	return oauthConfig
}

func AuthRoutes(route fiber.Router) {
	handler := NewAuthHandler()
	route.Get("/google/login", GetGoogleAuth)
	route.Get("/google/callback", middlewares.WithTransaction(GetGoogleCallback))
	route.Post("/register", middlewares.WithTransaction(RegisterUserHandler(handler)))
	route.Post("/login", middlewares.WithTransaction(LoginUserHandler(handler)))
	route.Post("/logout", middlewares.WithTransaction(LogoutUserHandler(handler)))
}

// GetGoogleAuth godoc
// @Summary Google Auth
// @Description Initiates Google OAuth2 login
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/auth/google/login [get]
func GetGoogleAuth(c *fiber.Ctx) error {
	OAuthConfig()
	state, err := randState()
	if err != nil {
		return fiber.NewError(500, "failed to create state")
	}

	// Optionally allow returnTo param to persist across redirects
	returnTo := c.Query("returnTo", "/")
	// Store state in cookie for CSRF protection
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state + "|" + url.QueryEscape(returnTo),
		HTTPOnly: true,
		Secure:   false, // set true in production with HTTPS
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
	})

	authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(authURL)
}

// GetGoogleCallback godoc
// @Summary Google Callback
// @Description Handles Google OAuth2 callback
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/auth/google/callback [get]
func GetGoogleCallback(c *fiber.Ctx) error {
	frontendURL := pkg.Cfg.Application.FrontendURL
	// appJWTSecret := pkg.Cfg.Application.JwtSecretKey

	OAuthConfig()
	ctx := context.Background()
	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		return fiber.NewError(400, "missing state or code")
	}

	// Validate state from cookie
	st := c.Cookies("oauth_state")
	parts := strings.SplitN(st, "|", 2)
	if st == "" || parts[0] != state {
		return fiber.NewError(400, "invalid state")
	}
	returnTo, _ := url.QueryUnescape(parts[1])
	if returnTo == "" {
		returnTo = "/"
	}

	tok, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Println("token exchange error:", err)
		return fiber.NewError(401, "failed to exchange code")
	}

	rawIDToken, _ := tok.Extra("id_token").(string)
	if rawIDToken == "" {
		return fiber.NewError(401, "id_token not found in token response")
	}

	// Verify ID Token signature & audience
	payload, err := idtoken.Validate(ctx, rawIDToken, oauthConfig.ClientID)
	if err != nil {
		log.Println("id token validate error:", err)
		return fiber.NewError(401, "invalid id_token")
	}

	// Extract user info from claims
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	sub, _ := payload.Claims["sub"].(string) // Google user id
	if email == "" || sub == "" {
		return fiber.NewError(400, "missing email/sub in id_token")
	}

	// === Upsert user in DB (pseudo) ===
	// In production, check if user exists by provider=google & provider_id=sub
	// If not exists, create user record linked to provider.

	tx := c.Locals(middlewares.TxContextKey).(*sql.Tx)
	userRepo, _ := repository.NewUserRepository(tx)
	roleRepo, _ := repository.NewRoleRepository(tx)
	userProviderRepo, _ := repository.NewUserProviderRepository(tx)

	userServiceWithTx := service.NewUserService(userRepo, roleRepo, userProviderRepo)

	fmt.Println("upserting google id:", sub)
	user, err := userServiceWithTx.UpsertGoogleUser(tx, sub, email, name, picture)
	if err != nil {
		return errors.InternalError(fmt.Sprintf("failed to upsert user: %v", err))
	}
	// userID := fmt.Sprintf("google:%s", sub)

	// Issue your own app session JWT
	// appToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub":     userID,
	// 	"email":   email,
	// 	"name":    name,
	// 	"picture": picture,
	// 	"iat":     time.Now().Unix(),
	// 	"exp":     time.Now().Add(appJWTExpiresIn).Unix(),
	// })

	// appJWT, err := appToken.SignedString(appJWTSecret)
	// if err != nil {
	// 	return fiber.NewError(500, "failed to sign app token")
	// }
	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return fiber.NewError(500, "failed to generate app token")
	}
	// Redirect back to frontend callback with token in query
	// r := fmt.Sprintf("%s/auth/callback?token=%s&returnTo=%s", strings.TrimRight(frontendURL, "/"), url.QueryEscape(accessToken), url.QueryEscape(returnTo))

	c.Cookie(&fiber.Cookie{
		Name:     "jwt_at",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "jwt_rt",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	// r := fmt.Sprintf(
	// 	"%s/auth/callback?accessToken=%s&refreshToken=%s&returnTo=%s",
	// 	strings.TrimRight(frontendURL, "/"),
	// 	url.QueryEscape(accessToken),
	// 	url.QueryEscape(refreshToken),
	// 	url.QueryEscape(returnTo),
	// )
	// return c.Redirect(r)
	return c.Redirect(frontendURL)
}

func RegisterUserHandler(handler *AuthHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var registerDto dto.UserRegisterRequest
		if errDto := c.BodyParser(&registerDto); errDto != nil {
			return pkg.ResponseApiErrorBadRequest(c, fmt.Sprintf("Failed to parse request body: %v", errDto))
		}

		if err := helper.ValidateRegisterUserRequest(&registerDto); err != nil {
			return pkg.ResponseApiErrorBadRequest(c, fmt.Sprintf("Validation failed: %v", err))
		}

		res, err := handler.RegisterUser(c, registerDto)
		if err != nil {
			// return pkg.ResponseApiError(c, fiber.StatusInternalServerError, "Failed to register user", err)
			// return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			return err
		}

		return pkg.ResponseApiOK(c, "User registered successfully...", res)
	}
}

func LoginUserHandler(handler *AuthHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginDto dto.UserLoginRequest
		if err := c.BodyParser(&loginDto); err != nil {
			return pkg.ResponseApiErrorBadRequest(c, fmt.Sprintf("Failed to parse request body: %v", err))
		}

		res, err := handler.LoginUser(c, loginDto)
		if err != nil {
			// return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			return err
		}

		return pkg.ResponseApiOK(c, "User logged in successfully...", res)
	}
}

// LogoutUserHandler handles user logout by clearing cookies.
// @summary Logout a user
// @tags Auth
// @produce json
// @security BearerAuth
// @router /v1/auth/logout [post]
func LogoutUserHandler(handler *AuthHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Clear cookies
		// c.ClearCookie("jwt_at", "/", "")
		// c.ClearCookie("jwt_rt", "/", "")

		expired := time.Now().Add(-time.Hour)

		c.Cookie(&fiber.Cookie{
			Name:     "jwt_at",
			Value:    "",
			Expires:  expired,
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})

		c.Cookie(&fiber.Cookie{
			Name:     "jwt_rt",
			Value:    "",
			Expires:  expired,
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})
		return pkg.ResponseApiOK(c, "User logged out successfully", nil)
	}
}

// RegisterUser handles user registration, including role assignment.
// @summary Register a new user with roles
// @description Register a new user and assign roles if provided.
// @tags Auth
// @accept json
// @produce json
// @param registerDto body dto.UserRegisterRequest true "User registration data"
// @success 200 {object} models.User "User registration successful"
// @failure 400 {object} map[string]interface{} "Bad request, validation errors"
// @failure 500 {object} map[string]interface{} "Internal server error, database or service errors"
// @router /v1/auth/register [post]
func (h *AuthHandler) RegisterUser(c *fiber.Ctx, regDto dto.UserRegisterRequest) (dto.UserResponse, error) {
	tx := c.Locals(middlewares.TxContextKey).(*sql.Tx)
	userRepo, _ := repository.NewUserRepository(tx)
	roleRepo, _ := repository.NewRoleRepository(tx)
	userProviderRepo, _ := repository.NewUserProviderRepository(tx)

	userServiceWithTx := service.NewUserService(userRepo, roleRepo, userProviderRepo)

	var (
		roleIDs []int64
		err     error
	)
	if len(regDto.Roles) > 0 {
		roleIDs, err = userServiceWithTx.ValidateRolesExist(tx, regDto.Roles)

		if err != nil {
			return dto.UserResponse{}, errors.RoleValidationFailed(fmt.Sprintf("role validation failed: %v", err))
		}
	}

	registerDto := dto.UserCreateRequest{
		Username: regDto.Username,
		Email:    regDto.Email,
		Password: regDto.Password,
		Roles:    []string{"user"},
	}

	res, err := userServiceWithTx.CreateUser(tx, &registerDto)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if len(roleIDs) > 0 {
		err = userServiceWithTx.AssignRolesToUserWithTx(tx, uint(res.ID), roleIDs)
		if err != nil {
			return dto.UserResponse{}, errors.InternalError(fmt.Sprintf("failed to assign roles to user: %v", err))
		}
	}

	return dto.UserResponse{
		ID:       uint(res.ID),
		Username: &registerDto.Username,
		Email:    res.Email,
		Roles:    registerDto.Roles,
	}, nil
}

// LoginUser handles user login and returns user details and token.
// @summary Login a user
// @description Login a user and return user details and token.
// @tags Auth
// @accept json
// @produce json
// @param loginDto body dto.UserLoginRequest true "User login data"
// @success 200 {object} dto.UserLoginResponse "User login successful"
// @failure 400 {object} map[string]interface{} "Bad request, validation errors"
// @failure 500 {object} map[string]interface{} "Internal server error, database or service errors"
// @router /v1/auth/login [post]
func (h *AuthHandler) LoginUser(c *fiber.Ctx, loginDto dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	tx := c.Locals(middlewares.TxContextKey).(*sql.Tx)
	userRepo, _ := repository.NewUserRepository(tx)
	roleRepo, _ := repository.NewRoleRepository(tx)
	userProviderRepo, _ := repository.NewUserProviderRepository(tx)

	userServiceWithTx := service.NewUserService(userRepo, roleRepo, userProviderRepo)

	res, err := userServiceWithTx.LoginUser(tx, loginDto)
	if err != nil {
		return dto.UserLoginResponse{}, err
	}

	isProd := pkg.Cfg.Application.Env == "production"
	// Set cookie Access Token
	c.Cookie(&fiber.Cookie{
		Name:  "jwt_at",
		Value: res.AccessToken,
		// Expires:  time.Now().Add(15 * time.Minute),
		Expires:     time.Now().Add(24 * time.Hour),
		HTTPOnly:    true,
		Secure:      isProd,
		SameSite:    "None",
		Path:        "/",
		SessionOnly: true,
	})

	// Set cookie Refresh Token
	c.Cookie(&fiber.Cookie{
		Name:        "jwt_rt",
		Value:       res.RefreshToken,
		Expires:     time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly:    true,
		Secure:      isProd,
		SameSite:    "None",
		Path:        "/",
		SessionOnly: true,
	})

	// Set Authorization header
	c.Set("Authorization", "Bearer "+res.AccessToken)

	return res, nil
}
