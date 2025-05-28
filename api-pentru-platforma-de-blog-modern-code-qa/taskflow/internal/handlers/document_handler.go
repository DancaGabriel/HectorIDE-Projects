package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-username/taskflow/internal/services" // Replace with your actual import path
)

// DocumentHandler handles the document update endpoint.
type DocumentHandler struct {
	gitService    *services.GitService
	apiService    *services.APIService
	dataModelService *services.DataModelService
	migrationService *services.DatabaseMigrationService
	gitRepoPath   string //Path to git repo
}

// NewDocumentHandler creates a new DocumentHandler instance.
func NewDocumentHandler(gitService *services.GitService, apiService *services.APIService, dataModelService *services.DataModelService, migrationService *services.DatabaseMigrationService, gitRepoPath string) *DocumentHandler {
	return &DocumentHandler{
		gitService:    gitService,
		apiService:    apiService,
		dataModelService: dataModelService,
		migrationService: migrationService,
		gitRepoPath:   gitRepoPath,
	}
}

// UpdateDocument handles the document update request.
func (dh *DocumentHandler) UpdateDocument(c *gin.Context) {
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
	filePath := dh.gitRepoPath + "/" + filename

	// Save new document content to file.
	err := os.WriteFile(filePath, []byte(requestBody.DocumentContent), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save document: %v", err)})
		return
	}

	// 1. Update API Documentation (if provided and enabled)
	if requestBody.UpdateAPIDocumentation && requestBody.APIDocumentation != "" {
		if err := dh.apiService.UpdateAPIDocumentation(requestBody.APIDocumentation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update API documentation: %v", err)})
			return
		}
	}

	// 2. Backup Database (if Data Model Changes are requested)
	if requestBody.DataModelChanges != "" { //Simple check if we have changes to models
		if err := dh.dataModelService.BackupDatabase(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to backup database: %v", err)})
			return
		}
	}

	// 3. Apply database migration (if Data Model Changes are requested)
	if requestBody.PerformDBMigration && requestBody.DataModelChanges != "" {
		if err := dh.migrationService.ApplyMigrations(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to apply database migrations: %v", err)})
			return
		}
	}

	// 4. Commit and Push changes to Git Repository
	if err := dh.gitService.CommitAndPush(requestBody.CommitMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to commit and push changes: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// Register registers the DocumentHandler's routes with the given Gin engine.
func (dh *DocumentHandler) Register(router *gin.Engine) {
	router.POST("/document/update", dh.UpdateDocument)
}

// Logger middleware example
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}