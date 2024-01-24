package main

import (
	"encoding/json"
	"example/evenorodd/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type User struct {
// 	Username string
// 	Password string
// }

type Student struct {
	Name      string `json:"name"`
	RollNo    int    `json:"rollno"`
	ContactNo int    `json:"contactno"`
	Email     string `json:"email"`
	ID        string `json:"id"`
}

func deleteStudent(ctx *gin.Context) {
	ID := ctx.Param("ID")

	tx, err := database.Db.Begin()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	// _, err = tx.Exec("DELETE FROM students WHERE \"ID\" = $1", ID)
	// if err != nil {
	// 	tx.Rollback()
	// 	fmt.Println(err)
	// 	ctx.AbortWithStatusJSON(400, "Couldn't delete user.")
	// 	return
	// }
	_, err = database.Db.Exec("DELETE FROM students WHERE \"ID\" = $1", ID)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't delete user.")
		return
	}

	// stu := Student{}
	// fmt.Println(stu, ID, ctx.Params)

	err = tx.Commit()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	ctx.JSON(http.StatusOK, "User is successfully deleted.")
}

func getAllStudents(ctx *gin.Context) {
	rows, err := database.Db.Query("SELECT * FROM students")
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var stu Student
		err := rows.Scan(&stu.Name, &stu.RollNo, &stu.ContactNo, &stu.Email, &stu.ID)
		if err != nil {
			ctx.AbortWithStatusJSON(500, "Internal Server Error")
			return
		}
		students = append(students, stu)
	}

	ctx.JSON(http.StatusOK, students)
}

func addUser(ctx *gin.Context) {
	stu := Student{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return
	}
	err = json.Unmarshal(data, &stu)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	//use Exec whenever we want to insert update or delete
	//Doing Exec(query) will not use a prepared statement, so lesser TCP calls to the SQL server

	_, err = database.Db.Exec(`INSERT INTO students VALUES (ARRAY[$1], $2, $3, $4, $5)`, stu.Name, stu.RollNo, stu.ContactNo, stu.Email, stu.ID)

	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create new user.")
	} else {
		ctx.JSON(http.StatusOK, "User is successfully created.")
	}
}

func updateUser(ctx *gin.Context) {
	// fmt.Println(ctx)
	ID := ctx.Param("ID") // Assuming you want to update a student based on their ID
	stu := Student{}
	data, err := ctx.GetRawData()
	fmt.Println(string(data))
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User data is not defined")
		return
	}
	err = json.Unmarshal(data, &stu)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}

	// fmt.Println(stu.Name)

	// Use a transaction to ensure atomicity
	tx, err := database.Db.Begin()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	_, err = tx.Exec(`
        UPDATE students 
        SET "Name" = ARRAY[$1], 
            "Roll No" = $2, 
            "Conatct No." = $3, 
            "Email" = $4 
        WHERE "ID" = $5`,
		// pq.Array([]string{stu.Name}), stu.RollNo, stu.ContactNo, stu.Email, ID)
		stu.Name, stu.RollNo, stu.ContactNo, stu.Email, ID)

	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't update user.")
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	ctx.JSON(http.StatusOK, "User is successfully updated.")
}

func main() {
	route := gin.Default()
	database.ConnectDatabase()
	route.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	route.GET("/getAllStudents", getAllStudents)
	route.POST("/addStudent", addUser)
	route.PUT("/updateStudent/:ID", updateUser)
	route.DELETE("/deleteStudent/:ID", deleteStudent)

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}

}

