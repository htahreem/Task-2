package main

import (
	"encoding/json"
	"example/evenorodd/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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



// .env file is as follows:
// HOST=localhost
// PORT=5432
// USER=htahreem
// DB_NAME=students
// PASSWORD=Lighttube32$
