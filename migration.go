package main

import "log"

func Migrate() {
	db := DB

	err := db.AutoMigrate(
		&Task{},
		&User{},
	)

	if err != nil {
		log.Fatal(err)
	}
}
