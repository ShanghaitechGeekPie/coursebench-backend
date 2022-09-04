package queries

import (
	"bytes"
	"coursebench-backend/3rd/captcha"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"image/png"
	"time"
)

const CaptchaLength = 6
const CaptchaExpire = 60 * 10

func GenerateCaptcha(ctx *fiber.Ctx) (string, error) {
	sess, err := session.GetSession(ctx)
	if err != nil {
		return "", errors.Wrap(err, errors.InternalServerError)
	}
	id := uuid.New().String()
	digit := captcha.RandomDigits(CaptchaLength)
	img := captcha.NewImage(id, digit, captcha.StdWidth, captcha.StdHeight)
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, img)
	if err != nil {
		return "", err
	}
	digit_s := ""
	for _, v := range digit {
		digit_s += string(v + '0')
	}
	sess.Set("captcha_code", digit_s)
	sess.Set("captcha_time", time.Now().Unix())
	if err := sess.Save(); err != nil {
		return "", errors.Wrap(err, errors.InternalServerError)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func VerifyCaptcha(ctx *fiber.Ctx, digits string) error {
	sess, err := session.GetSession(ctx)
	if err != nil {
		return errors.New(errors.NoCaptchaToken)
	}
	id := sess.Get("captcha_code")
	if id == nil {
		return errors.New(errors.NoCaptchaToken)
	}
	id_s := ""
	var gen_time int64
	ok := false
	if id_s, ok = id.(string); !ok {
		return errors.New(errors.InternalServerError)
	}
	if gen_time, ok = sess.Get("captcha_time").(int64); !ok {
		return errors.New(errors.InternalServerError)
	}
	if gen_time+CaptchaExpire < time.Now().Unix() {
		return errors.New(errors.CaptchaExpired)
	}
	if id_s != digits {
		return errors.New(errors.CaptchaMismatch)
	}
	sess.Set("captcha_time", int64(-10000))
	if err := sess.Save(); err != nil {
		return errors.New(errors.InternalServerError)
	}
	return nil
}
