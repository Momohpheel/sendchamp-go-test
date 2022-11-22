package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello! World")
	})

	Database()

	app.Post("/login", Login)
	task := app.Group("task")
	task.Post("/create", CreateTask)
	task.Post("/update/:id", UpdateTask)
	task.Delete("/delete/:id", DeleteTask)

	log.Fatal(app.Listen(":3000"))
}
