package pkg

import (
	"github.com/DiansSopandi/media_stream/dto"
	"github.com/DiansSopandi/media_stream/models"
)

func ToUserResponse(user models.User, roles []string) dto.UserResponse {
	return dto.UserResponse{
		ID:    uint(user.ID),
		Email: user.Email,
		Roles: roles,
	}
}
