package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 成就触发器配置结构
type CounterTriggerConfig struct {
	CounterTarget int    `json:"counter_target"` // 目标数量
	Action        string `json:"action"`         // 动作类型: "comment", "like_received", etc.
}

type TimeBasedTriggerConfig struct {
	TimeWindow    string `json:"time_window"`    // 时间窗口: "24h", "7d", etc.
	RequiredCount int    `json:"required_count"` // 要求的数量
	Action        string `json:"action"`         // 动作类型
}

type ComplexTriggerConfig struct {
	Conditions []map[string]interface{} `json:"conditions"` // 多个条件
	Logic      string                   `json:"logic"`      // "AND" 或 "OR"
}

// CheckAndGrantAchievements 检查并授予用户成就（主入口函数）
// 这个函数应该在用户执行某个动作后被调用
func CheckAndGrantAchievements(db *gorm.DB, userID uint, action string, value int) error {
	if db == nil {
		db = database.GetDB()
	}

	// 获取所有相关的成就定义
	var achievements []models.Achievement
	err := db.Where("is_active = ? AND type IN (?)", true,
		[]models.AchievementType{
			models.AchievementTypeCounter,
			models.AchievementTypeTimeBased,
		}).Find(&achievements).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	// 遍历检查每个成就
	for _, achievement := range achievements {
		// 检查用户是否已经获得该成就
		hasAchievement, err := HasUserAchievement(db, userID, achievement.ID)
		if err != nil {
			continue // 记录错误但继续处理其他成就
		}
		if hasAchievement {
			continue
		}

		// 根据成就类型进行检查
		switch achievement.Type {
		case models.AchievementTypeCounter:
			err = checkCounterAchievement(db, userID, achievement, action, value)
		case models.AchievementTypeTimeBased:
			err = checkTimeBasedAchievement(db, userID, achievement, action, value)
		}

		if err != nil {
			// 记录错误但继续处理
			continue
		}
	}

	return nil
}

// checkCounterAchievement 检查计数类型成就
func checkCounterAchievement(db *gorm.DB, userID uint, achievement models.Achievement, action string, value int) error {
	var config CounterTriggerConfig
	if err := json.Unmarshal([]byte(achievement.TriggerConfig), &config); err != nil {
		return err
	}

	// 检查动作是否匹配
	if config.Action != action {
		return nil
	}

	// 获取或创建进度记录
	var progress models.AchievementProgress
	err := db.Where("user_id = ? AND achievement_code = ?", userID, achievement.Code).
		FirstOrCreate(&progress, models.AchievementProgress{
			UserID:          userID,
			AchievementCode: achievement.Code,
			CurrentValue:    0,
			TargetValue:     config.CounterTarget,
		}).Error
	if err != nil {
		return err
	}

	// 更新进度
	progress.CurrentValue += value
	progress.LastUpdated = time.Now()

	if err := db.Save(&progress).Error; err != nil {
		return err
	}

	// 检查是否达成
	if progress.CurrentValue >= config.CounterTarget {
		return grantAchievement(db, userID, achievement)
	}

	return nil
}

// checkTimeBasedAchievement 检查时间相关成就
func checkTimeBasedAchievement(db *gorm.DB, userID uint, achievement models.Achievement, action string, value int) error {
	var config TimeBasedTriggerConfig
	if err := json.Unmarshal([]byte(achievement.TriggerConfig), &config); err != nil {
		return err
	}

	if config.Action != action {
		return nil
	}

	// 解析时间窗口
	duration, err := time.ParseDuration(config.TimeWindow)
	if err != nil {
		return err
	}

	now := time.Now()
	windowStart := now.Add(-duration)
	windowEnd := now

	// 获取或创建进度记录
	var progress models.AchievementProgress
	err = db.Where("user_id = ? AND achievement_code = ?", userID, achievement.Code).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新的进度记录
		progress = models.AchievementProgress{
			UserID:          userID,
			AchievementCode: achievement.Code,
			CurrentValue:    value,
			TargetValue:     config.RequiredCount,
			WindowStart:     &windowStart,
			WindowEnd:       &windowEnd,
			LastUpdated:     now,
		}
		if err := db.Create(&progress).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// 检查是否需要重置时间窗口
		if progress.WindowEnd != nil && now.After(*progress.WindowEnd) {
			// 时间窗口已过期，重置
			progress.CurrentValue = value
			progress.WindowStart = &windowStart
			progress.WindowEnd = &windowEnd
		} else {
			// 在时间窗口内，累加
			progress.CurrentValue += value
		}
		progress.LastUpdated = now

		if err := db.Save(&progress).Error; err != nil {
			return err
		}
	}

	// 检查是否达成
	if progress.CurrentValue >= config.RequiredCount {
		return grantAchievement(db, userID, achievement)
	}

	return nil
}

