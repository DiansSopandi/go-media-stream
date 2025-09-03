package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/DiansSopandi/goride_be/errors"
	"github.com/DiansSopandi/goride_be/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthGuard(c *fiber.Ctx) error {
	var tokenString string

	// publicRoutes := []string{
	// 	"/",
	// 	"/v1",
	// 	"/v1/health",
	// 	"/swagger/*",
	// 	"/swagger/index.html",
	// 	"/swagger/doc.json",
	// 	"/v1/auth/login",
	// 	"/v1/auth/register",
	// 	"/v1/auth/google/login",
	// 	"/v1/auth/google/callback",
	// }

	// for _, r := range publicRoutes {
	// 	if c.Path() == r || strings.HasPrefix(c.Path(), "/swagger/") {
	// 		return c.Next()
	// 	}
	// }

	if isPublicRoute(c.Path()) {
		return c.Next()
	}

	tokenString, err := extractToken(c)
	if err != nil {
		return err
	}

	// Ambil token dari Authorization header
	// authHeader := c.Get("Authorization")
	// if authHeader == "" {
	// 	return fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
	// }
	// if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
	// 	tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	// }

	// if tokenString == "" {
	// 	tokenString = c.Cookies("jwt_at")
	// }

	// if tokenString == "" {
	// 	return errors.Unauthorized("Missing token")
	// }

	// tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	// if tokenString == authHeader {
	// 	return fiber.NewError(fiber.StatusUnauthorized, "Invalid token format")
	// }

	// Parse & validasi JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan algoritma sesuai
		// SigningMethodHS256
		// if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Unauthorized("Unexpected signing method")
		}
		return []byte(pkg.Cfg.Application.JwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return errors.Unauthorized(fmt.Sprintf("Invalid token: %v", err))
	}

	claims := token.Claims.(jwt.MapClaims)
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return errors.Unauthorized("token expired")
		}
	}

	// Simpan user info ke context (opsional)
	// claims := token.Claims.(jwt.MapClaims)
	c.Locals("user", claims)

	return c.Next()
}

func GetPublicRoutes() []string {
	return pkg.Cfg.Application.PublicRoutes
}

func isPublicRoute(path string) bool {
	publicRoutes := GetPublicRoutes()
	for _, r := range publicRoutes {
		// kalau pakai wildcard swagger/*
		if strings.HasSuffix(r, "/*") {
			prefix := strings.TrimSuffix(r, "/*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}

		if path == r {
			return true
		}
	}
	return false
}

func extractToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	tokenString := c.Cookies("jwt_at")
	if tokenString != "" {
		return tokenString, nil
	}

	return "", errors.Unauthorized("Missing token")
}
