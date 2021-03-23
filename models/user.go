package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User Model

type User struct {
	ID              primitive.ObjectID `json:"-" bson:"_id"`
	FirstName       string             `json:"first_name" bson:"first_name" validate:"required,min=2"`
	LastName        string             `json:"last_name" bson:"last_name" validate:"required,min=2"`
	Email           string             `json:"email" bson:"email" validate:"required,email"`
	Password        string             `json:"password" bson:"password" validate:"required,min=6"`
	PhoneNo         string             `json:"phone_no" bson:"phone_no" validate:"required,min=11"`
	Address         string             `json:"address" bson:"address" validate:"required"`
	Token           string             `json:"token" bson:"token"`
	RefreshToken    string             `json:"refresh_token" bson:"refresh_token"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
	UserID          string             `json:"user_id" bson:"user_id"`
	PermissionGroup string             `json:"permission_group_id" bson:"permission_group_id"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
