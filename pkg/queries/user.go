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

package queries

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/mail"
	"coursebench-backend/pkg/models"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ResetPassword(db *gorm.DB, email string) error {
	if db == nil {
		db = database.GetDB()
	}
	if !CheckEmail(email) {
		return errors.New(errors.InvalidArgument)
	}
	user := models.User{}
	if err := db.Where("email = ?", email).Take(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.UserNotExists)
		}
		return errors.Wrap(err, errors.DatabaseError)
	}
	if !user.IsActive {
		return errors.New(errors.UserNotActive)
	}
	body := fmt.Sprintf(`<html><body><h1>您正在重置您在%s的密码</h1> <p>请点击该链接继续完成密码重置:</p><a href="{activeURL}">密码重置链接 </a> <br> <p>如果链接无法点击，请手动复制该链接并粘贴至浏览器：{activeURL} </p><br><br> <p>如果您没有注册过我们的服务或您没有进行过密码重置，请无视该邮件</p> </body></html>`, config.Text.ServiceName)
	return mail.PostMail(&user, "reset_password_mail_code", config.Text.ServiceName+"用户密码重置", "reset_password_active", body)
}

func ResetPasswordActive(db *gorm.DB, id uint, code string, password string) (err error) {
	if db == nil {
		db = database.GetDB()
	}
	if !CheckPassword(password) {
		return errors.New(errors.InvalidArgument)
	}
	user := &models.User{}
	result := db.Where("id = ?", id).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return errors.New(errors.UserNotExists)
	}
	if !user.IsActive {
		return errors.New(errors.UserNotActive)
	}
	ok, err := mail.CheckCode(user, code, "reset_password_mail_code")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(errors.MailCodeInvalid)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return errors.Wrap(err, errors.InternalServerError)
	}
	user.Password = string(hash)
	db.Select("password").Save(user)
	return nil
}

func Register(db *gorm.DB, u *models.User, invitation_code string) error {
	if db == nil {
		db = database.GetDB()
	}

	// 检查输入合法
	if !CheckPassword(u.Password) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckYear(u.Year) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckGrade(u.Grade) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckEmail(u.Email) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckNickName(u.NickName) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckRealName(u.RealName) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckInvitationCode(invitation_code) {
		return errors.New(errors.InvalidArgument)
	}

	// check if the invitation code is valid
	if invitation_code != "" {
		inviter, err := GetUserByInvitationCode(db, invitation_code)
		if err != nil {
			if errors.Is(err, errors.UserNotExists) {
				return errors.New(errors.InvitationCodeInvalid)
			}
			return err
		}

		u.InvitedByUserID = inviter.ID
		// TODO: only once for the inviter?
		inviter.Reward += 100
		db.Save(inviter)
	}

	// 检查邮箱是否已存在
	user := &models.User{}
	result := db.Where("email = ?", u.Email).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected > 0 {
		if user.IsActive {
			return errors.New(errors.UserEmailDuplicated)
		} else {
			db.Delete(user)
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	u.Password = string(hash)
	u.IsActive = false
	u.IsAdmin = false

	code, err := createInvitationCode(db)
	if err != nil {
		return err
	}
	u.InvitationCode = code

	if err = db.Create(u).Error; err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	body := fmt.Sprintf(`<html><body>
<h1>欢迎注册%s</h1> <p>我们已经接收到您的电子邮箱验证申请，请点击以下链接完成注册。</p>
<p>验证完成后，您将能够即刻发布课程评价，并与其他用户互动。</p>
<a href="{activeURL}">注册链接 </a> <br> 
<p>如无法点击，请手动复制该链接并粘贴至浏览器地址栏以完成注册：{activeURL} </p>
<p>预祝您在 CourseBench 玩得开心！</p> <br>
<p>如果您需要其它任何帮助，欢迎随时联系我们。</p>
<p>电邮地址：zhaoqch1@shanghaitech.edu.cn</p>
<p>如您并未注册CourseBench账号，请无视本邮件。</p><br>
<p>此致</p>
<p>%s 团队</p>
<br>
<h1>Thank you for registering for %s</h1> <p>We have received your application for verifying this email address. Please click on the link below to accomplish the process.</p>
<p>Once the registration is done, you will be able to post your comments on courses and interact with other users.</p>
<a href="{activeURL}">Register Link </a> <br> 
<p>If the link above isn’t working, in order to verify your account, please copy this URL and paste it on the address bar of your browser:  {activeURL} </p>
<p>Have fun at CourseBench!</p><br>
<p>If you need any help, please don’t hesitate to contact us.</p>
<p>Email: zhaoqch1@shanghaitech.edu.cn </p>
<p>If you are not registering for CourseBench, please ignore this email.</p><br>
<p>Yours,</p>
<p>%s Team</p>



 </body></html>`, config.Text.ServiceName, config.Text.ServiceName, config.Text.ServiceNameEN, config.Text.ServiceNameEN)
	err = mail.PostMail(u, "register_mail_code", config.Text.ServiceName+"用户注册验证", "active", body)
	if err != nil {
		return err
	}
	return nil
}

func RegisterActive(db *gorm.DB, id uint, code string) (err error) {
	if db == nil {
		db = database.GetDB()
	}
	user := &models.User{}
	// 检查邮箱是否已存在
	result := db.Where("id = ?", id).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return errors.New(errors.UserNotExists)
	}
	ok, err := mail.CheckCode(user, code, "register_mail_code")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(errors.MailCodeInvalid)
	}
	user.IsActive = true
	db.Select("is_active").Save(user)
	return nil
}

