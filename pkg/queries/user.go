package queries

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/mail"
	"coursebench-backend/pkg/models"
	"fmt"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
	"unicode"
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

func Register(db *gorm.DB, u *models.User) error {
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

// id: 被查询用户的id
// uid: 查询用户的id
// ip: 查询用户的ip
func GetProfile(db *gorm.DB, id uint, uid uint) (models.ProfileResponse, error) {
	if db == nil {
		db = database.GetDB()
	}
	user, err := GetUserByID(db, id)
	if err != nil {
		return models.ProfileResponse{}, err
	}
	avatar := ""
	if user.Avatar != "" {
		avatar = fmt.Sprintf("https://%s/%s/avatar/%s", database.GetEndpoint(), database.MinioConf.Bucket, user.Avatar)
	}
	if user.IsAnonymous && id != uid {
		return models.ProfileResponse{ID: id, NickName: user.NickName, Avatar: avatar, IsAnonymous: user.IsAnonymous, IsAdmin: user.IsAdmin}, nil
	} else {
		return models.ProfileResponse{ID: id, Email: user.Email, Year: user.Year, Grade: user.Grade, NickName: user.NickName, RealName: user.RealName, IsAnonymous: user.IsAnonymous, Avatar: avatar, IsAdmin: user.IsAdmin}, nil
	}
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