// package main

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type Student struct {
// 	ID        string `json:"id"`
// 	Name      string `json:"name"`
// 	RollNo    int    `json:"rollno"`
// 	ContactNo int    `json:"contactno"`
// 	Email     string `json:"email"`
// }

// var students = []Student{
// 	{ID: "1", Name: "John", RollNo: 1, ContactNo: 1234, Email: "john@gmail.com"},
// 	{ID: "2", Name: "Alice", RollNo: 2, ContactNo: 2245, Email: "alice@gmail.com"},
// 	{ID: "3", Name: "Bob", RollNo: 3, ContactNo: 3566, Email: "bob@gmail.com"},
// 	{ID: "4", Name: "Gwen", RollNo: 4, ContactNo: 7654, Email: "gwen@gmail.com"},
// }

// func getStudents(context *gin.Context) {
// 	context.IndentedJSON(http.StatusOK, students)
// }

// func addStudent(context *gin.Context) {
// 	var newStudent Student
// 	if err := context.BindJSON(&newStudent); err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	students = append(students, newStudent)
// 	context.IndentedJSON(http.StatusCreated, students)
// }

// func getStudentByID(id string) (*Student, error) {
// 	for ind, val := range students {
// 		if val.ID == id {
// 			return &students[ind], nil
// 		}
// 	}
// 	return nil, errors.New("Student doesn't exist")
// }

// func getStudent(context *gin.Context) {
// 	id := context.Param("id")
// 	student, err := getStudentByID(id)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student not found"})
// 		return
// 	}

// 	context.IndentedJSON(http.StatusOK, student)
// }

// func updateStudent(context *gin.Context) {
// 	var newStudent Student
// 	if err := context.BindJSON(&newStudent); err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id := context.Param("id")
// 	currStudent, err := getStudentByID(id)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student not found"})
// 		return
// 	}

// 	currStudent.Name = newStudent.Name
// 	currStudent.RollNo = newStudent.RollNo
// 	currStudent.ContactNo = newStudent.ContactNo
// 	currStudent.Email = newStudent.Email

// 	context.IndentedJSON(http.StatusOK, currStudent)
// }

// func main() {
// 	router := gin.Default()
// 	router.GET("/students", getStudents)
// 	router.GET("/students/:id", getStudent)
// 	router.POST("/students", addStudent)
// 	router.PATCH("/students/:id", updateStudent)
// 	router.Run("localhost:9090")
// }

// package main

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/lib/pq"
// )

// type Student struct {
// 	ID        string `json:"id"`
// 	Name      string `json:"name"`
// 	RollNo    int    `json:"rollno"`
// 	ContactNo int    `json:"contactno"`
// 	Email     string `json:"email"`
// }

// var db *sql.DB

// func initDB() {
// 	var err error
// 	// Replace the placeholders with your PostgreSQL database credentials
// 	db, err = sql.Open("postgres", "user=htahreem dbname=htahreem sslmode=disable")
// 	if err != nil {
// 		fmt.Println("Error opening database:", err)
// 		return
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println("Error connecting to database:", err)
// 		return
// 	}
// 	fmt.Println("Connected to PostgreSQL database")
// }

// func getStudents(context *gin.Context) {
// 	rows, err := db.Query("SELECT * FROM students")
// 	if err != nil {
// 		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	var students []Student
// 	for rows.Next() {
// 		var student Student
// 		err := rows.Scan(&student.ID, &student.Name, &student.RollNo, &student.ContactNo, &student.Email)
// 		if err != nil {
// 			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		students = append(students, student)
// 	}

// 	context.IndentedJSON(http.StatusOK, students)
// }

// func addStudent(context *gin.Context) {
// 	var newStudent Student
// 	if err := context.BindJSON(&newStudent); err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	_, err := db.Exec("INSERT INTO students (id, name, rollno, contactno, email) VALUES ($1, $2, $3, $4, $5)",
// 		newStudent.ID, newStudent.Name, newStudent.RollNo, newStudent.ContactNo, newStudent.Email)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	context.IndentedJSON(http.StatusCreated, newStudent)
// }

// func getStudentByID(id string) (*Student, error) {
// 	var student Student
// 	err := db.QueryRow("SELECT * FROM students WHERE id = $1", id).
// 		Scan(&student.ID, &student.Name, &student.RollNo, &student.ContactNo, &student.Email)

// 	if err != nil {
// 		return nil, errors.New("Student doesn't exist")
// 	}

// 	return &student, nil
// }

// func getStudent(context *gin.Context) {
// 	id := context.Param("id")
// 	student, err := getStudentByID(id)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student not found"})
// 		return
// 	}

// 	context.IndentedJSON(http.StatusOK, student)
// }

// func updateStudent(context *gin.Context) {
// 	var newStudent Student
// 	if err := context.BindJSON(&newStudent); err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id := context.Param("id")
// 	_, err := db.Exec("UPDATE students SET name=$1, rollno=$2, contactno=$3, email=$4 WHERE id=$5",
// 		newStudent.Name, newStudent.RollNo, newStudent.ContactNo, newStudent.Email, id)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	newStudent.ID = id
// 	context.IndentedJSON(http.StatusOK, newStudent)
// }

// func main() {
// 	initDB()

// 	router := gin.Default()
// 	router.GET("/students", getStudents)
// 	router.GET("/students/:id", getStudent)
// 	router.POST("/students", addStudent)
// 	router.PATCH("/students/:id", updateStudent)
// 	router.
