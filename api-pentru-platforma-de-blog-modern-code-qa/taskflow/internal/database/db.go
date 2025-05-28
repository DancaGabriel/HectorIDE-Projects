package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// InitDB initializes the database connection.
func InitDB() (*sql.DB, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: ", err)
		// Proceed without .env, assuming environment variables are set
	}

	// Retrieve database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // Default to localhost if not set
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432" // Default port if not set
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres" // Default user if not set
	}
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "taskflow" // Default database name if not set
	}

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	// Open a database connection
	var errOpen error
	db, errOpen = sql.Open("postgres", connStr)
	if errOpen != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", errOpen)
	}

	// Test the connection
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("failed to ping database: %w", errPing)
	}

	log.Println("Connected to the database")

	return db, nil
}

// GetDB returns the database connection instance.
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection not initialized. Call InitDB() first.")
	}
	return db
}

// EnsureTablesExist creates the necessary tables if they don't exist.
func EnsureTablesExist(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			registration_date TIMESTAMP WITH TIME ZONE DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author_id INTEGER REFERENCES users(id),
			publication_date TIMESTAMP WITH TIME ZONE DEFAULT now(),
			slug VARCHAR(255) UNIQUE NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	// Create comments table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER REFERENCES posts(id),
			user_id INTEGER REFERENCES users(id),
			content TEXT NOT NULL,
			created TIMESTAMP WITH TIME ZONE DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create comments table: %w", err)
	}

	// Create categories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}

	// Create post_categories table (many-to-many relationship)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER REFERENCES posts(id),
			category_id INTEGER REFERENCES categories(id),
			PRIMARY KEY (post_id, category_id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create post_categories table: %w", err)
	}

	log.Println("Tables created (if not already exist)")
	return nil
}