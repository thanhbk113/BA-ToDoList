package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id"`
	Name          *string            `json:"name" validate:"required,min=2,max=10"`
	Email         *string            `json:"email" validate:"required,email"`
	Password      *string            `json:"password"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"` //related SQL
}

type Todo struct {
	Id         primitive.ObjectID `json:"id"`
	Task       string             `json:"task" validate:"required"`
	Done       bool               `json:"done"`
	Created_at time.Time          `json:"created_at"`
}
type TodoList struct {
	User_id    string    `json:"user_id"` //related SQL
	Created_at time.Time `json:"created_at"`
}
