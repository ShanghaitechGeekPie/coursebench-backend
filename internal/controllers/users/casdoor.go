package users

import (
	"bytes"
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	casdoorStateKey      = "casdoor_oauth_state"
	casdoorBindUserIDKey = "casdoor_bind_user_id"
	casdoorReturnURLKey  = "casdoor_return_url"
)

type casdoorTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func CasdoorLogin(c *fiber.Ctx) (err error) {
	return startCasdoorFlow(c, false)
}

func CasdoorBind(c *fiber.Ctx) (err error) {
	return startCasdoorFlow(c, true)
}

func CasdoorCallback(c *fiber.Ctx) (err error) {
	sess, err := session.GetStore().Get(c)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}

	state := c.Query("state")
	if state == "" || state != readStringSession(sess.Get(casdoorStateKey)) {
		return errors.New(errors.InvalidArgument)
	}

	code := c.Query("code")
	if code == "" {
		return errors.New(errors.InvalidArgument)
	}

	redirectURI := getCasdoorRedirectURI()
	token, err := getCasdoorOAuthToken(code, redirectURI)
	if err != nil {
		return err
	}
	if token.AccessToken == "" {
		return errors.New(errors.InternalServerError)
	}

	claims, err := parseJWTClaims(token.AccessToken)
	if err != nil {
		return err
	}

	casdoorSub := readClaimAsString(claims, "sub")
	if casdoorSub == "" {
		return errors.New(errors.InternalServerError)
	}
	email := strings.ToLower(strings.TrimSpace(readClaimAsString(claims, "email")))
	nickname := strings.TrimSpace(readClaimAsString(claims, "name"))
	if nickname == "" {
		nickname = strings.TrimSpace(readClaimAsString(claims, "preferred_username"))
	}
	realname := nickname

	bindUserID := readUintSession(sess.Get(casdoorBindUserIDKey))
	var user *models.User
	if bindUserID != 0 {
		err = queries.BindCasdoorIdentity(nil, bindUserID, casdoorSub)
		if err != nil {
			return err
		}
		user, err = queries.GetUserByID(nil, bindUserID)
		if err != nil {
			return err
		}
	} else {
		user, err = queries.GetUserByCasdoorSub(nil, casdoorSub)
		if err != nil {
			if !errors.Is(err, errors.UserNotExists) {
				return err
			}

			if email != "" {
				user, err = queries.GetUserByEmail(nil, email)
				if err != nil && !errors.Is(err, errors.UserNotExists) {
					return err
				}
				if err == nil {
					err = queries.BindCasdoorIdentity(nil, user.ID, casdoorSub)
					if err != nil {
						return err
					}
				}
			}

			if user == nil {
				user, err = queries.CreateOAuthUser(nil, email, nickname, realname, casdoorSub)
				if err != nil {
					return err
				}
			}
		}
	}

	sess.Set("user_id", user.ID)
	sess.Delete(casdoorStateKey)
	sess.Delete(casdoorBindUserIDKey)
	returnURL := readStringSession(sess.Get(casdoorReturnURLKey))
	sess.Delete(casdoorReturnURLKey)
	if err = sess.Save(); err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}

	redirectURL := buildFrontendOAuthCallbackURL(returnURL)
	return c.Redirect(redirectURL, http.StatusFound)
}

func CasdoorUnbind(c *fiber.Ctx) (err error) {
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	if err = queries.UnbindCasdoorIdentity(nil, id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{},
		Error: false,
	})
}

