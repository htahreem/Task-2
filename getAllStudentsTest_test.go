package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	main "task-2" // Update with your actual package name
	"task-2/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllStudents(t *testing.T) {
	// Initialize the database and Gin router
	database.ConnectDatabase()
	router := gin.Default()
	router.GET("/getAllStudents", main.GetAllStudents)

	// Create a GET request
	req, err := http.NewRequest("GET", "/getAllStudents", nil)
	assert.NoError(t, err)

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, resp.Code)

	// Unmarshal the JSON response into a slice of Student structs
	var students []main.Student
	err = json.Unmarshal(resp.Body.Bytes(), &students)
	assert.NoError(t, err)

	// Perform additional assertions based on your expected data
	// For example, you can check the length of the returned students slice
	assert.NotEmpty(t, students)
}
