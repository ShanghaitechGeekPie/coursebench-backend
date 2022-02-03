package courses

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type InfoResponse struct {
	Name       string              `json:"name"`
	Code       string              `json:"code"`
	ID         int                 `json:"id"`
	Institute  string              `json:"institute"`
	Credit     int                 `json:"credit"`
	Score      []float64           `json:"score"`
	CommentNum int                 `json:"comment_num"`
	Groups     []InfoResponseGroup `json:"groups"`
}

type InfoResponseGroup struct {
	ID         int       `json:"id"`
	Code       string    `json:"code"`
	Score      []float64 `json:"score"`
	CommentNum int       `json:"comment_num"`
	Teachers   []struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"teachers"`
}

func Info(c *fiber.Ctx) (err error) {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id := uint(id_i)
	db := database.GetDB()
	course := &models.Course{}
	err = db.Preload("Groups").Preload("Groups.Teachers").Where("id = ?", id).First(course).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.CourseNotExists)
	}
	response := &InfoResponse{Name: course.Name, Code: course.Code, ID: int(course.ID), Institute: course.Institute, Credit: course.Credit, CommentNum: course.CommentCount}
	response.Score = make([]float64, course.CommentCount)
	if course.CommentCount != 0 {
		for i := 0; i < course.CommentCount; i++ {
			response.Score[i] = float64(course.Scores[i]) / float64(course.CommentCount)
		}
	}
	response.Groups = make([]InfoResponseGroup, len(course.Groups))
	for i, g := range course.Groups {
		response.Groups[i].ID = int(g.ID)
		response.Groups[i].Code = g.Code
		response.Groups[i].CommentNum = g.CommentCount
		response.Groups[i].Score = make([]float64, g.CommentCount)
		if g.CommentCount != 0 {
			for j := 0; j < g.CommentCount; j++ {
				response.Groups[i].Score[j] = float64(g.Scores[j]) / float64(g.CommentCount)
			}
		}
		response.Groups[i].Teachers = make([]struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		}, len(g.Teachers))
		for j, t := range g.Teachers {
			response.Groups[i].Teachers[j].Name = t.Name
			response.Groups[i].Teachers[j].ID = int(t.ID)
		}
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