func Login(db *gorm.DB, email, password string) (*models.User, error) {
	if db == nil {
		db = database.GetDB()
	}

	// 检查输入合法
	if !CheckEmail(email) {
		return nil, errors.New(errors.InvalidArgument)
	}
	if !CheckPassword(password) {
		return nil, errors.New(errors.InvalidArgument)
	}

	user := &models.User{}
	// 检查邮箱是否已存在
	result := db.Where("email = ?", email).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(errors.UserNotExists)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.Wrap(err, errors.UserPasswordIncorrect)
	}

	if !user.IsActive {
		return nil, errors.New(errors.UserNotActive)
	}

	if user.InvitationCode == "" {
		code, err := createInvitationCode(db)
		if err != nil {
			return nil, err
		}

		user.InvitationCode = code
		err = db.Select("invitation_code").Save(user).Error
		if err != nil {
			return nil, errors.Wrap(err, errors.DatabaseError)
		}
	}

	return user, nil
}

func UpdateProfile(db *gorm.DB, id uint, year int, grade models.GradeType, nickname string, realname string, isAnonymous bool) (err error) {
	if db == nil {
		db = database.GetDB()
	}
	user, err := GetUserByID(db, id)
	if err != nil {
		return err
	}
	if !CheckYear(year) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckGrade(grade) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckNickName(nickname) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckRealName(realname) {
		return errors.New(errors.InvalidArgument)
	}
	user.Year = year
	user.Grade = grade
	user.NickName = nickname
	user.RealName = realname
	user.IsAnonymous = isAnonymous
	err = db.Select("year", "grade", "nick_name", "real_name", "is_anonymous").Save(user).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	return
}

func UpdatePassword(db *gorm.DB, id uint, oldPassword string, newPassword string) (err error) {
	if db == nil {
		db = database.GetDB()
	}
	user, err := GetUserByID(db, id)
	if err != nil {
		return err
	}
	if !CheckPassword(oldPassword) {
		return errors.New(errors.InvalidArgument)
	}
	if !CheckPassword(newPassword) {
		return errors.New(errors.InvalidArgument)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.Wrap(err, errors.UserPasswordIncorrect)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	user.Password = string(hash)
	err = db.Select("password").Save(user).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	return
}

func GetUserByID(db *gorm.DB, id uint) (*models.User, error) {
	if db == nil {
		db = database.GetDB()
	}

	user := &models.User{}
	result := db.Where("id = ?", id).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(errors.UserNotExists)
	}

	return user, nil
}

