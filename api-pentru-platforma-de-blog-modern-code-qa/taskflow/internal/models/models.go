// Backend Code (Go)

package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the user data model.
type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username       string    `gorm:"type:varchar(50);unique;not null"`
	Email          string    `gorm:"type:varchar(255);unique;not null"`
	Password       string    `gorm:"type:varchar(255);not null"`
	FirstName      string    `gorm:"type:varchar(100)"`
	LastName       string    `gorm:"type:varchar(100)"`
	RegistrationDate time.Time `gorm:"type:timestamptz;default:now()"`
	Posts          []Post    `gorm:"foreignKey:AuthorID"` // One-to-Many relationship with Post
	Comments       []Comment `gorm:"foreignKey:AuthorID"` // One-to-Many relationship with Comment
}

// BeforeCreate hook to hash password before saving.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Post represents the post data model.
type Post struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Content        string    `gorm:"type:text;not null"`
	Title          string    `gorm:"type:varchar(255);not null"`
	AuthorID       uuid.UUID `gorm:"type:uuid;not null"` // Foreign key referencing User
	Author         User      `gorm:"foreignKey:AuthorID"` // Relationship with User
	Categories     []Category `gorm:"many2many:post_categories;"`
	PublicationDate time.Time `gorm:"type:timestamptz;default:now()"`
	Slug           string    `gorm:"type:varchar(255);unique;not null"`
	Comments       []Comment `gorm:"foreignKey:PostID"` // One-to-Many relationship with Comment
	PostCategories []PostCategory `gorm:"foreignKey:PostID"`

}

// Comment represents the comment data model.
type Comment struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Text         string    `gorm:"type:text;not null"`
	AuthorID     uuid.UUID `gorm:"type:uuid;not null"` // Foreign key referencing User
	Author       User      `gorm:"foreignKey:AuthorID"` // Relationship with User
	PostID       uuid.UUID `gorm:"type:uuid;not null"` // Foreign key referencing Post
	Post         Post      `gorm:"foreignKey:PostID"`   // Relationship with Post
	CreationDate time.Time `gorm:"type:timestamptz;default:now()"`
}

// Category represents the category data model.
type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(100);unique;not null"`
	Description string    `gorm:"type:text"`
	Posts          []Post    `gorm:"many2many:post_categories;"`
}

// PostCategory represents the many-to-many relationship between Post and Category.
type PostCategory struct {
	PostID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Post       Post      `gorm:"foreignKey:PostID"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
}

// BeforeSave hook for Post to generate slug.
func (p *Post) BeforeSave(tx *gorm.DB) (err error) {
	// Implement slug generation logic here, ensuring uniqueness.
	// This is a placeholder. Consider using a library like "github.com/gosimple/slug"
	// and adding a retry mechanism to handle collisions by appending a unique identifier.
	if p.Slug == "" {
		p.Slug = generateSlug(p.Title) // Replace with actual slug generation
		// Check for collision and append a unique identifier if necessary (not implemented here)
	}
	return nil
}

// generateSlug is a placeholder function for slug generation.
// Replace this with actual slug generation logic.
func generateSlug(title string) string {
	// Placeholder implementation: lowercase and replace spaces with hyphens.
	slug := ""
	for _, r := range title {
		if r >= 'a' && r <= 'z' {
			slug += string(r)
		} else if r >= 'A' && r <= 'Z' {
			slug += string(r + 32) // Convert to lowercase
		} else if r == ' ' {
			slug += "-"
		}
	}
	return slug
}

// Migrate performs the database migration.
// Requires a gorm.DB instance
func Migrate(db *gorm.DB) error {
	// AutoMigrate will create tables, add missing columns, and add missing indexes.
	// It will NOT delete or change existing columns due to safety.
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Category{}, &PostCategory{})
	if err != nil {
		return err
	}

	// Create indexes manually.  AutoMigrate doesn't always handle index creation perfectly, especially for composite indexes.
	// NOTE: Consider using GORM's `HasIndex` and `Migrator().CreateIndex` methods
	//       to prevent errors if the indexes already exist (for idempotent migrations).
	db.Exec("CREATE INDEX IF NOT EXISTS idx_post_author_id ON posts (author_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comment_author_id ON comments (author_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comment_post_id ON comments (post_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_post_publication_date ON posts (publication_date)")

	return nil
}