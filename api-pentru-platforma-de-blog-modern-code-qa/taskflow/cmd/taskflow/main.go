package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv" // For loading .env files
)

// Config struct to hold application configurations
type Config struct {
	GitRepoPath          string
	APIDocumentationPath string
	DatabaseBackupPath   string
}

// GitService provides git-related functionalities
type GitService struct {
	repoPath string
}

// NewGitService creates a new GitService instance.
func NewGitService(repoPath string) *GitService {
	return &GitService{repoPath: repoPath}
}

// InitializeGitRepo initializes a new Git repository if one doesn't exist.
func (g *GitService) InitializeGitRepo() error {
	if _, err := os.Stat(g.repoPath + "/.git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", g.repoPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %w", err)
		}
		fmt.Println("Initialized new Git repository at", g.repoPath)
	} else if err != nil {
		return fmt.Errorf("error checking git repository: %w", err)
	}
	return nil
}

// CommitAndPush commits changes and pushes them to a remote repository (if configured).
func (g *GitService) CommitAndPush(message string) error {
	cmdAdd := exec.Command("git", "-C", g.repoPath, "add", ".")
	if err := cmdAdd.Run(); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	cmdCommit := exec.Command("git", "-C", g.repoPath, "commit", "-m", message)
	if err := cmdCommit.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	// Attempt to push changes.  This might fail if no remote is configured.
	cmdPush := exec.Command("git", "-C", g.repoPath, "push")
	if err := cmdPush.Run(); err != nil {
		fmt.Println("Warning: Failed to push changes (no remote configured or other issues):", err)
		// We don't want to fail the entire operation if the push fails.
		// Pushing to remote could be an optional feature.
	}
	return nil
}

// APIService handles API Documentation updates.
type APIService struct {
	docPath string
}

// NewAPIService creates a new APIService instance.
func NewAPIService(docPath string) *APIService {
	return &APIService{docPath: docPath}
}

// UpdateAPIDocumentation updates the API documentation (e.g., Swagger/OpenAPI).
func (a *APIService) UpdateAPIDocumentation(newContent string) error {
	// Simple file overwrite.  Replace with more sophisticated logic
	// if the documentation format requires it (e.g., merging changes).
	err := os.WriteFile(a.docPath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update API documentation: %w", err)
	}
	return nil
}

// DataModelService handles Data Model modifications.
type DataModelService struct {
	dbBackupPath string
}

// NewDataModelService creates a new DataModelService instance.
func NewDataModelService(dbBackupPath string) *DataModelService {
	return &DataModelService{dbBackupPath: dbBackupPath}
}

// BackupDatabase creates a backup of the database.
func (d *DataModelService) BackupDatabase() error {
	// Placeholder: Replace with actual database backup command (e.g., pg_dump)
	backupFilename := fmt.Sprintf("backup_%s.dump", time.Now().Format("20060102150405"))
	backupPath := d.dbBackupPath + "/" + backupFilename

	//For Posgres:
	//cmd := exec.Command("pg_dump", "-U", "user", "-d", "database", "-f", backupPath)

	// Simulate a backup (replace with real backup command).
	file, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString("Simulated database backup data.")
	if err != nil {
		return fmt.Errorf("failed to write to backup file: %w", err)
	}

	fmt.Println("Simulated database backup created at", backupPath)
	return nil
}

// DatabaseMigrationService handles schema migrations
type DatabaseMigrationService struct{}

func NewDatabaseMigrationService() *DatabaseMigrationService {
	return &DatabaseMigrationService{}
}

// ApplyMigrations applies database migrations.
func (d *DatabaseMigrationService) ApplyMigrations() error {
	// Placeholder: Replace with actual migration tool command (e.g., alembic upgrade head)
	fmt.Println("Simulating applying database migrations.")
	return nil
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file") // Not fatal, allow defaults.
	}

	// Configuration
	config := Config{
		GitRepoPath:          os.Getenv("GIT_REPO_PATH"), //Expect the env variable GIT_REPO_PATH to exist.
		APIDocumentationPath: os.Getenv("API_DOC_PATH"),  //Expect the env variable API_DOC_PATH to exist.
		DatabaseBackupPath:   os.Getenv("DB_BACKUP_PATH"), //Expect the env variable DB_BACKUP_PATH to exist.
	}

	if config.GitRepoPath == "" {
		log.Fatal("GIT_REPO_PATH is not set. Please set it in .env or environment variables.")
		os.Exit(1) // Exit if GIT_REPO_PATH is not set
	}

	// Initialize Services
	gitService := NewGitService(config.GitRepoPath)
	apiService := NewAPIService(config.APIDocumentationPath)
	dataModelService := NewDataModelService(config.DatabaseBackupPath)
	migrationService := NewDatabaseMigrationService()

	// Initialize Git Repo
	err = gitService.InitializeGitRepo()
	if err != nil {
		log.Fatalf("Failed to initialize Git repository: %v", err)
	}

	// Gin setup
	r := gin.Default()

	// Serve frontend files (HTML, CSS, JavaScript) from the "web" directory
    r.Static("/static", "./web/static") // Serve static assets (CSS, JS, images)
    r.LoadHTMLGlob("web/templates/*") // Load HTML templates

    // Define a route to serve the index page (entry point to your frontend)
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
            "title": "TaskFlow",
            "message": "Welcome to TaskFlow!",
        })
    })

	// API Endpoints
	r.POST("/document/update", func(c *gin.Context) {
		// Request body should contain the new document content,
		// changes to API endpoints (optional), and data model changes (optional).
		var requestBody struct {
			DocumentContent     string `json:"document_content"`
			APIDocumentation    string `json:"api_documentation"`
			DataModelChanges    string `json:"data_model_changes"` //Description of changes needed
			CommitMessage       string `json:"commit_message"`
			PerformDBMigration  bool   `json:"perform_db_migration"` //Indicate if we must perform a DB migration
			UpdateAPIDocumentation bool  `json:"update_api_documentation"`//Indicate if we want to perform API documentation update
		}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a unique file name.
		filename := uuid.New().String() + ".md" // Or any relevant extension
		filePath := config.GitRepoPath + "/" + filename

		// Save new document content to file.
		err := os.WriteFile(filePath, []byte(requestBody.DocumentContent), 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save document: %v", err)})
			return
		}

		// 1. Update API Documentation (if provided and enabled)
		if requestBody.UpdateAPIDocumentation && requestBody.APIDocumentation != "" {
			if err := apiService.UpdateAPIDocumentation(requestBody.APIDocumentation); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update API documentation: %v", err)})
				return
			}
		}

		// 2. Backup Database (if Data Model Changes are requested)
		if requestBody.DataModelChanges != "" { //Simple check if we have changes to models
			if err := dataModelService.BackupDatabase(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to backup database: %v", err)})
				return
			}
		}

		// 3. Apply database migration (if Data Model Changes are requested)
		if requestBody.PerformDBMigration && requestBody.DataModelChanges != "" {
			if err := migrationService.ApplyMigrations(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to apply database migrations: %v", err)})
				return
			}
		}

		// 4. Commit and Push changes to Git Repository
		if err := gitService.CommitAndPush(requestBody.CommitMessage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to commit and push changes: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
	})

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Server listening on port %s...\n", port)
	r.Run(":" + port)
}