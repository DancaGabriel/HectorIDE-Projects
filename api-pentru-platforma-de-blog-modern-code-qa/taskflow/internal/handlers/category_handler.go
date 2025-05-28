package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category-related API requests.
type CategoryHandler struct{}

// NewCategoryHandler creates a new CategoryHandler instance.
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

// CreateCategory handles the creation of a new category.  Requires authentication.
// @Summary Create a new category
// @Description Creates a new category with the provided name.  Authentication is required, and the user must have admin rights.
// @Tags categories
// @Accept json
// @Produce json
// @Param name body string true "Category Name"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// In a real application, check user roles/permissions for admin access.
	// Placeholder: Assuming authentication middleware has set user information.
	/*
		user, exists := c.Get("user")
		if !exists || !user.(User).IsAdmin { //Assuming User struct has an IsAdmin field
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}
	*/

	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulate category creation (replace with database interaction).
	// In a real application, validate that the category name is unique.
	// id := generateUniqueID() // Replace with a unique ID generation mechanism (UUID or similar)
	newCategoryID := 1 //Autoincrement
	category := map[string]interface{}{
		"id":   newCategoryID,
		"name": requestBody.Name,
	}

	// Simulate successful creation.
	c.JSON(http.StatusCreated, category)
}

// GetCategories retrieves a list of all categories.
// @Summary Get all categories
// @Description Retrieves a list of all categories. No authentication required.
// @Tags categories
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	// Simulate retrieving categories from the database (replace with actual database query).
	categories := []map[string]interface{}{
		{"id": 1, "name": "Technology"},
		{"id": 2, "name": "Travel"},
		{"id": 3, "name": "Food"},
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory retrieves a specific category by ID.
// @Summary Get a category by ID
// @Description Retrieves a specific category by its ID. No authentication required.
// @Tags categories
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{category_id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("category_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Simulate retrieving the category from the database (replace with actual database query).
	var category map[string]interface{}
	switch id {
	case 1:
		category = map[string]interface{}{"id": 1, "name": "Technology"}
	case 2:
		category = map[string]interface{}{"id": 2, "name": "Travel"}
	case 3:
		category = map[string]interface{}{"id": 3, "name": "Food"}
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory updates an existing category.  Requires authentication.
// @Summary Update an existing category
// @Description Updates an existing category with the provided name. Authentication is required, and the user must have admin rights.
// @Tags categories
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Param name body string true "Category Name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{category_id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	// In a real application, check user roles/permissions for admin access.
	// Placeholder: Assuming authentication middleware has set user information.
	/*
		user, exists := c.Get("user")
		if !exists || !user.(User).IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}
	*/
	idStr := c.Param("category_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulate updating the category in the database (replace with actual database update).
	// In a real application, validate that the category name is unique (excluding the current category).
	// Find current Category with id
	category := map[string]interface{}{
		"id":   id,
		"name": requestBody.Name,
	}

	// Simulate successful update.
	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category.  Requires authentication.
// @Summary Delete a category
// @Description Deletes a category by ID. Authentication is required, and the user must have admin rights.
// @Tags categories
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 204 {string} string ""
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{category_id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	// In a real application, check user roles/permissions for admin access.
	// Placeholder: Assuming authentication middleware has set user information.
	/*
		user, exists := c.Get("user")
		if !exists || !user.(User).IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}
	*/
	idStr := c.Param("category_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Simulate deleting the category from the database (replace with actual database deletion).
	// In a real application, check if the category is associated with any posts and handle accordingly.
	// (e.g., prevent deletion or re-assign the posts to a default category).

	// Simulate successful deletion.
	c.Status(http.StatusNoContent)
}