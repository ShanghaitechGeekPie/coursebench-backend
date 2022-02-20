package session

import (
	"coursebench-backend/pkg/events"
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

func GetUserID(ctx *fiber.Ctx) (uint, *events.AttributedEvent) {
	sess, err := store.Get(ctx)
	if err != nil {
		return 0, events.Wrap(err, events.InternalServerError)
	}
	t := sess.Get("user_id")
	id, ok := t.(uint)
	if !ok {
		return 0, events.New(events.UserNotLogin)
	}
	return id, nil
}
