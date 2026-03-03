package achievement

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// List 获取用户的所有成就列表
func List(c *fiber.Ctx) error {
	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	db := database.GetDB()
	achievements, err := queries.GetUserAchievements(db, uid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  achievements,
		Error: false,
	})
}

// Stats 获取用户的成就统计
func Stats(c *fiber.Ctx) error {
	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	db := database.GetDB()
	stats, err := queries.GetUserAchievementStats(db, uid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  stats,
		Error: false,
	})
}

// UserAchievements 获取指定用户的公开成就信息
func UserAchievements(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return errors.New(errors.InvalidArgument)
	}

	db := database.GetDB()

	// 检查用户是否存在
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return errors.Wrap(err, errors.UserNotExists)
	}

	// 如果用户设置了匿名，只返回统计数据
	achievements, err := queries.GetUserAchievements(db, uint(userID))
	if err != nil {
		return err
	}

	// 如果用户是匿名的，只返回已解锁的成就（不显示进度）
	if user.IsAnonymous {
		achievements.InProgress = []models.AchievementResponse{}
		achievements.Locked = []models.AchievementResponse{}
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  achievements,
		Error: false,
	})
}

// AdminInitialize 管理员初始化默认成就（仅管理员可访问）
func AdminInitialize(c *fiber.Ctx) error {
	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	db := database.GetDB()
	var user models.User
	if err := db.First(&user, uid).Error; err != nil {
		return errors.Wrap(err, errors.UserNotExists)
	}

	if !user.IsAdmin {
		return errors.New(errors.PermissionDenied)
	}

	err = queries.InitializeDefaultAchievements(db)
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]string{"message": "Default achievements initialized"},
		Error: false,
	})
}

// AdminRecalculate 管理员重新计算用户成就（用于修复数据）
func AdminRecalculate(c *fiber.Ctx) error {
	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	db := database.GetDB()
	var user models.User
	if err := db.First(&user, uid).Error; err != nil {
		return errors.Wrap(err, errors.UserNotExists)
	}

	if !user.IsAdmin {
		return errors.New(errors.PermissionDenied)
	}

	// 获取要重新计算的用户ID
	targetUserIDStr := c.Query("user_id")
	if targetUserIDStr == "" {
		return errors.New(errors.InvalidArgument)
	}

	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return errors.New(errors.InvalidArgument)
	}

	// 获取用户的实际统计数据
	stats, err := queries.GetActualUserStatistics(db, uint(targetUserID))
	if err != nil {
		return err
	}

	// 根据统计数据重新检查成就
	// 评论成就
	if commentCount, ok := stats["total_comments"]; ok {
		err = queries.CheckAndGrantAchievements(db, uint(targetUserID), "comment", commentCount)
		if err != nil {
			return err
		}
	}

	// 点赞成就
	if likesReceived, ok := stats["total_likes_received"]; ok {
		err = queries.CheckAndGrantAchievements(db, uint(targetUserID), "like_received", likesReceived)
		if err != nil {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"recalculated": true, "stats": stats},
		Error: false,
	})
}
