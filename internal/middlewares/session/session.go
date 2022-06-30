package session

import (
	"coursebench-backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

var store *session.Store

func init() {
	store = session.New(session.Config{Expiration: time.Hour * 24 * 2, CookieHTTPOnly: false, CookieSecure: false})
}

func GetStore() *session.Store {
	return store
}

func GetUserID(ctx *fiber.Ctx) (uint, error) {
	sess, err := store.Get(ctx)
	if err != nil {
		return 0, errors.Wrap(err, errors.InternalServerError)
	}
	t := sess.Get("user_id")
	id, ok := t.(uint)
	if !ok {
		return 0, errors.New(errors.UserNotLogin)
	}
	return id, nil
}

func GetSession(ctx *fiber.Ctx) (*session.Session, error) {
	return store.Get(ctx)
}
