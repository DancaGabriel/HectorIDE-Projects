package services

import (
	"fmt"
	"os"
)

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
	fmt.Printf("API Documentation updated at %s\n", a.docPath)
	return nil
}