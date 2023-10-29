package security

import (
	"errors"
	"net/http"

	"github.com/shgxzybaba/go_web01/data"
)



func session(w http.ResponseWriter, r *http.Request) (session data.Session, err error) {

	cookie, err := r.Cookie("_sessionId")
	if err != nil {
		session = data.Session{
			Uuid: cookie.Value,
		}
		if ok, _ := session.Check(); !ok {
			err = errors.New("Invalid session id")
		}

	}
	return
}
