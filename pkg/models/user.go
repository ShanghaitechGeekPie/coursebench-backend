package models

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/modelRegister"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type GradeType int

const (
	Undergraduate GradeType = 1
	Postgraduate  GradeType = 2
	PhDStudent    GradeType = 3
)

type User struct {
	gorm.Model
	Email    string
	Password string
	NickName string
	RealName string
	Year     int
	Grade    GradeType
}

func init() {
	modelRegister.Register(&User{})
}

func (u *User) Register() error {
	db := database.GetDB()

	// 检查输入合法
	if !CheckPassword(u.Password) {
		return errors.InvalidArgument
	}
	if !CheckYear(u.Year) {
		return errors.InvalidArgument
	}
	if !CheckGrade(u.Grade) {
		return errors.InvalidArgument
	}
	if !CheckEmail(u.Email) {
		return errors.InvalidArgument
	}

	// 检查邮箱是否已存在
	result := db.Where("email = ?", u.Email).Take(&User{})
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected > 0 {
		return errors.UserEmailDuplicated
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	u.Password = string(hash)

	if err := db.Create(u).Error; err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	return nil
}

func CheckYear(year int) bool {
	return year >= 2014 && year <= time.Now().Year()
}

func CheckGrade(grade GradeType) bool {
	return grade == Undergraduate || grade == Postgraduate || grade == PhDStudent
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
