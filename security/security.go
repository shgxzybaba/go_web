package security

import (
	"errors"
	"net/http"

	"github.com/shgxzybaba/go_web01/data"
)

func Session(w http.ResponseWriter, r *http.Request) (session data.Session, err error) {

	cookie, err := r.Cookie("_sessionId")
	if err != nil {
		session = data.Session{
			Uuid: cookie.Value,
		}
		if ok, _ := session.Check(); !ok {
			err = errors.New("invalid session id")
		}

	}
	return
}

func Login(username, password string) (student data.Student, session data.Session, err error) {
	rows, err := data.DB.Query("SELECT id, username, password FROM student WHERE username = $1", username)
	student = data.Student{}
	if err != nil {
		if err = rows.Scan(&student.Id, &student.Username, &student.Password); err != nil {
			if ok := validatePassword(password, student.Password); !ok {
				err = errors.New("invalid username/password combination")
				return
			}
			session, err = student.CreateSession()
			return
		}
	}
	return
}

func validatePassword(sent, expected string) (ok bool) {
	ok = sent == expected
	return
}