func startCasdoorFlow(c *fiber.Ctx, bind bool) error {
	if config.GlobalConf.CasdoorEndpoint == "" || config.GlobalConf.CasdoorClientID == "" || config.GlobalConf.CasdoorClientSecret == "" {
		return errors.New(errors.InvalidArgument)
	}

	sess, err := session.GetStore().Get(c)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}

	if bind {
		id, e := session.GetUserID(c)
		if e != nil {
			return e
		}
		sess.Set(casdoorBindUserIDKey, id)
	} else {
		sess.Delete(casdoorBindUserIDKey)
	}

	state, err := generateOAuthState(16)
	if err != nil {
		return err
	}
	sess.Set(casdoorStateKey, state)

	returnURL := sanitizeReturnURL(c.Query("return_url"))
	if returnURL == "" {
		returnURL = "/"
	}
	sess.Set(casdoorReturnURLKey, returnURL)
	if err = sess.Save(); err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}

	authURL := buildCasdoorSigninURL(state, getCasdoorRedirectURI())
	return c.Redirect(authURL, http.StatusFound)
}

func buildCasdoorSigninURL(state string, redirectURI string) string {
	q := url.Values{}
	q.Set("client_id", config.GlobalConf.CasdoorClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", "read")
	q.Set("state", state)
	return strings.TrimRight(config.GlobalConf.CasdoorEndpoint, "/") + "/login/oauth/authorize?" + q.Encode()
}

func getCasdoorRedirectURI() string {
	if config.GlobalConf.CasdoorRedirectURI != "" {
		return config.GlobalConf.CasdoorRedirectURI
	}
	return strings.TrimRight(config.GlobalConf.ServerURL, "/") + "/v1/user/casdoor/callback"
}

func getCasdoorOAuthToken(code string, redirectURI string) (*casdoorTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", config.GlobalConf.CasdoorClientID)
	form.Set("client_secret", config.GlobalConf.CasdoorClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	endpoint := strings.TrimRight(config.GlobalConf.CasdoorEndpoint, "/") + "/api/login/oauth/access_token"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, errors.InternalServerError)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errors.InternalServerError)
	}
	defer func() { _ = resp.Body.Close() }()

	token := &casdoorTokenResponse{}
	if err = json.NewDecoder(resp.Body).Decode(token); err != nil {
		return nil, errors.Wrap(err, errors.InternalServerError)
	}
	return token, nil
}

func parseJWTClaims(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.New(errors.InvalidArgument)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.Wrap(err, errors.InvalidArgument)
	}
	claims := map[string]interface{}{}
	if err = json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.Wrap(err, errors.InvalidArgument)
	}
	return claims, nil
}

func readClaimAsString(claims map[string]interface{}, key string) string {
	if claims == nil {
		return ""
	}
	v, ok := claims[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func sanitizeReturnURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if !u.IsAbs() {
		if strings.HasPrefix(raw, "/") {
			return raw
		}
		return ""
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	return u.String()
}

func buildFrontendOAuthCallbackURL(returnURL string) string {
	q := url.Values{}
	q.Set("return_url", returnURL)
	callbackPath := "/oauth/casdoor"
	if base := strings.TrimSpace(config.GlobalConf.CasdoorFrontendURL); base != "" {
		return strings.TrimRight(base, "/") + callbackPath + "?" + q.Encode()
	}

	if parsed, err := url.Parse(returnURL); err == nil && parsed.IsAbs() {
		origin := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
		return origin + callbackPath + "?" + q.Encode()
	}

	return callbackPath + "?" + q.Encode()
}

func readStringSession(value interface{}) string {
	if value == nil {
		return ""
	}
	s, ok := value.(string)
	if ok {
		return s
	}
	return ""
}

func readUintSession(value interface{}) uint {
	switch v := value.(type) {
	case uint:
		return v
	case uint64:
		return uint(v)
	case int:
		if v > 0 {
			return uint(v)
		}
	}
	return 0
}

func generateOAuthState(length int) (string, error) {
	if length <= 0 {
		return "", errors.New(errors.InvalidArgument)
	}
	b := make([]byte, length)
	if _, err := crand.Read(b); err != nil {
		return "", errors.Wrap(err, errors.InternalServerError)
	}
	return hex.EncodeToString(b), nil
}
