package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`        // Unique identifier for the user.
	FirstName     *string            `json:"first_name" validate:"required,min=2,max=100" bson:"first_name"`
	LastName      *string            `json:"last_name" validate:"required,min=2,max=100" bson:"last_name"`
	Email         *string            `json:"email" validate:"email,required" bson:"email"` // Unique email for the user.
	PasswordHash  *string            `json:"password" validate:"required" bson:"password"` // Hashed password for security.
	Phone         *string            `json:"phone" validate:"required" bson:"phone"`       // Contact information.
	Company       *string            `json:"company,omitempty" bson:"company"`            // Company name if applicable.
	Status        *string            `json:"status" bson:"status"`                        // User's status, e.g., active, inactive.
	UserType      *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER|eq=MANAGER" bson:"user_type"` // Role-based access control.
	Token         *string            `json:"token,omitempty" bson:"token"`                // JWT token for session management.
	RefreshToken  *string            `json:"refresh_token,omitempty" bson:"refresh_token"`// Refresh token for extended sessions.
	Notes         *string            `json:"notes,omitempty" bson:"notes"`                // Additional information about the user.
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`                // Timestamp for user creation.
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`                // Timestamp for the last update.
	UserID        string             `json:"user_id" bson:"user_id"`                      // Unique identifier for business logic.
	LastLogin     time.Time          `json:"last_login,omitempty" bson:"last_login"`      // Timestamp for last login.
	CompanyIDs     []primitive.ObjectID `json:"company_ids"`
}
