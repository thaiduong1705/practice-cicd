package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	appdb "todolist/internal/database"
	"todolist/internal/handlers"
	"todolist/internal/models"
)

var (
	e   *echo.Echo
	gdb *gorm.DB
)

// Validator custom
type CustomValidator struct{ v *validator.Validate }

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

// Setup trước khi chạy test
func TestMain(m *testing.M) {
	// Đảm bảo các ENV này trùng với docker-compose hoặc docker run của bạn
	// ví dụ: docker run -e POSTGRES_USER=taskuser -e POSTGRES_PASSWORD=taskpass -e POSTGRES_DB=taskdb -p 5432:5432 postgres:15
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "abc123")
	os.Setenv("DB_NAME", "todolist")
	os.Setenv("DB_PORT", "5432")

	var err error
	gdb, err = appdb.ConnectDB()
	if err != nil {
		panic(fmt.Errorf("cannot connect db: %w", err))
	}

	e = echo.New()
	e.Validator = &CustomValidator{v: validator.New()}

	e.GET("/tasks", handlers.GetAllTasks)
	e.GET("/tasks/:id", handlers.GetTaskByID)
	e.POST("/tasks", handlers.CreateTask)
	e.PUT("/tasks/:id", handlers.UpdateTask)
	e.DELETE("/tasks/:id", handlers.DeleteTask)

	code := m.Run()
	os.Exit(code)
}

// --- CRUD + Negative Tests ---
func TestTaskCRUD_AndNegative(t *testing.T) {
	// Xóa sạch DB trước test
	_ = gdb.Exec("DELETE FROM tasks")

	// ===== CREATE =====
	create := map[string]string{"title": "Learn CI/CD", "description": "Pipeline with GitLab + Go"}
	body, _ := json.Marshal(create)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rec.Code, rec.Body.String())
	}
	var created models.Task
	_ = json.Unmarshal(rec.Body.Bytes(), &created)

	// ===== LIST =====
	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// ===== GET BY ID =====
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// ===== UPDATE =====
	up := map[string]string{"title": "Learn GitLab CI/CD"}
	body, _ = json.Marshal(up)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", created.ID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// ===== DELETE =====
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// ===== NEGATIVE CASES =====
	// GET not found
	req = httptest.NewRequest(http.MethodGet, "/tasks/999999", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	// UPDATE not found
	req = httptest.NewRequest(http.MethodPut, "/tasks/999999", bytes.NewBufferString(`{"title":"x"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	// CREATE thiếu title
	req = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"description":"no title"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
