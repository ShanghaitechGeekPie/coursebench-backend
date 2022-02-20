package fiber

import (
	"coursebench-backend/pkg/events"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var maxByteSize int

func init() {
	/*
		var err error
		maxByteSize, err = strconv.Atoi(os.Getenv("MAX_LOG_DATA_STRING_SIZE"))
		if err != nil {
			panic(err)
		}*/
}

func errorHandler(c *fiber.Ctx, err error) error {

	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	var event *events.AttributedEvent

	if err == fiber.ErrMethodNotAllowed {
		event = events.Wrap(err, events.InvalidRequest)
	}

	attributedError, ok := err.(*events.AttributedError)

	if !ok {
		event = events.Wrap(err, events.UnCaughtError)
	} else {
		event = attributedError.Event
	}

	if event == nil {
		event = events.New(events.LogicError).Add(zap.String("message", "event is nil, which should be filtered in context middleware"))
	}

	// TODO: Log error
	/*
		wrappedLog := wrapLog(c, userError)
		if wrappedLog != nil {
			logs.GetLogger().Log(wrappedLog)
		}*/

	// Return status code with error message
	return c.Status(event.HttpStatus()).JSON(models.ErrorResponse{
		Timestamp: event.Time(),
		Errno:     event.Name(),
		Message:   event.Message(),
		Error:     true,
	})
}

func trimStringWithMaxLength(origin []byte) string {
	length := len(origin)
	if length > maxByteSize {
		length = maxByteSize
	}
	return string(origin[:length])
}

/*
func wrapLog(c *fiber.Ctx, ue errors.UserError) *logs.LogWrapper {
	switch ue.LogLevel() {
	case errors.Info:
		return infoLogWrapper(c, ue)
	case errors.Error:
		fallthrough
	case errors.Fatal:
		return errorLogWrapper(c, ue)
	default:
		return nil
	}
}

func errorLogWrapper(c *fiber.Ctx, ue errors.UserError) *logs.LogWrapper {
	return logs.Wrap(ue, map[string]interface{}{
		"body":   trimStringWithMaxLength(utils.CopyBytes(c.Body())),
		"header": trimStringWithMaxLength(utils.CopyBytes(c.Request().Header.Header())),
		"url":    utils.CopyString(c.OriginalURL()),
		"ip":     c.IP(),
		"method": c.Method(),
	})
}

func infoLogWrapper(c *fiber.Ctx, ue errors.UserError) *logs.LogWrapper {
	return logs.Wrap(ue, map[string]interface{}{
		"ip":     c.IP(),
		"header": trimStringWithMaxLength(utils.CopyBytes(c.Request().Header.Header())),
	})
}
*/
