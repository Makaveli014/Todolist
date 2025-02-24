package routes

import (
	"github.com/Makaveli014/Todolist/controllers"
	"github.com/gofiber/fiber/v2"
)

func TodoRoute(route fiber.Router) {
	route.Get("", controllers.GetTodos)
}
