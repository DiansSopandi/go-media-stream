package repository

import (
	"database/sql"

	"github.com/DiansSopandi/media_stream/db"
	"github.com/DiansSopandi/media_stream/models"
)

type UserProviderRepository struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewUserProviderRepository(tx *sql.Tx) (*UserProviderRepository, error) {
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

	return &UserProviderRepository{
		DB: db.InitDatabase(),
		TX: tx,
	}, nil

}

func (r *UserProviderRepository) BeginUserProviderTransaction() (*sql.Tx, error) {
	return r.DB.Begin()
}

func (r *UserProviderRepository) GetUserByProviderID(providerID string) (*models.UserProvider, error) {
	var userProvider models.UserProvider

	query := `SELECT id,user_id, provider,provider_id 
			  FROM user_providers 
			  WHERE provider_id = $1`
	err := r.DB.QueryRow(query, providerID).Scan(&userProvider.ID, &userProvider.UserID, &userProvider.Provider, &userProvider.ProviderID)
	if err != nil {
		return nil, err
	}
	return &userProvider, nil
}

func (r *UserProviderRepository) CreateUserProvider(tx *sql.Tx, userProvider *models.UserProvider) error {
	query := "INSERT INTO user_providers (user_id, provider, provider_id, provider_email) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(query, userProvider.UserID, userProvider.Provider, userProvider.ProviderID, userProvider.ProviderEmail)
	return err
}
