package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // don't forget to add it. It doesn't be added automatically
)

var Db *sql.DB //created outside to make it global.

// make sure your function start with uppercase to call outside of the directory.
func ConnectDatabase() {

	err := godotenv.Load() //by default, it is .env so we don't have to write
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	//we read our .env file
	// the env file should have the following layout:
	// DB_HOST=localhost
	// DB_PORT=5432
	// DB_USER=your_username
	// DB_PASSWORD=your_password
	// DB_NAME=your_database_name
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // don't forget to convert int since port is int type.
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	fmt.Print(psqlSetup)

	var err2 error
	db, err2 := sql.Open("postgres", psqlSetup)
	if err2 != nil {
		log.Fatal(err2)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	Db = db

	// db, errSql := sql.Open("postgres", psqlSetup)
	// if errSql != nil {
	// 	fmt.Println("There is an error while connecting to the database ", err)
	// 	panic(err)
	// } else {
	// Db = db
	// 	fmt.Println("Successfully connected to database!")
	// }
}

func RunMigration() {

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS students (
			name VARCHAR(255),
			roll_no INTEGER,
			contact_no INTEGER,
			email VARCHAR(255),
			id VARCHAR(255) PRIMARY KEY
		);
	`

	_, err := Db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error creating 'students' table:", err)
	}
	fmt.Println("Migration successful!")
}
