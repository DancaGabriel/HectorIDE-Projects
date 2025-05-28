package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log" // Using zerolog for structured logging

	"taskflow/internal/models"
	"taskflow/internal/services"
)

// PostHandler handles HTTP requests related to posts.
type PostHandler struct {
	postService    *services.PostService
	userService    *services.UserService // to fetch user details, if needed
	categoryService *services.CategoryService
}

// NewPostHandler creates a new PostHandler instance.
func NewPostHandler(postService *services.PostService, userService *services.UserService, categoryService *services.CategoryService) *PostHandler {
	return &PostHandler{
		postService:    postService,
		userService:    userService,
		categoryService: categoryService,
	}
}

// CreatePost handles the creation of a new post.
// @Summary Create a new post
// @Description Creates a new post with the provided details.
// @Tags posts
// @Accept json
// @Produce json
// @Param request body models.Post true "Post details"
// @Success 201 {object} models.Post
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.BindJSON(&post); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	//Basic Validation
	if post.Title == "" || post.Content == "" || post.AuthorID == uuid.Nil {
		log.Error().Msg("Missing required fields for post creation.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields (Title, Content, AuthorID)"})
		return
	}

	// Example Validation: Verify that the author exists (if userService is available)
	_, err := h.userService.GetUserByID(post.AuthorID) // Assuming a GetUserByID method
	if err != nil {
		log.Error().Err(err).Str("author_id", post.AuthorID.String()).Msg("Author not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
		return
	}

	createdPost, err := h.postService.CreatePost(&post)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	log.Info().Str("post_id", createdPost.ID.String()).Msg("Post created successfully")
	c.JSON(http.StatusCreated, createdPost)
}

// GetPost retrieves a specific post by ID.
// @Summary Get a post by ID
// @Description Retrieves a post with the specified ID.
// @Tags posts
// @Produce json
// @Param id path string true "Post ID (UUID)" Format(uuid)
// @Success 200 {object} models.Post
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("post_id", idStr).Msg("Invalid post ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format (UUID required)"})
		return
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		log.Error().Err(err).Str("post_id", id.String()).Msg("Post not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// ListPosts retrieves a list of posts with optional pagination.
// @Summary List posts with pagination
// @Description Retrieves a list of posts with optional pagination.
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} models.Post
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 500 {object} map[string]string "Server error"
// @Router /posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		log.Error().Err(err).Str("page", pageStr).Msg("Invalid page number")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit number")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	posts, err := h.postService.ListPosts(page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// UpdatePost handles updating an existing post.
// @Summary Update a post
// @Description Updates an existing post with the provided details.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID (UUID)" Format(uuid)
// @Param request body models.Post true "Post details"
// @Success 200 {object} models.Post
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("post_id", idStr).Msg("Invalid post ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format (UUID required)"})
		return
	}

	var post models.Post
	if err := c.BindJSON(&post); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Ensure the ID in the path matches the ID in the request body.
	if post.ID != uuid.Nil && post.ID != id {
		log.Error().Str("post_id_path", id.String()).Str("post_id_body", post.ID.String()).Msg("Post ID mismatch")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID in path and body must match"})
		return
	}
	post.ID = id // Use ID from the path to ensure consistency

	updatedPost, err := h.postService.UpdatePost(&post)
	if err != nil {
		log.Error().Err(err).Str("post_id", id.String()).Msg("Failed to update post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	if updatedPost == nil { // Assuming service returns nil if not found
		log.Warn().Str("post_id", id.String()).Msg("Post not found during update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	log.Info().Str("post_id", updatedPost.ID.String()).Msg("Post updated successfully")
	c.JSON(http.StatusOK, updatedPost)
}

// DeletePost handles deleting a post by ID.
// @Summary Delete a post
// @Description Deletes a post with the specified ID.
// @Tags posts
// @Produce json
// @Param id path string true "Post ID (UUID)" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("post_id", idStr).Msg("Invalid post ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format (UUID required)"})
		return
	}

	err = h.postService.DeletePost(id)
	if err != nil {
		log.Error().Err(err).Str("post_id", id.String()).Msg("Failed to delete post")
		// Differentiate between "not found" and other errors.
		if err.Error() == "Post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})

		return
	}

	log.Info().Str("post_id", id.String()).Msg("Post deleted successfully")
	c.Status(http.StatusNoContent)
}

