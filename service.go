package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"
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
	dbs := DB
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

	taskData := Task{
		Name:        task.Name,
		Description: task.Description,
	}

	// taskData := map[string]interface{}{
	// 	"name":        task.Name,
	// 	"description": task.Description,
	// }

	jsonReq, _ := json.Marshal(taskData)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Task",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(q)

	err = ch.Publish(
		"",
		"Task",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(jsonReq),
		},
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Success Publishing Message to Queue")
	dbs.Create(&taskData)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"status": true, "task": taskData, "message": "Task Created Successfully"})

}

func UpdateTask(c *fiber.Ctx) error {
	dbs := DB
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

	dbs.Model(&Task{}).Where("id = ?", id).Updates(&taskData)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"status": true, "task": taskData, "message": "Task Updated Successfully"})

}

func DeleteTask(c *fiber.Ctx) error {
	dbs := DB
	idstring := c.Params("id")
	id, _ := strconv.Atoi(idstring)
	dbs.Where("id = ?", id).Delete(&Task{})
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"status": true, "message": "Task Deleted Successfully"})
}

func Login(c *fiber.Ctx) error {
	model := new(User)
	dbs := DB
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

	result := dbs.Where("email = ?", user.Email).First(&model)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error":  result.Error.Error(),
			"status": false})
	}

	status, msg := VerifyPassword(model.Password, user.Password)

	if status {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"status":  true,
			"message": "User logged in Successfully",
			"user":    user.Email,
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status":  false,
			"message": msg,
		})
	}

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Password is incorrect"
		check = false
	}

	return check, msg
}
