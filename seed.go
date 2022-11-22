package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seed(db *gorm.DB) {

	SeedData(db)

}

func SeedData(db *gorm.DB) {

	database := DB

	user := User{

		Email:    "philip@sendchamp.com",
		Password: HashPassword("philip"),
	}

	database.FirstOrCreate(&User{}, user)

}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}
