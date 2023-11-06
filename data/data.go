package data

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v2"
	_ "github.com/lib/pq"
	"time"
)

var (
	DB           *sql.DB
	SessionStore *session.Store
)

type SessionData struct {
	UserId int `json:"user_id"`
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

	if err != nil {
		fmt.Println("Could not create pool from database", err)
		panic(err)
	}

	SessionStore = session.New(session.Config{
		Expiration:   24 * time.Hour, // Session expiration time
		CookieSecure: true,           // Enable secure cookies (HTTPS)
		Storage: postgres.New(postgres.Config{
			//DB:       Pool,
			Table:    "sessions_v2",
			Host:     "127.0.0.1",
			Username: "akindurooluwasegun",
			Database: "go_web01",
			Port:     5432,
		}),
	})
	SessionStore.RegisterType(SessionData{})

	return

}

func ShutDown() {
	fmt.Println("Closing database connections!")
	DB.Close()
}
