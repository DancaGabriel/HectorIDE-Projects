package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct to hold application configurations
type Config struct {
	GitRepoPath          string
	APIDocumentationPath string
	DatabaseBackupPath   string
	DatabaseHost         string
	DatabasePort         string
	DatabaseUser         string
	DatabasePass         string
	DatabaseName         string
}

// LoadConfig loads the configuration from environment variables and .env file.
func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file") // Not fatal, allow defaults.
	}

	// Default values
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "postgres"
	dbName := "mydb"

	// Override defaults with environment variables if set
	if os.Getenv("DB_HOST") != "" {
		dbHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		dbPort = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_USER") != "" {
		dbUser = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_NAME") != "" {
		dbName = os.Getenv("DB_NAME")
	}

	config := &Config{
		GitRepoPath:          os.Getenv("GIT_REPO_PATH"),
		APIDocumentationPath: os.Getenv("API_DOC_PATH"),
		DatabaseBackupPath:   os.Getenv("DB_BACKUP_PATH"),
		DatabaseHost:         dbHost,
		DatabasePort:         dbPort,
		DatabaseUser:         dbUser,
		DatabasePass:         os.Getenv("DB_PASS"),
		DatabaseName:         dbName,
	}

	if config.GitRepoPath == "" {
		log.Fatal("GIT_REPO_PATH is not set. Please set it in .env or environment variables.")
		os.Exit(1) // Exit if GIT_REPO_PATH is not set
	}

	return config, nil
}