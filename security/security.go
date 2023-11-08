package security

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/utils"
)

type User struct {
	Username string
	Id       int
}

func (u *User) user(s data.Student) {
	u.Username = s.Username
	u.Id = int(s.Id)
}

func Login(username, password string) (student data.Student, err error) {
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

	}

	return
}

func validatePassword(sent, expected string) (ok bool) {
	ok = sent == expected
	return
}
func LoginHandler(c *fiber.Ctx) (err error) {
	user, password := c.FormValue("username"), c.FormValue("password")
	student, err := Login(user, password)
	if err != nil {
		return
	} else {
		var sess *session.Session
		sess, err = data.SessionStore.Get(c)
		if err != nil {
			return err
		}
		sessionData := data.SessionData{UserId: int(student.Id), Email: student.Username}
		sess.Set("session_data", sessionData)
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Unable to store session")
		}

		return c.Redirect("/dashboard", 301)
	}

}

func LogoutHandler(c *fiber.Ctx) error {

	sess, e := data.SessionStore.Get(c)
	if e != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Attempt to logout failed, try again!")
	}
	if e = sess.Destroy(); e != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Attempt to logout failed, try again!")
	}
	return c.Redirect("/", fiber.StatusFound)
}

func GetLoginPage(c *fiber.Ctx) error {
	return c.Render("login", utils.DefaultResponse(c), "layout")
}

func init() {

	// ConfigDefault is the default config
}

func GetSessionData(c *fiber.Ctx) (sessionData data.SessionData, err error) {

	sess, e := data.SessionStore.Get(c)
	if e != nil {

		err = errors.New("unable to get session data from request")
		return
	}
	s := sess.Get("session_data")

	if s == nil {
		err = errors.New("unable to get session data from request")
		return
	}

	sessionData = s.(data.SessionData)
	return
}
