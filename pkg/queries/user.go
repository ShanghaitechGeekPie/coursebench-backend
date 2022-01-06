package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Register(u *models.User) error {
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
	result := db.Where("email = ?", u.Email).Take(&models.User{})
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

func Login(email, password string) (*models.User, error) {
	db := database.GetDB()

	// 检查输入合法
	if !CheckEmail(email) {
		return nil, errors.InvalidArgument
	}
	if !CheckPassword(password) {
		return nil, errors.InvalidArgument
	}

	user := &models.User{}
	// 检查邮箱是否已存在
	result := db.Where("email = ?", email).Take(user)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, errors.UserDoNotExist
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.Wrap(err, errors.UserPasswordIncorrect)
	}

	return user, nil
}

func CheckYear(year int) bool {
	return year >= 2014 && year <= time.Now().Year()
}

func CheckGrade(grade models.GradeType) bool {
	return grade == models.Undergraduate || grade == models.Postgraduate || grade == models.PhDStudent
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
