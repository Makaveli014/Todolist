package controllers

import (
	"Todolist/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strconv"
)

// Получить все задачи
func GetTodos(c *fiber.Ctx, db *pgxpool.Pool) error {
	rows, err := db.Query(context.Background(), "SELECT id, title, description, status, created_at, updated_at FROM tasks")
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Ошибка при получении задач",
			"error":   err.Error(),
		})
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Printf("Error scanning task data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Ошибка при сканировании данных",
				"error":   err.Error(),
			})
		}
		todos = append(todos, todo)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    todos,
	})
}

// Создать задачу
// Создать задачу
func CreateTodo(c *fiber.Ctx, db *pgxpool.Pool) error {
	// Логируем начало запроса
	log.Println("Получен POST-запрос на создание задачи")

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	// Парсим JSON из тела запроса
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Некорректный JSON",
			"error":   err.Error(),
		})
	}

	// Логируем входные данные
	log.Printf("Попытка вставки: Title=%s, Description=%s, Status=%s", body.Title, body.Description, body.Status)

	// Проверяем обязательные поля
	if body.Title == "" || body.Description == "" {
		log.Println("Ошибка: пустые поля title или description")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Поля title и description обязательны",
		})
	}

	// Если статус пустой, устанавливаем значение по умолчанию
	if body.Status == "" {
		body.Status = "new"
	}

	// Выполняем SQL-запрос на вставку
	_, err := db.Exec(context.Background(),
		"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)",
		body.Title, body.Description, body.Status)

	if err != nil {
		log.Printf("Ошибка при вставке в БД: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Ошибка при создании задачи",
			"error":   err.Error(),
		})
	}

	log.Println("Задача успешно добавлена!")

	// Возвращаем успешный ответ
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Задача создана",
	})
}

// Обновить задачу
func UpdateTodo(c *fiber.Ctx, db *pgxpool.Pool) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID должен быть числом",
		})
	}

	var body struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
	}

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Некорректный JSON",
			"error":   err.Error(),
		})
	}

	// Логирование данных перед обновлением
	log.Printf("Updating task ID %d with title=%s, description=%s, status=%s", id, *body.Title, *body.Description, *body.Status)

	_, err = db.Exec(context.Background(),
		"UPDATE tasks SET title = COALESCE($1, title), description = COALESCE($2, description), status = COALESCE($3, status), updated_at = now() WHERE id = $4",
		body.Title, body.Description, body.Status, id)
	if err != nil {
		log.Printf("Error updating task ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Ошибка при обновлении задачи",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Задача обновлена",
	})
}

// Удалить задачу
func DeleteTodo(c *fiber.Ctx, db *pgxpool.Pool) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID должен быть числом",
		})
	}

	// Логирование данных перед удалением
	log.Printf("Deleting task ID %d", id)

	_, err = db.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		log.Printf("Error deleting task ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Ошибка при удалении задачи",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Задача удалена",
	})
}