// GetPostsByCategory retrieves posts belonging to a specific category.
// @Summary Get posts by category ID
// @Description Retrieves a list of posts associated with a specific category ID.
// @Tags posts
// @Produce json
// @Param category_id path string true "Category ID (UUID)" Format(uuid)
// @Success 200 {array} models.Post
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /categories/{category_id}/posts [get]
func (h *PostHandler) GetPostsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		log.Error().Err(err).Str("category_id", categoryIDStr).Msg("Invalid category ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format (UUID required)"})
		return
	}

	// Optional: Verify that the category exists (if categoryService is available)
	_, err = h.categoryService.GetCategoryByID(categoryID) // Assuming a GetCategoryByID method
	if err != nil {
		log.Error().Err(err).Str("category_id", categoryID.String()).Msg("Category not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	posts, err := h.postService.GetPostsByCategory(categoryID)
	if err != nil {
		log.Error().Err(err).Str("category_id", categoryID.String()).Msg("Failed to get posts by category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts by category"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// AddCategoryToPost adds a category to a post.
// @Summary Add a category to a post
// @Description Adds a category to a specific post.
// @Tags posts
// @Produce json
// @Param post_id path string true "Post ID (UUID)" Format(uuid)
// @Param category_id path string true "Category ID (UUID)" Format(uuid)
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]string "Invalid post ID or category ID"
// @Failure 404 {object} map[string]string "Post or Category not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /posts/{post_id}/categories/{category_id} [post]
func (h *PostHandler) AddCategoryToPost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	categoryIDStr := c.Param("category_id")

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		log.Error().Err(err).Str("post_id", postIDStr).Msg("Invalid post ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format (UUID required)"})
		return
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		log.Error().Err(err).Str("category_id", categoryIDStr).Msg("Invalid category ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format (UUID required)"})
		return
	}

	// Verify that both the post and category exist (if services are available)
	_, err = h.postService.GetPostByID(postID)
	if err != nil {
		log.Error().Err(err).Str("post_id", postID.String()).Msg("Post not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	_, err = h.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		log.Error().Err(err).Str("category_id", categoryID.String()).Msg("Category not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	err = h.postService.AddCategoryToPost(postID, categoryID)
	if err != nil {
		log.Error().Err(err).Str("post_id", postID.String()).Str("category_id", categoryID.String()).Msg("Failed to add category to post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add category to post"})
		return
	}

	log.Info().Str("post_id", postID.String()).Str("category_id", categoryID.String()).Msg("Category added to post successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Category added to post successfully"})
}

// RemoveCategoryFromPost removes a category from a post.
// @Summary Remove a category from a post
// @Description Removes a category from a specific post.
// @Tags posts
// @Produce json
// @Param post_id path string true "Post ID (UUID)" Format(uuid)
// @Param category_id path string true "Category ID (UUID)" Format(uuid)
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]string "Invalid post ID or category ID"
// @Failure 404 {object} map[string]string "Post or Category not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /posts/{post_id}/categories/{category_id} [delete]
func (h *PostHandler) RemoveCategoryFromPost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	categoryIDStr := c.Param("category_id")

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		log.Error().Err(err).Str("post_id", postIDStr).Msg("Invalid post ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format (UUID required)"})
		return
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		log.Error().Err(err).Str("category_id", categoryIDStr).Msg("Invalid category ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format (UUID required)"})
		return
	}

	// Verify that both the post and category exist (if services are available) - for safety and good API design
	_, err = h.postService.GetPostByID(postID)
	if err != nil {
		log.Error().Err(err).Str("post_id", postID.String()).Msg("Post not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	_, err = h.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		log.Error().Err(err).Str("category_id", categoryID.String()).Msg("Category not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	err = h.postService.RemoveCategoryFromPost(postID, categoryID)
	if err != nil {
		log.Error().Err(err).Str("post_id", postID.String()).Str("category_id", categoryID.String()).Msg("Failed to remove category from post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove category from post"})
		return
	}

	log.Info().Str("post_id", postID.String()).Str("category_id", categoryID.String()).Msg("Category removed from post successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Category removed from post successfully"})
}

// AddRoutes defines the routes for the PostHandler. This is an alternative way
// to set up routes, especially useful when you want to group routes under a specific path.
func (h *PostHandler) AddRoutes(router *gin.RouterGroup) {
	router.POST("", h.CreatePost)
	router.GET("", h.ListPosts)
	router.GET("/:id", h.GetPost)
	router.PUT("/:id", h.UpdatePost)
	router.DELETE("/:id", h.DeletePost)
	router.GET("/categories/:category_id", h.GetPostsByCategory) // Get posts by category

	// Add/Remove category from post
	router.POST("/:post_id/categories/:category_id", h.AddCategoryToPost)
	router.DELETE("/:post_id/categories/:category_id", h.RemoveCategoryFromPost)

	fmt.Println("Post routes registered.") //Indicate post route handler has been called
}