package service

import (
	"database/sql"
	"fmt"

	"github.com/DiansSopandi/media_stream/dto"
	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/models"
	"github.com/DiansSopandi/media_stream/pkg/utils"
	"github.com/DiansSopandi/media_stream/repository"
)

type UserService struct {
	UserRepo         *repository.UserRepository
	RoleRepo         *repository.RoleRepository
	UserProviderRepo *repository.UserProviderRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository, userProviderRepo *repository.UserProviderRepository) *UserService {
	return &UserService{
		UserRepo:         userRepo,         // repository.NewUserRepository(),
		RoleRepo:         roleRepo,         // repository.NewRoleRepository(),
		UserProviderRepo: userProviderRepo, // repository.NewUserProviderRepository(),
	}
}

func (s *UserService) BeginTransaction() (*sql.Tx, error) {
	return s.UserRepo.BeginTransaction()
}

func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// func (s *UserService) CreateUser(tx *sql.Tx, user *models.User) (models.User, error) {
func (s *UserService) CreateUser(tx *sql.Tx, createUserDto *dto.UserCreateRequest) (models.User, error) {
	password, _ := utils.HashPassword(createUserDto.Password)
	user := models.User{
		Username:   createUserDto.Username,
		Email:      createUserDto.Email,
		Password:   password,
		Provider:   "local",
		ProviderID: nil,
	}

	exists, err := s.UserRepo.CheckUsernameExistsWithTx(tx, user.Username)
	if err != nil {
		return models.User{}, errors.InternalError(fmt.Sprintf("failed to check username: %v", err))
	}

	if exists {
		return models.User{}, errors.UsernameAlreadyExists("username already exists")
	}

	exists, err = s.UserRepo.CheckEmailExistsWithTx(tx, user.Email)
	if err != nil {
		return models.User{}, errors.InternalError(fmt.Sprintf("failed to check email: %v", err))
	}

	if exists {
		return models.User{}, errors.EmailAlreadyExists("username already exists")
	}

	res, err := s.UserRepo.CreateUser(tx, &user)
	if err != nil {
		return models.User{}, errors.InternalError(fmt.Sprintf("failed to create user: %v", err))
	}

	userProvider := &models.UserProvider{
		UserID:        uint(res.ID),
		Provider:      "local",
		ProviderID:    fmt.Sprintf("%d", res.ID),
		ProviderEmail: res.Email,
	}

	s.UserProviderRepo.CreateUserProvider(tx, userProvider)
	return res, nil
}

func (s *UserService) LoginUser(tx *sql.Tx, loginDto dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	var roleNames []string

	user, errUser := s.UserRepo.GetUserByEmail(loginDto.Email)
	if errUser != nil {
		return dto.UserLoginResponse{}, errors.InternalError(fmt.Sprintf("failed to get user by email: %v", errUser))
	}

	if user == nil {
		return dto.UserLoginResponse{}, errors.UserNotFound("user not found")
	}

	if !utils.CheckPasswordHash(loginDto.Password, user.Password) {
		return dto.UserLoginResponse{}, errors.InvalidCredential("invalid email or password")
	}

	role, errRole := s.RoleRepo.GetRoleByUserID(int(user.ID))
	if errRole != nil {
		return dto.UserLoginResponse{}, errors.InternalError(fmt.Sprintf("failed to get role by user id: %v", errRole))
	}

	if role == nil {
		return dto.UserLoginResponse{}, errors.RoleNotFound("role not found")
	}

	for _, r := range role {
		roleNames = append(roleNames, r.Name)
	}

	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return dto.UserLoginResponse{}, errors.InternalError(fmt.Sprintf("failed to generate token: %v", err))
	}

	// userMap := map[string]interface{}{}
	// userBytes, _ := json.Marshal(user)
	// json.Unmarshal(userBytes, &userMap)
	// userMap["token"] = token

	userResponse := dto.UserResponse{
		ID:       uint(user.ID),
		Username: &user.Username,
		Email:    user.Email,
		Roles:    roleNames,
	}

	return dto.UserLoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) UpdateUser(tx *sql.Tx, user *models.User) error {
	err := s.UserRepo.UpdateUser(tx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ValidateRolesExist(tx *sql.Tx, roleNames []string) ([]int64, error) {
	return s.UserRepo.ValidateRolesExist(tx, roleNames)
}

func (s *UserService) AssignRolesToUserWithTx(tx *sql.Tx, userID uint, roleIDs []int64) error {
	return s.UserRepo.AssignRolesToUserWithTx(tx, userID, roleIDs)
}

func (s *UserService) UpsertGoogleUser(tx *sql.Tx, googleID, email, name, picture string) (*models.User, error) {
	user, errUser := s.UserRepo.GetUserByEmail(email)
	if errUser != nil && errUser != sql.ErrNoRows {
		return nil, errUser
	}

	userProvider, errUserProvider := s.UserProviderRepo.GetUserByProviderID(googleID)
	if errUserProvider != nil && errUserProvider != sql.ErrNoRows {
		return nil, errUserProvider
	}

	if user != nil {
		if userProvider == nil {
			userProvider = &models.UserProvider{
				UserID:        uint(user.ID),
				Provider:      "google",
				ProviderID:    googleID,
				ProviderEmail: email,
			}

			if errCreateUserProvider := s.UserProviderRepo.CreateUserProvider(tx, userProvider); errCreateUserProvider != nil {
				return nil, errCreateUserProvider
			}
		}
	}

	if user == nil {
		user = &models.User{
			Email:    email,
			Username: name,
			Picture:  picture,
		}

		user, err := s.UserRepo.CreateUser(tx, user)
		if err != nil {
			return nil, err
		}

		userProvider = &models.UserProvider{
			UserID:        uint(user.ID),
			Provider:      "google",
			ProviderID:    googleID,
			ProviderEmail: email,
		}

		if errCreateUserProvider := s.UserProviderRepo.CreateUserProvider(tx, userProvider); errCreateUserProvider != nil {
			return nil, errCreateUserProvider
		}
	}

	return user, nil
}