// grantAchievement 授予成就（事务安全，幂等性保证）
func grantAchievement(db *gorm.DB, userID uint, achievement models.Achievement) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 使用 ON CONFLICT DO NOTHING 实现幂等性
		userAchievement := models.UserAchievement{
			UserID:        userID,
			AchievementID: achievement.ID,
			UnlockedAt:    time.Now(),
			Progress:      100, // 已完成
			IsNotified:    false,
		}

		// 使用 Clauses 添加 ON CONFLICT 子句
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "achievement_id"}},
			DoNothing: true, // 如果已存在则不做任何操作
		}).Create(&userAchievement).Error

		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}

		// 检查是否实际插入了新记录（RowsAffected > 0）
		// 如果 RowsAffected == 0，说明记录已存在，不需要更新统计
		if tx.RowsAffected == 0 {
			return nil // 成就已存在，直接返回
		}

		// 更新用户成就统计
		err = updateUserAchievementStat(tx, userID, achievement)
		if err != nil {
			return err
		}

		// 如果成就有积分奖励，更新用户积分
		if achievement.RewardPoints > 0 {
			err = tx.Model(&models.User{}).Where("id = ?", userID).
				Update("reward", gorm.Expr("reward + ?", achievement.RewardPoints)).Error
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
		}

		return nil
	})
}

// updateUserAchievementStat 更新用户成就统计
func updateUserAchievementStat(db *gorm.DB, userID uint, achievement models.Achievement) error {
	var stat models.UserAchievementStat
	err := db.Where("user_id = ?", userID).FirstOrCreate(&stat, models.UserAchievementStat{
		UserID: userID,
	}).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	stat.TotalAchievements++
	stat.TotalPoints += achievement.RewardPoints

	// 根据稀有度更新统计
	switch achievement.Rarity {
	case models.RarityCommon:
		stat.CommonCount++
	case models.RarityUncommon:
		stat.UncommonCount++
	case models.RarityRare:
		stat.RareCount++
	case models.RarityEpic:
		stat.EpicCount++
	case models.RarityLegendary:
		stat.LegendaryCount++
	}

	return db.Save(&stat).Error
}

// HasUserAchievement 检查用户是否已拥有某个成就
func HasUserAchievement(db *gorm.DB, userID uint, achievementID uint) (bool, error) {
	var count int64
	err := db.Model(&models.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error
	if err != nil {
		return false, errors.Wrap(err, errors.DatabaseError)
	}
	return count > 0, nil
}

// GetUserAchievements 获取用户的所有成就（包括进度）
func GetUserAchievements(db *gorm.DB, userID uint) (*models.UserAchievementListResponse, error) {
	if db == nil {
		db = database.GetDB()
	}

	// 获取所有成就定义
	var allAchievements []models.Achievement
	err := db.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&allAchievements).Error
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	// 获取用户已解锁的成就
	var userAchievements []models.UserAchievement
	err = db.Where("user_id = ?", userID).Preload("Achievement").Find(&userAchievements).Error
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	// 创建已解锁成就的映射
	unlockedMap := make(map[uint]models.UserAchievement)
	for _, ua := range userAchievements {
		unlockedMap[ua.AchievementID] = ua
	}

	// 获取进度信息
	var progressList []models.AchievementProgress
	err = db.Where("user_id = ?", userID).Find(&progressList).Error
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	// 创建进度映射
	progressMap := make(map[string]models.AchievementProgress)
	for _, p := range progressList {
		progressMap[p.AchievementCode] = p
	}

	// 分类成就
	response := &models.UserAchievementListResponse{
		Unlocked:   []models.AchievementResponse{},
		InProgress: []models.AchievementResponse{},
		Locked:     []models.AchievementResponse{},
	}

	for _, achievement := range allAchievements {
		achievementResp := models.AchievementResponse{
			ID:           achievement.ID,
			Code:         achievement.Code,
			Name:         achievement.Name,
			Description:  achievement.Description,
			Icon:         achievement.Icon,
			Type:         achievement.Type,
			Category:     achievement.Category,
			Rarity:       achievement.Rarity,
			RewardPoints: achievement.RewardPoints,
			IsHidden:     achievement.IsHidden,
			IsUnlocked:   false,
		}

		// 检查是否已解锁
		if ua, unlocked := unlockedMap[achievement.ID]; unlocked {
			achievementResp.IsUnlocked = true
			achievementResp.UnlockedAt = &ua.UnlockedAt
			progress := 100
			achievementResp.Progress = &progress
			response.Unlocked = append(response.Unlocked, achievementResp)
		} else if progress, hasProgress := progressMap[achievement.Code]; hasProgress {
			// 有进度但未解锁
			achievementResp.Progress = &progress.CurrentValue
			achievementResp.TargetValue = &progress.TargetValue
			response.InProgress = append(response.InProgress, achievementResp)
		} else if !achievement.IsHidden {
			// 未开始且非隐藏成就
			response.Locked = append(response.Locked, achievementResp)
		}
	}

	return response, nil
}

