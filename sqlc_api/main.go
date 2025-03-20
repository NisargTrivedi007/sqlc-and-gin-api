package main

import (
	"context"
	"database/sql"
	"fmt"
	"sqlc_api/api"
	"sqlc_api/todos"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	connStr := "postgres://postgres:Beyondthewall007@localhost:5432/golang_db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	context := context.Background()
	queries := todos.New(db)

	todo, err := queries.GetTodos(context)
	if err != nil {
		panic(err)
	}
	fmt.Println(todo)

	api.SetupRoutes(router)
	// router.GET("/user", homePage)
	router.Run()
}
