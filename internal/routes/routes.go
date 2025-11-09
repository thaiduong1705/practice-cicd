package routes

import (
	"todolist/internal/handlers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "todolist/docs" // Import docs để load swagger spec
)

func InitRoutes(e *echo.Echo) {
	// Define your routes here
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	e.GET("/tasks", handlers.GetAllTasks)

	e.GET("/tasks/:id", handlers.GetTaskByID)

	e.POST("/tasks", handlers.CreateTask)

	e.PUT("/tasks/:id", handlers.UpdateTask)

	e.DELETE("/tasks/:id", handlers.DeleteTask)

	// Swagger UI - cấu hình URL để tránh double slash
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.URL("swagger/doc.json")))
}
