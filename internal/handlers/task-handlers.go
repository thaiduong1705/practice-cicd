package handlers

import (
	"net/http"
	"todolist/internal/database"
	"todolist/internal/models"

	"github.com/labstack/echo/v4"
)

// GetAllTasks godoc
// @Summary      Lấy danh sách tất cả task
// @Description  Trả về toàn bộ task trong hệ thống
// @Tags         Công việc
// @Produce      json
// @Success      200 {array} models.Task
// @Failure      500 {object} map[string]string
// @Router       /tasks [get]
func GetAllTasks(ctx echo.Context) error {
	db, err := database.ConnectDB()
	defer database.CloseDB(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database connection error")
	}

	var tasks []*models.Task

	if err := db.Find(&tasks).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve tasks")
	}

	return ctx.JSON(http.StatusOK, tasks)
}

// GetTaskByID godoc
// @Summary      Lấy chi tiết task
// @Description  Lấy thông tin task theo ID
// @Tags         Công việc
// @Produce      json
// @Param        id   path string true "ID task"
// @Success      200 {object} models.Task
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [get]
func GetTaskByID(ctx echo.Context) error {
	db, err := database.ConnectDB()
	defer database.CloseDB(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database connection error")
	}
	id := ctx.Param("id")
	var task models.Task
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}

	return ctx.JSON(http.StatusOK, task)
}

// CreateTask godoc
// @Summary      Tạo task mới
// @Description  Tạo một task mới với dữ liệu gửi lên
// @Tags         Công việc
// @Accept       json
// @Produce      json
// @Param        task body models.Task true "Dữ liệu task"
// @Success      201 {object} models.Task
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks [post]
func CreateTask(ctx echo.Context) error {
	db, err := database.ConnectDB()
	defer database.CloseDB(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database connection error")
	}
	var task models.Task
	if err := ctx.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := ctx.Validate(task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	if err := db.Create(&task).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create task")
	}
	return ctx.JSON(http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary      Cập nhật task
// @Description  Cập nhật thông tin task theo ID
// @Tags         Công việc
// @Accept       json
// @Produce      json
// @Param        id   path string true "ID task"
// @Param        task body models.Task true "Dữ liệu cập nhật"
// @Success      200 {object} models.Task
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [put]
func UpdateTask(ctx echo.Context) error {
	db, err := database.ConnectDB()
	defer database.CloseDB(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database connection error")
	}
	id := ctx.Param("id")
	var task models.Task
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}

	if err := ctx.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := ctx.Validate(task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	if err := db.Save(&task).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update task")
	}
	return ctx.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary      Xóa task
// @Description  Xóa task theo ID
// @Tags         Công việc
// @Produce      json
// @Param        id   path string true "ID task"
// @Success      200 {string} string "Xóa thành công"
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [delete]
func DeleteTask(ctx echo.Context) error {
	db, err := database.ConnectDB()
	defer database.CloseDB(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database connection error")
	}
	id := ctx.Param("id")
	var task models.Task
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}
	if err := db.Delete(&task).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete task")
	}
	return ctx.NoContent(http.StatusOK)
}
