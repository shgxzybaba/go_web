package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"time"
)

var (
	DB *sql.DB
)

type Session struct {
	Uuid      string
	UserId    int
	CreatedAt time.Time
}

func (session *Session) Check() (ok bool, err error) {

	rows := DB.QueryRow("SELECT * FROM sessions where uuid = $1", session.Uuid)
	if err != nil {
		return false, err // Return nil slice and error
	}

	ok = rows != nil
	return
}

func init() {
	var err error
	fmt.Println("Setting up connection to database")
	DB, err = sql.Open("postgres", "host=127.0.0.1 port=5432 user=akindurooluwasegun dbname=go_web01 sslmode=disable")

	if err != nil {

		fmt.Println("Could not open database", err)
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("Could not reach database", err)
		panic(err)
	}
	fmt.Println("Database connection successful")
	return

}

func ShutDown() {
	fmt.Println("Closing database connections!")
	DB.Close()
}

func AllCourses() (courses []string) {
	return
}