// GetUserAchievementStats 获取用户成就统计
func GetUserAchievementStats(db *gorm.DB, userID uint) (*models.AchievementStatResponse, error) {
	if db == nil {
		db = database.GetDB()
	}

	// 获取统计数据
	var stat models.UserAchievementStat
	err := db.Where("user_id = ?", userID).First(&stat).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	// 获取最近解锁的成就
	var recentUnlocked []models.UserAchievement
	err = db.Where("user_id = ?", userID).
		Preload("Achievement").
		Order("unlocked_at DESC").
		Limit(5).
		Find(&recentUnlocked).Error
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	recentList := []models.AchievementResponse{}
	for _, ua := range recentUnlocked {
		progress := 100
		recentList = append(recentList, models.AchievementResponse{
			ID:           ua.Achievement.ID,
			Code:         ua.Achievement.Code,
			Name:         ua.Achievement.Name,
			Description:  ua.Achievement.Description,
			Icon:         ua.Achievement.Icon,
			Type:         ua.Achievement.Type,
			Category:     ua.Achievement.Category,
			Rarity:       ua.Achievement.Rarity,
			RewardPoints: ua.Achievement.RewardPoints,
			IsUnlocked:   true,
			UnlockedAt:   &ua.UnlockedAt,
			Progress:     &progress,
		})
	}

	response := &models.AchievementStatResponse{
		TotalAchievements: stat.TotalAchievements,
		UnlockedCount:     stat.TotalAchievements,
		TotalPoints:       stat.TotalPoints,
		RarityBreakdown: map[models.AchievementRarity]int{
			models.RarityCommon:    stat.CommonCount,
			models.RarityUncommon:  stat.UncommonCount,
			models.RarityRare:      stat.RareCount,
			models.RarityEpic:      stat.EpicCount,
			models.RarityLegendary: stat.LegendaryCount,
		},
		RecentUnlocked: recentList,
	}

	return response, nil
}

// GetActualUserStatistics 获取用户的实际统计数据（用于手动验证）
// 这些统计数据可以用来补发成就或验证系统准确性
func GetActualUserStatistics(db *gorm.DB, userID uint) (map[string]int, error) {
	if db == nil {
		db = database.GetDB()
	}

	stats := make(map[string]int)

	// 评论数
	var commentCount int64
	db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentCount)
	stats["total_comments"] = int(commentCount)

	// 收到的点赞数
	var totalLikes int64
	db.Model(&models.Comment{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(like), 0)").Scan(&totalLikes)
	stats["total_likes_received"] = int(totalLikes)

	// 给出的点赞数
	var likesGiven int64
	db.Model(&models.CommentLike{}).Where("user_id = ? AND is_like = ?", userID, true).Count(&likesGiven)
	stats["total_likes_given"] = int(likesGiven)

	// 回复数
	var replyCount int64
	db.Model(&models.Reply{}).Where("user_id = ?", userID).Count(&replyCount)
	stats["total_replies"] = int(replyCount)

	return stats, nil
}

