package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{
		Id:        1,
		Title:     "Walk the dog 🦮",
		Completed: false,
	},
	{
		Id:        2,
		Title:     "Walk the cat 🐈",
		Completed: false,
	},
}

// get all todos
func GetTodos(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todos": todos,
		},
	})
}
