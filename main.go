package main

import (
	"Todolist/routes"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// Инициализация базы данных
func initDB() {
	dsn := "postgres://postgres:bioroot@localhost:5432/todoapp?sslmode=disable" // Укажи правильные параметры подключения
	var err error

	// Создаем пул соединений
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Проверяем соединение с базой
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	log.Println("✅ Успешное подключение к базе данных!")
}

func main() {
	initDB()
	app := fiber.New()
	app.Use(logger.New())

	// Группируем маршруты
	todo := app.Group("/tasks")
	routes.TodoRoute(todo, db)

	// Запускаем сервер
	log.Println("🚀 Сервер запущен на http://localhost:8000")
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
