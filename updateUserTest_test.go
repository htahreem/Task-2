package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	main "example/evenorodd" // Update with your actual package name
	"example/evenorodd/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateStudent(t *testing.T) {
	// Initialize the database and Gin router
	database.ConnectDatabase()
	router := gin.Default()
	router.PUT("/updateStudent/:ID", main.UpdateUser)

	// Create a sample student for testing
	updateStudent := main.Student{
		Name:      "Updated Name",
		RollNo:    456,
		ContactNo: 12345,
		Email:     "updated.email@example.com",
		ID:        "12345",
	}

	// Convert the student struct to JSON
	payload, err := json.Marshal(updateStudent)
	assert.NoError(t, err)

	// Create a PUT request with the JSON payload
	req, err := http.NewRequest("PUT", "/updateStudent/12345", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "\"User is successfully updated.\"", resp.Body.String())
}
