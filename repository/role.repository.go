package repository

import (
	"database/sql"

	"github.com/DiansSopandi/media_stream/db"
	"github.com/DiansSopandi/media_stream/models"
)

type RoleRepository struct {
	DB *sql.DB
	TX *sql.Tx
}

// func NewRoleRepository(withTransaction bool) (*RoleRepository, error) {
func NewRoleRepository(tx *sql.Tx) (*RoleRepository, error) {
	// dbConn := db.InitDatabase()
	// var trx *sql.Tx

	// if withTransaction {
	// 	var err error
	// 	trx, err = dbConn.Begin()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	// 	}
	// }

	return &RoleRepository{
		DB: db.InitDatabase(),
		TX: tx,
	}, nil
}

func (r *RoleRepository) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role

	query := `SELECT id, name, description, created_at, updated_at 
			  FROM roles 
			  WHERE deleted_at IS NULL
			  ORDER BY id DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepository) CreateRoles(tx *sql.Tx, role *models.Role) (models.Role, error) {
	query := `INSERT INTO roles (name, description) VALUES ($1, $2) 
	RETURNING id, name, description, created_at, updated_at`

	err := tx.QueryRow(query, role.Name, role.Description).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	return *role, err
}

func (r *RoleRepository) GetRoleByUserID(userID int) ([]models.Role, error) {
	var roles []models.Role

	query := `SELECT r.id, r.name, r.description, r.created_at, r.updated_at 
			  FROM roles r 
			  JOIN user_roles ur ON r.id = ur.role_id 
			  WHERE ur.user_id = $1 AND r.deleted_at IS NULL`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
