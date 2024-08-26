package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name" validate:"required"`
    Email       string             `bson:"email" json:"email" validate:"required,email"`
    Phone       string             `bson:"phone,omitempty" json:"phone"`
    CompanyID   primitive.ObjectID `bson:"company_id" json:"company_id" validate:"required"`
    Status      string             `bson:"status" json:"status" validate:"required,eq=NEW|eq=CONTACTED|eq=QUALIFIED|eq=CONVERTED|eq=LOST"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
    Notes       string             `bson:"notes,omitempty" json:"notes"`
}
