package main

import (
	"gorm.io/gorm"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	gorm.Model
}
type Task struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	gorm.Model
}
