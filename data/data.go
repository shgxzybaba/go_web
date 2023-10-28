package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
)

func Setup() (err error) {
	DB, err = sql.Open("postgres", "host=127.0.0.1 port=5432 user=akindurooluwasegun dbname=go_web01 sslmode=disable")

	if err != nil {

		fmt.Println("Could not open database", err)
	}
	err = DB.Ping()
	if err != nil {
		fmt.Println("Could not reach database", err)
	}
	return

}

func ShutDown() {
	fmt.Println("Closing database connections!")
	DB.Close()
}


func AllCourses() (courses []string) {
	return 
}