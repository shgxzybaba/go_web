package security

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/utils"
	"net/http"
	"time"
)

type User struct {
	Username string
	Id       int
}

func (u *User) user(s data.Student) {
	u.Username = s.Username
	u.Id = int(s.Id)
}

func Session(w http.ResponseWriter, r *http.Request) (session data.Session, err error) {

	cookie, err := r.Cookie("session-id")
	if err == nil {
		session = data.Session{
			Uuid: cookie.Value,
		}
		if ok, _ := session.Check(); !ok {
			err = errors.New("invalid session id")
		}

	}
	return
}

func BasicSecurity(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Validating session cookie")
		_, err := Session(w, r)
		if err != nil {
			utils.BasicErrorHandle(w, r)
			return
		}
		h(w, r)

	}
}

func Login(username, password string) (student data.Student, session data.Session, err error) {
	// Execute the query
	rows := data.DB.QueryRow("SELECT id, username, password FROM student WHERE username = $1", username)
	student = data.Student{}
	err = rows.Scan(&student.Id, &student.Username, &student.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("no matching user found")
		} else {
			// Other scanning errors, return the error
			return
		}
	} else {
		if ok := validatePassword(password, student.Password); !ok {
			err = errors.New("invalid username/password combination")
			return
		}

		session, err = student.CreateSession()
	}

	return
}

func validatePassword(sent, expected string) (ok bool) {
	ok = sent == expected
	return
}
func LoginHandler(c *fiber.Ctx) (err error) {
	user, password := c.FormValue("username"), c.FormValue("password")
	student, sess, err := Login(user, password)
	var response = utils.Data{}
	if err != nil {
		return
	} else {
		u := User{}
		u.user(student)
		response.Response = u
		response.Err = ""
		cookie := &fiber.Cookie{
			Name:     "session-id",
			Value:    sess.Uuid,
			HTTPOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		}
		c.Cookie(cookie)
		return c.Redirect("/dashboard", 301)
	}

}

func GetLoginPage(c *fiber.Ctx) error {
	response := utils.Data{}
	return c.Render("login", response, "layout")
}
