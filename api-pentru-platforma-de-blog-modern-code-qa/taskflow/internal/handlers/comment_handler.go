package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-username/taskflow/internal/models" // Replace with your actual module path
)

// CommentHandler handles comment related API requests.
type CommentHandler struct {
	// db *gorm.DB // Inject database dependency (replace with your actual DB instance)
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		// db: db, // Assign the injected database connection
	}
}

// CreateComment handles the creation of a new comment.
func (ch *CommentHandler) CreateComment(c *gin.Context) {
	// Extract post ID from URL parameters.
	postIDStr := c.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Bind the request body to a Comment struct.
	var comment models.Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the PostID of the comment to the extracted post ID.
	comment.PostID = postID

	// Basic validation.  Replace with more comprehensive validation as needed.
	if comment.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment text is required"})
		return
	}

	// In a real application, you would:
	// 1. Authenticate the user (get UserID from JWT or session)
	// 2. Validate that the post exists.
	// 3. Save the comment to the database, associating it with the post and user.
	// 4. Handle errors appropriately.

	//Simulate create.

	comment.ID = uuid.New() //Fake ID

	// Respond with the newly created comment.
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully", "comment": comment})
}

// GetComments retrieves all comments for a specific post.
func (ch *CommentHandler) GetComments(c *gin.Context) {
	// Extract post ID from URL parameters
	postIDStr := c.Param("post_id")
	postID, err := uuid.Parse(postIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// In a real application, you would:
	// 1. Query the database to retrieve all comments associated with the post ID.
	// 2. Handle errors appropriately.
	// Simulate fetching comments

	comments := []models.Comment{
		{ID: uuid.New(), PostID: postID, Text: "Great post!", AuthorID: uuid.New()},
		{ID: uuid.New(), PostID: postID, Text: "I learned a lot.", AuthorID: uuid.New()},
	}

	// Respond with the list of comments.
	c.JSON(http.StatusOK, comments)
}

// UpdateComment handles updating a specific comment.
func (ch *CommentHandler) UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Bind the request body to a Comment struct.
	var updatedComment models.Comment
	if err := c.BindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation (replace with more comprehensive validation).
	if updatedComment.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment text is required"})
		return
	}

	// In a real application, you would:
	// 1. Authenticate the user (verify they have permission to update the comment).
	// 2. Query the database to find the comment with the given ID.
	// 3. Update the comment's fields with the values from the request body.
	// 4. Handle errors appropriately.

	// Simulate database update.
	updatedComment.ID = commentID

	// Respond with the updated comment.
	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully", "comment": updatedComment})
}

// DeleteComment handles deleting a specific comment.
func (ch *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// In a real application, you would:
	// 1. Authenticate the user (verify they have permission to delete the comment).
	// 2. Query the database to find the comment with the given ID.
	// 3. Delete the comment from the database.
	// 4. Handle errors appropriately.

	// Simulate database delete
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully", "comment_id": commentID})
}

// GetComment retrieves a specific comment by ID (Less common than GetComments, but included for completeness).
func (ch *CommentHandler) GetComment(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	//In a real application, you would query the database for a specific comment
	//Simulate returning a comment

	comment := models.Comment{
		ID:      commentID,
		PostID:    uuid.New(),
		AuthorID:  uuid.New(),
		Text:    "This is a specific comment",
	}

	c.JSON(http.StatusOK, comment)

}

// GetCommentsByUser retrieves all comments for a specific user.
func (ch *CommentHandler) GetCommentsByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// In a real application, you would:
	// 1. Query the database to retrieve all comments associated with the user ID.
	// 2. Handle errors appropriately.
	// Simulate fetching comments

	comments := []models.Comment{
		{ID: uuid.New(), PostID: uuid.New(), Text: "User's comment 1", AuthorID: userID},
		{ID: uuid.New(), PostID: uuid.New(), Text: "User's comment 2", AuthorID: userID},
	}

	// Respond with the list of comments.
	c.JSON(http.StatusOK, comments)
}

// parseIntParam is a helper function to parse an integer parameter from the context.
func parseIntParam(c *gin.Context, paramName string) (int, error) {
	paramStr := c.Param(paramName)
	paramInt, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	return paramInt, nil
}