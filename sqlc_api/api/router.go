package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	middleware "sqlc_api/Middleware"
	"sqlc_api/dtos"
	"sqlc_api/jwt_code"
	"sqlc_api/todos"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DB is a global database connection pool
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	connStr := "postgres://postgres:Beyondthewall007@localhost:5432/golang_db?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return err
	}

	return nil
}
func CloseDB() {
	if err := DB.Close(); err != nil {
		fmt.Println("Error closing database connection:", err)
	} else {
		fmt.Println("Database connection closed")
	}
}

// GetQueries returns a new queries object using the global connection
func GetQueries() *todos.Queries {
	return todos.New(DB)
}

func SetupRoutes(router *gin.Engine) {
	// add all this methods with auth middleware
	router.GET("/todos", middleware.AuthMiddleware(), getTodos)
	router.POST("/todos", middleware.AuthMiddleware(), addTodo)
	router.PUT("/todos", middleware.AuthMiddleware(), updateTodo)
	router.DELETE("/todos/:id", middleware.AuthMiddleware(), deleteTodo)
	router.GET("/todos/:id", middleware.AuthMiddleware(), getTodoByID)
	router.GET("/users", middleware.AuthMiddleware(), getUsers)
	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

}

func getTodos(c *gin.Context) {
	queries := GetQueries()
	todo, err := queries.GetTodos(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	todoDTOs := convertToTodoDTO(todo)
	c.JSON(http.StatusOK, todoDTOs)
}

func convertToTodoDTO(todos []todos.Todo) []dtos.TodoDTO {
	var todoDTOs []dtos.TodoDTO
	for _, todo := range todos {
		var createdDate string
		if todo.CreatedDate.Valid {
			createdDate = todo.CreatedDate.Time.String()
		}
		var updated_date string
		if todo.UpdatedDate.Valid {
			updated_date = todo.UpdatedDate.Time.String()
		}

		todoDTO := dtos.TodoDTO{
			ID:          todo.ID,
			Todo:        todo.Task.String,
			CreatedDate: createdDate,
			CreatedBy:   int(todo.CreatedBy.Int64),
			UpdatedDate: updated_date,
		}
		todoDTOs = append(todoDTOs, todoDTO)
	}
	return todoDTOs
}

func addTodo(c *gin.Context) {
	queries := GetQueries()
	var todo dtos.CreateTodoDTO
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newTodo := todos.CreateTodoParams{
		Task:      sql.NullString{String: todo.Todo, Valid: true},
		CreatedBy: sql.NullInt64{Int64: int64(todo.CreatedBy), Valid: true},
	}
	err := queries.CreateTodo(context.Background(), newTodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *gin.Context) {
	queries := GetQueries()
	var todo dtos.UpdateTodoDTO
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	updateTodo := todos.UpdateTodoParams{
		ID:          todo.ID,
		Task:        sql.NullString{String: todo.Todo, Valid: true},
		UpdatedDate: sql.NullTime{Valid: true, Time: time.Now()},
	}
	err := queries.UpdateTodo(context.Background(), updateTodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

func deleteTodo(c *gin.Context) {
	queries := GetQueries()
	todoID := c.Param("id")

	id, err := strconv.Atoi(todoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = queries.DeleteTodo(context.Background(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
func getTodoByID(c *gin.Context) {
	queries := GetQueries()
	todoID := c.Param("id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	todo, err := queries.GetTodo(context.Background(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todo"})
		return
	}
	todoDTO := dtos.TodoDTO{
		ID:          todo.ID,
		Todo:        todo.Task.String,
		CreatedDate: todo.CreatedDate.Time.String(),
		CreatedBy:   int(todo.CreatedBy.Int64),
		UpdatedDate: todo.UpdatedDate.Time.String(),
	}

	c.JSON(http.StatusOK, todoDTO)
}

func getUsers(c *gin.Context) {
	queries := GetQueries()
	users, err := queries.GetAllUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	userDTOs := convertToUserDTO(users)
	c.JSON(http.StatusOK, userDTOs)
}
func convertToUserDTO(users []todos.User) []dtos.UserDTO {
	var userDTOs []dtos.UserDTO
	for _, user := range users {
		userDTO := dtos.UserDTO{
			ID:          user.ID,
			Username:    user.Username,
			EmailId:     user.EmailID,
			PhoneNo:     user.PhoneNo.String,
			CreatedDate: user.CreatedDate.String(),
		}
		userDTOs = append(userDTOs, userDTO)
	}
	return userDTOs
}

func registerUser(c *gin.Context) {
	queries := GetQueries()
	var user dtos.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newuser := todos.RegisterParams{
		Username: user.Username,
		EmailID:  user.EmailId,
		PhoneNo:  sql.NullString{String: user.PhoneNo, Valid: true},
		Password: user.Password,
	}
	err := queries.Register(context.Background(), newuser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func validateUser(user dtos.UserDTO) error {
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}
	if user.EmailId == "" {
		return fmt.Errorf("email is required")
	}
	if user.PhoneNo == "" {
		return fmt.Errorf("phone number is required")
	}
	if len(user.PhoneNo) != 10 {
		return fmt.Errorf("phone number must be 10 digits")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(user.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	} else {
		// Check if password contains at least one uppercase letter, one lowercase letter, one digit, and one special character
		hasUpper := false
		hasLower := false
		hasDigit := false
		for _, char := range user.Password {
			switch {
			case char >= 'A' && char <= 'Z':
				hasUpper = true
			case char >= 'a' && char <= 'z':
				hasLower = true
			case char >= '0' && char <= '9':
				hasDigit = true
			}
		}
		if !hasUpper || !hasLower || !hasDigit {
			return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, one digit")
		}
	}
	return nil
}

func loginUser(c *gin.Context) {
	queries := GetQueries()
	var user dtos.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	loginUser := todos.LoginParams{
		Username: user.Username,
		Password: user.Password,
	}
	returnUser, err := queries.Login(context.Background(), loginUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid username or password! please try again"})
		return
	}

	token, err := jwt_code.GenerateToken(int(returnUser.ID), returnUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": returnUser, "token": token})
}
