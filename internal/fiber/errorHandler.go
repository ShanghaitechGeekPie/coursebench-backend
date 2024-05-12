// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

package fiber

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
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

	if err == fiber.ErrMethodNotAllowed {
		err = errors.Wrap(err, errors.InvalidRequest)
	}

	userError, ok := err.(errors.UserError)

	if !ok {
		userError = errors.Wrap(err, errors.UnCaughtError).(errors.UserError)
	}
	errorMsg := fmt.Sprintf("%s, %s, %s, %s", userError.Error(), userError.Name(), userError.Cause().Error(), userError.Stacktrace())
	log.Println(errorMsg)

	// TODO: Log error
	/*
		wrappedLog := wrapLog(c, userError)
		if wrappedLog != nil {
			logs.GetLogger().Log(wrappedLog)
		}*/

	if !config.GlobalConf.InDevelopment {
		errorMsg = ""
	}
	// Return status code with error message
	return c.Status(userError.StatusCode()).JSON(models.ErrorResponse{
		Timestamp:   userError.Time(),
		Errno:       userError.Name(),
		Message:     userError.Error(),
		Error:       true,
		FullMessage: errorMsg,
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
