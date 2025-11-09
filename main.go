package main

import (
	"log"
	"net"
	"todolist/internal/routes"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title           TodoList API
// @version         1.0
// @description     API TodoList với Echo, GORM và SQLite.
// @host            localhost:1323
// @BasePath
// A todolist application to practice CI/CD

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Load biến môi trường từ file .env (nếu tồn tại)
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] Không tìm thấy file .env hoặc lỗi khi load:", err)
	}
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize routes
	routes.InitRoutes(e)

	l, err := net.Listen("tcp", ":1323")
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Listener = l
	e.Logger.Fatal(e.StartServer(e.Server))
}
