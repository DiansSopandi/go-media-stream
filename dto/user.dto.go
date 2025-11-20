package dto

type UserCreateRequest struct {
	Username string   `json:"username" example:"John Doe"`
	Email    string   `json:"email" example:"Q2Sb9@example.com"`
	Password string   `json:"password" example:"Cilok99!@"`
	Roles    []string `json:"roles" example:"admin,driver,user,superadmin"` // multiple roles
	// AvatarUrl  string `json:"avatarUrl" db:"avatar_url"`
	// AvatarName string `json:"avatarName" db:"avatar_name"`
	// FirstName  string `json:"firstName" db:"first_name"`
	// LastName   string `json:"lastName" db:"last_name"`
	// Phone      string `json:"phone" db:"phone"`
	// Address    string `json:"address" db:"address"`
}

type UserRegisterRequest struct {
	Username        string   `json:"username" validate:"required,min=3,max=50" example:"John Doe"`
	Email           string   `json:"email" validate:"required,email" example:"Q2Sb9@example.com"`
	Password        string   `json:"password" validate:"required,min=8" example:"Cilok99!@"`
	PasswordConfirm string   `json:"password_confirm" validate:"required,eqfield=Password" example:"Cilok99!@"`
	Roles           []string `json:"roles" example:"user" validate:"required"`
	// AvatarUrl  string `json:"avatarUrl" db:"avatar_url"`
	// AvatarName string `json:"avatarName" db:"avatar_name"`
	// FirstName  string `json:"firstName" db:"first_name"`
	// LastName   string `json:"lastName" db:"last_name"`
	// Phone      string `json:"phone" db:"phone"`
	// Address    string `json:"address" db:"address"`
}

type UserUpdateRequest struct {
	Username string   `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string   `json:"email,omitempty" validate:"omitempty,email"`
	Password string   `json:"password,omitempty" validate:"omitempty,min=8"`
	Roles    []string `json:"roles,omitempty"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"admin@redtech.co.id"`
	Password string `json:"password" validate:"required" example:"redtech.co.id"`
}

type UserResponse struct {
	ID       uint     `json:"id"`
	Username *string  `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

type UserLoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
