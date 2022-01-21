package session

import (
	"coursebench-backend/pkg/errors"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var sessionMap map[uint]*session.Session

func init() {
	sessionMap = make(map[uint]*session.Session)
}

func SetSession(id uint, sess *session.Session) (err error) {
	sess.ID()
	if t, ok := sessionMap[id]; ok {
		err = t.Destroy()
		if err != nil {
			return errors.Wrap(err, errors.InternalServerError)
		}
	}
	sessionMap[id] = sess
	return nil
}

func GetSession(id uint) *session.Session {
	return sessionMap[id]
}
