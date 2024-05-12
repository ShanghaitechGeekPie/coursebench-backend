package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"

	"gorm.io/gorm"
)

func Ranklist(db *gorm.DB) ([]models.RanklistResponse, error) {
	if db == nil {
		db = database.GetDB()
	}

	var ranklist []models.RanklistResponse
	result := db.Model(&models.User{}).
		Select("nick_name", "reward", "is_anonymous").
		Order("reward DESC").
		Limit(30).
		Scan(&ranklist)

	for i := range ranklist {
		if ranklist[i].IsAnonymous {
			ranklist[i].NickName = ""
		}
	}

	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	return ranklist, nil
}
