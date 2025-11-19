package utils

import (
	"time"

	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int, email string) (string, string, error) {
	atClaims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		// "exp":   time.Now().Add(7 * time.Hour).Unix(), // Token expires in 7 hours
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Token expires in 15 minutes
		// "exp":  time.Now().Add(60 * time.Second).Unix(), // Token expires in 60 seconds
		// "exp":  time.Now().Add(15 * 60 * time.Second).Unix(), // Token expires in 15 minutes
		"type": "access_token",
	}
	rtClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		// "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // Token expires in 7 * 24 hours jwt.TimeFunc().Add(time.Hour * 24).Unix(),
		"exp":  time.Now().Add(7 * time.Hour).Unix(), // Token expires in 7 hours
		"type": "refresh_token",
		// "exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours jwt.TimeFunc().Add(time.Hour * 24).Unix(),
	}

	jwtSecret := pkg.Cfg.Application.JwtSecretKey
	atToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	rtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	// jwtSecret := pkg.Cfg.Application.JwtSecret
	// SigningKey: jwtware.SigningKey{Key: []byte(pkg.Cfg.Application.SsoJwtSecret)},

	accessToken, err := atToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := rtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
