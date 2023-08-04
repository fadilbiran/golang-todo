package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Finished_at string `json:"finished_at"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Deleted_at  string `json:"deleted_at"`
}

var todos = []todo{}

func generateID() string {
	largestID := 0
	for _, todo := range todos {
		id, _ := strconv.Atoi(todo.ID)
		if id > largestID {
			largestID = id
		}
	}

	// Generate a new unique ID based on the largest ID found
	newID := strconv.Itoa(largestID + 1)

	return newID
}

func createTodo(c *gin.Context) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	now := time.Now().Format("02-1-2006 15:04:05")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data provided"})
		return
	}

	newTodo := todo{
		ID:          generateID(), // You can implement your own ID generation logic
		Title:       input.Title,
		Description: input.Description,
		Created_at:  now,
	}

	// Save the todo to the database or slice (todos) here
	todos = append(todos, newTodo)

	c.JSON(http.StatusCreated, newTodo)
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func getTodo(c *gin.Context) {
	todoID := c.Param("id")

	// Find the todo by ID in the todos slice
	var foundTodo todo
	for i := range todos {
		if todos[i].ID == todoID {
			foundTodo = todos[i]
			break
		}
	}

	// Check if the todo is found
	if foundTodo.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, foundTodo)
}

func updateTodo(c *gin.Context) {
	todoID := c.Param("id")

	// Find the todo by ID in the todos slice
	var updatedTodo todo
	for i := range todos {
		if todos[i].ID == todoID {
			updatedTodo = todos[i]
			break
		}
	}

	// Check if the todo is found
	if updatedTodo.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	// Bind the JSON request body to the updatedTodo variable
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data provided"})
		return
	}

	// Update the todo in the todos slice
	for i := range todos {
		if todos[i].ID == todoID {
			todos[i].Title = input.Title
			todos[i].Description = input.Description
			now := time.Now().Format("02-1-2006 15:04:05")
			todos[i].Updated_at = now
			updatedTodo = todos[i]
			break
		}
	}

	c.JSON(http.StatusOK, updatedTodo)
}

func finishTodo(c *gin.Context) {
	todoID := c.Param("id")
	var finishedTodo todo
	for i := range todos {
		if todos[i].ID == todoID {
			finishedTodo = todos[i]
			break
		}
	}
	if finishedTodo.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	now := time.Now().Format("02-1-2006 15:04:05")
	finishedTodo.Finished_at = now

	for i := range todos {
		if todos[i].ID == todoID {
			todos[i] = finishedTodo
			break
		}
	}

	c.JSON(http.StatusOK, finishedTodo)

}

func deleteTodo(c *gin.Context) {
	todoID := c.Param("id")
	var deletedTodo todo
	for i := range todos {
		if todos[i].ID == todoID {
			now := time.Now().Format("02-1-2006 15:04:05")
			todos[i].Deleted_at = now
			deletedTodo = todos[i]
			break
		}
	}
	if deletedTodo.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not Found"})
	}

	c.JSON(http.StatusOK, deletedTodo)
}

func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PUT("/todos/:id", updateTodo)
	router.POST("/todos", createTodo)
	router.PATCH("/todos/:id/finish", finishTodo)
	router.PATCH("/todos/:id/delete", deleteTodo)
	router.Run("localhost:8080")
}
