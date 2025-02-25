package routes

import (
	"Todolist/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Определяем маршруты для задач
func TodoRoute(api fiber.Router, db *pgxpool.Pool) {
	api.Get("/", func(c *fiber.Ctx) error { return controllers.GetTodos(c, db) })         // Получить все задачи
	api.Post("/", func(c *fiber.Ctx) error { return controllers.CreateTodo(c, db) })      // Создать задачу
	api.Put("/:id", func(c *fiber.Ctx) error { return controllers.UpdateTodo(c, db) })    // Обновить задачу
	api.Delete("/:id", func(c *fiber.Ctx) error { return controllers.DeleteTodo(c, db) }) // Удалить задачу
}
