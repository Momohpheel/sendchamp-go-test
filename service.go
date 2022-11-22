package main

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var db = DB

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var validate = validator.New()

func ValidateStruct(user UserRequest) []*ErrorResponse {

	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateTask(task TaskRequest) []*ErrorResponse {

	var errors []*ErrorResponse
	err := validate.Struct(task)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func CreateTask(c *fiber.Ctx) error {
	task := new(TaskRequest)

	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error":  err,
			"status": false})
	}

	errors := ValidateTask(*task)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	taskData := TaskRequest{
		Name:        task.Name,
		Description: task.Description,
	}

	db.Create(&taskData)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"task": taskData, "message": "Task Created Successfully"})

}

func UpdateTask(c *fiber.Ctx) error {
	idstring := c.Params("id")
	id, _ := strconv.Atoi(idstring)
	task := new(TaskRequest)

	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error":  err,
			"status": false})
	}

	errors := ValidateTask(*task)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	taskData := map[string]interface{}{
		"name":        task.Name,
		"description": task.Description,
	}

	db.Model(&Task{}).Where("id = ?", id).Updates(&taskData)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"status": true, "task": taskData, "message": "Task Updated Successfully"})

}

func DeleteTask(c *fiber.Ctx) error {
	idstring := c.Params("id")
	id, _ := strconv.Atoi(idstring)
	db.Where("id = ?", id).Delete(&Task{})
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"status": true, "message": "Task Deleted Successfully"})
}

func Login(c *fiber.Ctx) error {
	user := new(UserRequest)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error":  err,
			"status": false})
	}

	errors := ValidateStruct(*user)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)

	}

	result := db.Where("email = ? AND password = ?", user.Email, user.Password).First(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error":  result.Error,
			"status": false})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": true,
		"user":   user,
	})
}