// InitializeDefaultAchievements 初始化默认成就（用于系统安装时）
func InitializeDefaultAchievements(db *gorm.DB) error {
	if db == nil {
		db = database.GetDB()
	}

	achievements := []models.Achievement{
		// 评论相关成就
		{
			Code:          "first_comment",
			Name:          "初次发声",
			Description:   "发布你的第一条评论",
			Type:          models.AchievementTypeCounter,
			Category:      models.CategoryComment,
			Rarity:        models.RarityCommon,
			TriggerConfig: `{"counter_target": 1, "action": "comment"}`,
			RewardPoints:  10,
			SortOrder:     1,
		},
		{
			Code:          "comment_veteran_10",
			Name:          "评论新手",
			Description:   "发布10条评论",
			Type:          models.AchievementTypeCounter,
			Category:      models.CategoryComment,
			Rarity:        models.RarityCommon,
			TriggerConfig: `{"counter_target": 10, "action": "comment"}`,
			RewardPoints:  50,
			SortOrder:     2,
		},
		{
			Code:          "comment_veteran_50",
			Name:          "评论达人",
			Description:   "发布50条评论",
			Type:          models.AchievementTypeCounter,
			Category:      models.CategoryComment,
			Rarity:        models.RarityUncommon,
			TriggerConfig: `{"counter_target": 50, "action": "comment"}`,
			RewardPoints:  200,
			SortOrder:     3,
		},
		{
			Code:          "comment_master_100",
			Name:          "评论大师",
			Description:   "发布100条评论",
			Type:          models.AchievementTypeCounter,
			Category:      models.CategoryComment,
			Rarity:        models.RarityRare,
			TriggerConfig: `{"counter_target": 100, "action": "comment"}`,
			RewardPoints:  500,
			SortOrder:     4,
		},

		// 点赞相关成就
		{
			Code:          "popular_10",
			Name:          "小有名气",
			Description:   "累计收到10个点赞",
			Type:          models.AchievementTypeCounter,
			Category:      models.CategoryLike,
			Rarity:        models.RarityCommon,
			TriggerConfig: `{"counter_target": 10, "action": "like_received"}`,
			RewardPoints:  30,
			SortOrder:     10,
		},
		{
			Code:          "popular_50",
			Name:          "人气之星",
			Description:   "累计收到50个点赞",
			Type:          models.AchievementTypeCounter,
			Category:      models.CategoryLike,
			Rarity:        models.RarityUncommon,
			TriggerConfig: `{"counter_target": 50, "action": "like_received"}`,
			RewardPoints:  100,
			SortOrder:     11,
		},
		{
			Code:          "popular_100",
			Name:          "万众瞩目",
			Description:   "累计收到100个点赞",
			Type:          models.AchievementTypeCounter,
			Category:      models.CategoryLike,
			Rarity:        models.RarityRare,
			TriggerConfig: `{"counter_target": 100, "action": "like_received"}`,
			RewardPoints:  300,
			SortOrder:     12,
		},

		// 时间相关成就
		{
			Code:          "active_day",
			Name:          "活跃的一天",
			Description:   "24小时内发布5条评论",
			Type:          models.AchievementTypeTimeBased,
			Category:      models.CategoryTime,
			Rarity:        models.RarityUncommon,
			TriggerConfig: `{"time_window": "24h", "required_count": 5, "action": "comment"}`,
			RewardPoints:  100,
			SortOrder:     20,
		},
		{
			Code:          "active_week",
			Name:          "活跃的一周",
			Description:   "7天内发布20条评论",
			Type:          models.AchievementTypeTimeBased,
			Category:      models.CategoryTime,
			Rarity:        models.RarityRare,
			TriggerConfig: `{"time_window": "168h", "required_count": 20, "action": "comment"}`,
			RewardPoints:  300,
			SortOrder:     21,
		},
	}

	for _, achievement := range achievements {
		// 使用 FirstOrCreate 避免重复插入
		err := db.Where("code = ?", achievement.Code).FirstOrCreate(&achievement).Error
		if err != nil {
			return fmt.Errorf("failed to create achievement %s: %w", achievement.Code, err)
		}
	}

	return nil
}
