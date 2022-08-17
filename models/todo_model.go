package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	Task       string             `json:"task" validate:"required"`
	Done       bool               `json:"done" validate:"required"`
	Created_at time.Time          `json:"created_at" validate:"required"`
}
type TodoList struct {
	Id       primitive.D `json:"id,omitempty"`
	TodoList []Todo      `json:"todolist,omitempty"`
}
