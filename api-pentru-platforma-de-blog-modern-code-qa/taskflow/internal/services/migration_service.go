package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// DatabaseMigrationService handles schema migrations
type DatabaseMigrationService struct {
	migrationToolPath string // Path to the migration tool executable (e.g., Flyway, Liquibase)
	databaseURL       string // Connection string for the database
}

// NewDatabaseMigrationService creates a new DatabaseMigrationService instance.
// It retrieves the migration tool path and database URL from environment variables.
func NewDatabaseMigrationService() *DatabaseMigrationService {
	migrationToolPath := os.Getenv("MIGRATION_TOOL_PATH")
	if migrationToolPath == "" {
		log.Println("Warning: MIGRATION_TOOL_PATH is not set. Using default 'flyway'. Ensure Flyway is in your PATH.")
		migrationToolPath = "flyway" // Default value if not set
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Error: DATABASE_URL is not set. Please set it in .env or environment variables.")
	}

	return &DatabaseMigrationService{
		migrationToolPath: migrationToolPath,
		databaseURL:       databaseURL,
	}
}

// ApplyMigrations applies database migrations using the configured migration tool.
// It assumes that the migration files are located in a default location (e.g., "migrations" directory).
// It's designed to be tool-agnostic by accepting the migration tool path as a configuration.
func (d *DatabaseMigrationService) ApplyMigrations() error {
	// Construct the command to execute the migration tool.
	// This example uses Flyway. Adapt the command based on the tool used (e.g., Liquibase).
	cmd := exec.Command(
		d.migrationToolPath,
		"-url="+d.databaseURL,
		"migrate",
	)

	// Set environment variables for the migration command (if needed).
	// Example: setting the Flyway user and password.
	// cmd.Env = append(os.Environ(), "FLYWAY_USER=myuser", "FLYWAY_PASSWORD=mypassword")

	// Capture the output of the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply database migrations: %w\nOutput: %s", err, string(output))
	}

	// Log the output of the migration tool.
	log.Println("Database migrations applied successfully:\n", string(output))

	return nil
}

// Info displays the migration status. This function can be extended to provide detailed migration information.
func (d *DatabaseMigrationService) Info() error {
    cmd := exec.Command(
        d.migrationToolPath,
        "-url="+d.databaseURL,
        "info",
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to get migration info: %w\nOutput: %s", err, string(output))
    }

    log.Println("Migration Info:\n", string(output))
    return nil
}

// Repair attempts to repair the migration history table.
func (d *DatabaseMigrationService) Repair() error {
    cmd := exec.Command(
        d.migrationToolPath,
        "-url="+d.databaseURL,
        "repair",
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to repair migration history: %w\nOutput: %s", err, string(output))
    }

    log.Println("Migration history repaired successfully:\n", string(output))
    return nil
}

// Validate validates the migration files against the database.
func (d *DatabaseMigrationService) Validate() error {
    cmd := exec.Command(
        d.migrationToolPath,
        "-url="+d.databaseURL,
        "validate",
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to validate migrations: %w\nOutput: %s", err, string(output))
    }

    log.Println("Migrations validated successfully:\n", string(output))
    return nil
}