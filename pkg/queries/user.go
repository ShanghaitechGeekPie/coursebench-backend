package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/mail"
	"coursebench-backend/pkg/models"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
	"unicode"
)

func Register(u *models.User) error {
	db := database.GetDB()

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

	// 检查邮箱是否已存在
	result := db.Where("email = ?", u.Email).Take(&models.User{})
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected > 0 {
		return errors.New(errors.UserEmailDuplicated)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	u.Password = string(hash)
	u.IsActive = false

	if err := db.Create(u).Error; err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	err = mail.PostMail(u)
	if err != nil {
		return err
	}
	return nil
}

func RegisterActive(id uint, code string) (err error) {
	db := database.GetDB()
	user := &models.User{}
	// 检查邮箱是否已存在
	result := db.Where("id = ?", id).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return errors.New(errors.UserNotExists)
	}
	ok, err := mail.CheckCode(user, code)
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

func Login(email, password string) (*models.User, error) {
	db := database.GetDB()

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

func UpdateProfile(id uint, year int, grade models.GradeType, nickname string, realname string, isAnonymous bool) (err error) {
	db := database.GetDB()
	user, err := GetUserByID(id)
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

func UpdatePassword(id uint, oldPassword string, newPassword string) (err error) {
	db := database.GetDB()
	user, err := GetUserByID(id)
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

func GetUserByID(id uint) (*models.User, error) {
	db := database.GetDB()

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
func GetProfile(id uint, uid uint) (models.ProfileResponse, error) {
	user, err := GetUserByID(id)
	if err != nil {
		return models.ProfileResponse{}, err
	}
	if user.IsAnonymous && id != uid {
		return models.ProfileResponse{ID: id, NickName: user.NickName, Avatar: user.Avatar, IsAnonymous: user.IsAnonymous}, nil
	} else {
		return models.ProfileResponse{ID: id, Email: user.Email, Year: user.Year, Grade: user.Grade, NickName: user.NickName, RealName: user.RealName, IsAnonymous: user.IsAnonymous, Avatar: user.Avatar}, nil
	}
}

func CheckYear(year int) bool {
	return year == 0 || (year >= 2014 && year <= time.Now().Year())
}

func CheckGrade(grade models.GradeType) bool {
	return grade == models.Undergraduate || grade == models.Postgraduate || grade == models.PhDStudent || grade == models.UnknownGrade
}

func CheckEmail(email string) bool {
	if len(email) > 30 {
		return false
	}
	if strings.Contains(email, "+") {
		return false
	}
	if err := checkmail.ValidateFormat(email); err != nil {
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
	if len(nickname) > 20 {
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
	if len(realname) > 20 {
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
