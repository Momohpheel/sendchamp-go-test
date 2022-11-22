package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	Database()
	Migrate()
	seed(db)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello! World")
	})

	app.Post("/login", Login)
	task := app.Group("task")
	task.Post("/create", CreateTask)
	task.Post("/update/:id", UpdateTask)
	task.Delete("/delete/:id", DeleteTask)

	log.Fatal(app.Listen(":3000"))
}
