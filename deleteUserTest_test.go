package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	main "example/evenorodd" // Update with your actual package name
	"example/evenorodd/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteStudent(t *testing.T) {
	// connect to database
	database.ConnectDatabase()
	router := gin.Default()
	router.DELETE("/deleteStudent/:ID", main.DeleteStudent)

	// deleteStudent := main.Student{
	// 	Name:      "Updated Name",
	// 	RollNo:    456,
	// 	ContactNo: 12345,
	// 	Email:     "updated.email@example.com",
	// 	ID:        "12345",
	// }

	req, err := http.NewRequest("DELETE", "/deleteStudent/12345", nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}
