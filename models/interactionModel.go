package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Interaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID   primitive.ObjectID `bson:"customerID" json:"customer_id"`
	UserID       primitive.ObjectID `bson:"userID" json:"user_id"`
	CompanyID    primitive.ObjectID `bson:"companyID" json:"company_id"`
	Type         string             `bson:"type" json:"type" validate:"required,eq=MEETING|eq=TICKET"`
	Status       string             `bson:"status" json:"status" validate:"required,eq=OPEN|eq=RESOLVED"`
	Description  string             `bson:"description,omitempty" json:"description"`
	ScheduledAt  time.Time          `bson:"scheduled_at,omitempty" json:"scheduled_at"` // For meetings
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
