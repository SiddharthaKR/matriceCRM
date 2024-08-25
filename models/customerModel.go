package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`        // Unique identifier for the customer.
	FirstName     *string            `json:"first_name" validate:"required,min=2,max=100" bson:"first_name"`
	LastName      *string            `json:"last_name" validate:"required,min=2,max=100" bson:"last_name"`
	Email         *string            `json:"email" validate:"email,required" bson:"email"` // Unique email for the customer.
	Phone         *string            `json:"phone" validate:"required" bson:"phone"`       // Contact information.
	Company       *string            `json:"company,omitempty" bson:"company"`            // Company name if applicable.
	Status        *string            `json:"status" validate:"required,eq=LEAD|eq=CUSTOMER|eq=PROSPECT" bson:"status"` // Status of the customer.
	Notes         *string            `json:"notes,omitempty" bson:"notes"`                // Additional information or notes about the customer.
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`                // Timestamp for customer creation.
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`                // Timestamp for the last update.
	CustomerID    string             `json:"customer_id" bson:"customer_id"`              // Unique identifier for business logic.
	LastInteraction time.Time        `json:"last_interaction,omitempty" bson:"last_interaction"` // Timestamp for the last interaction with the customer.
	CompanyID       primitive.ObjectID `bson:"companyID" json:"company_id"`
	PasswordHash  *string            `json:"password" validate:"required" bson:"password"` // Hashed password for security.
	Token         *string            `json:"token,omitempty" bson:"token"`                // JWT token for session management.
	RefreshToken  *string            `json:"refresh_token,omitempty" bson:"refresh_token"`// Refresh token for extended sessions.
}
