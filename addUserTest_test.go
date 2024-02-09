package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// Update with your actual package name
	main "task-2"
	"task-2/database"

	// students "example/evenorodd/students"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	// Initialize the database and Gin router
	database.ConnectDatabase()
	router := gin.Default()
	router.POST("/addStudent", main.AddUser)

	// Create a sample student for testing
	newStudent := main.Student{
		Name:      "John Doe",
		RollNo:    123,
		ContactNo: 12345,
		Email:     "john.doe@example.com",
		ID:        "12345",
	}

	// Convert the student struct to JSON
	payload, err := json.Marshal(newStudent)
	assert.NoError(t, err)

	// Create a POST request with the JSON payload
	req, err := http.NewRequest("POST", "/addStudent", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "\"User is successfully created.\"", resp.Body.String())
}
