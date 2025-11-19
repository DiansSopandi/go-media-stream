package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/DiansSopandi/media_stream/pkg"
	_ "github.com/lib/pq"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func CreateDatabaseIfNotExists() {
	// host := pkg.GetEnv("DB_HOST")
	// port := pkg.GetEnv("DB_PORT")
	// user := pkg.GetEnv("DB_USER")
	// password := pkg.GetEnv("DB_PASSWORD")
	// dbname := pkg.GetEnv("DB_NAME")
	host := pkg.Cfg.Database.Host
	port := pkg.Cfg.Database.Port
	user := pkg.Cfg.Database.User
	password := pkg.Cfg.Database.Password
	dbname := pkg.Cfg.Database.DBName

	// Connect to postgres database (default database)
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to postgres database: %v", err)
	}

	defer db.Close()

	// Check if database exists
	var exists bool
	query := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)`

	err = db.QueryRow(query, dbname).Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking if database exists connection: %v", err)
	}

	if !exists {
		// Create database
		createQuery := fmt.Sprintf("CREATE DATABASE %s", dbname)
		_, err = db.Exec(createQuery)
		if err != nil {
			log.Fatalf("Error creating database: %v", err)
		}
		log.Printf("Database '%s' created successfully", dbname)
	}
}

func Connect() *sql.DB {
	// Baca dari environment variables
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// dbname := os.Getenv("DB_NAME")
	host := pkg.Cfg.Database.Host
	port := pkg.Cfg.Database.Port
	user := pkg.Cfg.Database.User
	password := pkg.Cfg.Database.Password
	dbname := pkg.Cfg.Database.DBName

	// Buat connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}
	fmt.Println("🔌 Connecting to PostgreSQL...")

	return db
}

// InitDatabase initializes database (create if not exists and connect)
func InitDatabase() *sql.DB {
	once.Do(func() {
		// First, create database if it doesn't exist
		CreateDatabaseIfNotExists()

		// Then connect to the database
		dbInstance = Connect()
	})
	return dbInstance
}

func InitTransaction() *sql.Tx {
	db := InitDatabase() // pastikan dbInstance sudah inisialisasi
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	return tx
}

func CloseDB() error {
	if dbInstance != nil {
		return dbInstance.Close()
	}
	return nil
}
