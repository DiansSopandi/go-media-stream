package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/DiansSopandi/media_stream/db"
	"github.com/DiansSopandi/media_stream/dto"
	"github.com/DiansSopandi/media_stream/models"
	"github.com/lib/pq"
)

type UserRepository struct {
	DB *sql.DB
	TX *sql.Tx
}

// NewUserRepository creates a new UserRepository instance.
// func NewUserRepository(tx *sql.Tx, withTransaction bool) (*UserRepository, error) {
func NewUserRepository(tx *sql.Tx) (*UserRepository, error) {
	// dbConn := db.InitDatabase()

	// var trx *sql.Tx
	// if withTransaction {
	// 	var err error
	// 	trx, err = dbConn.Begin()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	// 	}
	// }

	// isTx := tx != nil
	// if isTx {
	// 	trx = tx
	// }

	return &UserRepository{
		DB: db.InitDatabase(),
		TX: tx,
	}, nil
}

func (r *UserRepository) BeginTransaction() (*sql.Tx, error) {
	return r.DB.Begin()
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, email, password,  created_at, updated_at, deleted_at 
	FROM users WHERE id = $1`

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		// &user.AvatarUrl,
		// &user.AvatarName,
		// &user.FirstName,
		// &user.LastName,
		// &user.Phone,
		// &user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password,  created_at, updated_at, deleted_at 
	FROM users WHERE email = $1`

	var user models.User
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(tx *sql.Tx, user *models.User) (models.User, error) {
	// query := `INSERT INTO users (username, email, password, avatar_url, avatar_name, first_name, last_name, phone, address, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	query := `INSERT INTO users (username, email, password, provider, provider_id, picture) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, email, created_at, updated_at`

	// err := r.DB.QueryRow(query,
	// err := tx.QueryRow(query,
	err := tx.QueryRow(query,
		user.Username,
		user.Email,
		user.Password,
		user.Provider,
		user.ProviderID,
		user.Picture,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	return *user, err
}
func (r *UserRepository) UpdateUser(tx *sql.Tx, user *models.User) error {
	// query := `UPDATE users SET username = $1, email = $2, password = $3, avatar_url = $4, avatar_name = $5, first_name = $6, last_name = $7, phone = $8, address = $9, role = $10, updated_at = NOW() WHERE id = $11`
	query := `UPDATE users SET username = $1, email = $2, password = $3, picture = $4, updated_at = NOW() WHERE id = $5`

	_, err := tx.Exec(query,
		user.Username,
		user.Email,
		user.Password,
		// user.AvatarUrl,
		// user.AvatarName,
		// user.FirstName,
		// user.LastName,
		// user.Phone,
		// user.Address,
		user.Picture,
		user.ID,
	)
	return err
}
func (r *UserRepository) DeleteUser(tx *sql.Tx, id int) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	_, err := tx.Exec(query, id)
	return err
}
func (r *UserRepository) GetAllUsers() ([]dto.UserResponse, error) {
	// query := `SELECT id, username, email, password, avatar_url, avatar_name, first_name, last_name, phone, address, role, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL`
	query := `SELECT u.id,  u.email, COALESCE(ARRAY_AGG(r.name), '{}') as roles 
	FROM users u
	INNER JOIN user_roles usr ON u.id = usr.user_id
	INNER JOIN roles r ON usr.role_id = r.id 
	WHERE u.deleted_at IS NULL
	GROUP BY u.id, u.username, u.email, u.created_at, u.updated_at `

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []dto.UserResponse
	var roles pq.StringArray
	for rows.Next() {
		var user dto.UserResponse //models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&roles,
			// &user.AvatarUrl,
			// &user.AvatarName,
			// &user.FirstName,
			// &user.LastName,
			// &user.Phone,
			// &user.Address,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		user.Roles = roles
		users = append(users, user)
	}
	return users, nil
}
func (r *UserRepository) CountUsers() (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	var count int
	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) ValidateRolesExist(tx *sql.Tx, roleNames []string) ([]int64, error) {
	if len(roleNames) == 0 {
		return nil, fmt.Errorf("at least one role is required")
	}

	// Remove duplicates
	uniqueRoles := make(map[string]bool)
	var cleanRoles []string

	for _, role := range roleNames {
		if !uniqueRoles[role] {
			uniqueRoles[role] = true
			cleanRoles = append(cleanRoles, role)
		}
	}

	// Build query dengan placeholders
	placeholders := make([]string, len(cleanRoles))
	args := make([]interface{}, len(cleanRoles))
	for i, name := range cleanRoles {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = name
	}

	query := fmt.Sprintf(`
        SELECT id, name 
        FROM roles 
        WHERE name IN (%s) 
        AND deleted_at IS NULL`,
		// AND is_active = true`,
		strings.Join(placeholders, ","))

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	foundRoles := make(map[string]int64)
	for rows.Next() {
		var roleID int64
		var roleName string
		if err := rows.Scan(&roleID, &roleName); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		foundRoles[roleName] = roleID
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading role rows: %w", err)
	}

	// Check missing roles
	if len(foundRoles) != len(cleanRoles) {
		var missingRoles []string
		for _, roleName := range cleanRoles {
			if _, exists := foundRoles[roleName]; !exists {
				missingRoles = append(missingRoles, roleName)
			}
		}
		return nil, fmt.Errorf("roles not found: %v", missingRoles)
	}

	// Return role IDs
	roleIDs := make([]int64, len(cleanRoles))
	for i, roleName := range cleanRoles {
		roleIDs[i] = foundRoles[roleName]
	}

	return roleIDs, nil
}

func (r *UserRepository) CheckUsernameExistsWithTx(tx *sql.Tx, username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1 AND deleted_at IS NULL`

	// err := tx.QueryRow(query, username).Scan(&count)
	err := tx.QueryRow(query, username).Scan(&count)
	return count > 0, err
}

func (r *UserRepository) CheckEmailExistsWithTx(tx *sql.Tx, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND deleted_at IS NULL`

	// err := tx.QueryRow(query, email).Scan(&count)
	err := tx.QueryRow(query, email).Scan(&count)
	return count > 0, err
}

func (r *UserRepository) AssignRolesToUserWithTx(tx *sql.Tx, userID uint, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	// Batch insert untuk user_roles
	valueStrings := make([]string, len(roleIDs))
	args := make([]interface{}, len(roleIDs)*2)

	for i, roleID := range roleIDs {
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		args[i*2] = userID
		args[i*2+1] = roleID
	}

	query := fmt.Sprintf(`
        INSERT INTO user_roles (user_id, role_id) 
        VALUES %s 
        ON CONFLICT (user_id, role_id) DO NOTHING`,
		strings.Join(valueStrings, ","))

	_, err := tx.Exec(query, args...)
	if err != nil {
		fmt.Println("error assign to roles user", err)
		return fmt.Errorf("failed to assign roles: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, provider, provider_id, picture
			  FROM users 
			  WHERE provider = 'google' AND provider_id = $1 AND deleted_at IS NULL`
	err := r.DB.QueryRow(query, googleID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Provider,
		&user.ProviderID,
		&user.Picture,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