func GetProfile(db *gorm.DB, queriedUserID uint, queryingUserID uint) (models.ProfileResponse, error) {
	if db == nil {
		db = database.GetDB()
	}

	user, err := GetUserByID(db, queriedUserID)
	if err != nil {
		return models.ProfileResponse{}, err
	}

	displayInvitationCode := queryingUserID == queriedUserID
	displayReward := queryingUserID == queriedUserID

	// don't query if not logged in
	if queryingUserID != 0 && queryingUserID != queriedUserID {
		queryingUser, err := GetUserByID(db, queryingUserID)
		if err != nil {
			return models.ProfileResponse{}, err
		}
		if queryingUser.IsAdmin || queryingUser.IsCommunityAdmin {
			displayReward = true
		}
	}

	avatar := ""
	if user.Avatar != "" {
		avatar = fmt.Sprintf("https://%s/%s/avatar/%s", database.GetEndpoint(), database.MinioConf.Bucket, user.Avatar)
	}
	r := models.ProfileResponse{
		ID:               user.ID,
		NickName:         user.NickName,
		Avatar:           avatar,
		IsAnonymous:      user.IsAnonymous,
		IsAdmin:          user.IsAdmin,
		IsCommunityAdmin: user.IsCommunityAdmin,
	}

	if !user.IsAnonymous || queryingUserID == queriedUserID {
		r.Email = user.Email
		r.Year = user.Year
		r.Grade = user.Grade
		r.RealName = user.RealName
	}

	if displayInvitationCode {
		r.InvitationCode = user.InvitationCode
	}

	if displayReward {
		r.Reward = user.Reward
	} else {
		r.Reward = -1
	}

	return r, nil
}

func CheckYear(year int) bool {
	return year == 0 || (year >= 2014 && year <= time.Now().Year())
}

func CheckGrade(grade models.GradeType) bool {
	return grade == models.Undergraduate || grade == models.Postgraduate || grade == models.PhDStudent || grade == models.UnknownGrade
}

func CheckEmail(email string) bool {
	if len(email) > 100 {
		return false
	}
	if strings.Contains(email, "+") {
		return false
	}
	if err := checkmail.ValidateFormat(email); err != nil {
		return false
	}
	if !strings.HasSuffix(email, config.GlobalConf.MailSuffix) {
		return false
	}
	return true
}

func CheckPassword(password string) bool {
	if len(password) < 6 || len(password) > 30 {
		return false
	}
	for _, c := range password {
		if (c < '0' || c > '9') && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && !strings.Contains("!@#$%^&*()-_=+{}[]|\\:;'<>,.?/~`", string(c)) {
			return false
		}
	}
	return true
}

func CheckNickName(nickname string) bool {
	if len(nickname) > 40 {
		return false
	}
	r := []rune(nickname)
	for _, c := range r {
		if !unicode.IsGraphic(c) {
			return false
		}
	}
	return true
}

func CheckRealName(realname string) bool {
	if len(realname) > 30 {
		return false
	}
	r := []rune(realname)
	for _, c := range r {
		if !unicode.IsGraphic(c) {
			return false
		}
	}
	return true
}

func CheckInvitationCode(code string) bool {
	if len(code) == 0 {
		return true
	}
	if len(code) != 5 {
		return false
	}
	for _, c := range code {
		if (c < '0' || c > '9') && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return true
}

func createInvitationCode(db *gorm.DB) (string, error) {
	// try a few times before giving up
	for i := 0; i < 5; i++ {
		codeRunes := make([]rune, 0, 5)
		for i := 0; i < 5; i++ {
			codeRunes = append(codeRunes, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")[rand.Intn(62)])
		}
		code := string(codeRunes)

		_, err := GetUserByInvitationCode(db, code)
		if err != nil {
			if errors.Is(err, errors.UserNotExists) {
				return code, nil
			}
			return "", err
		}
	}

	return "", errors.New(errors.InternalServerError)
}

func GetUserByInvitationCode(db *gorm.DB, code string) (*models.User, error) {
	if db == nil {
		db = database.GetDB()
	}

	user := &models.User{}
	result := db.Where("invitation_code = ?", code).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(errors.UserNotExists)
	}

	return user, nil
}
