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

package models

import (
	"coursebench-backend/pkg/modelRegister"
	"time"

	"gorm.io/gorm"
)

// AchievementType 成就类型
type AchievementType string

const (
	// 计数类型成就
	AchievementTypeCounter AchievementType = "counter"
	// 时间相关成就
	AchievementTypeTimeBased AchievementType = "time_based"
	// 复合条件成就
	AchievementTypeComplex AchievementType = "complex"
	// 一次性成就
	AchievementTypeOneTime AchievementType = "one_time"
)

// AchievementCategory 成就分类
type AchievementCategory string

const (
	CategoryComment AchievementCategory = "comment" // 评论相关
	CategoryLike    AchievementCategory = "like"    // 点赞相关
	CategorySocial  AchievementCategory = "social"  // 社交相关
	CategoryExplore AchievementCategory = "explore" // 探索相关
	CategoryTime    AchievementCategory = "time"    // 时间相关
	CategorySpecial AchievementCategory = "special" // 特殊成就
)

// AchievementRarity 稀有度
type AchievementRarity string

const (
	RarityCommon    AchievementRarity = "common"    // 普通
	RarityUncommon  AchievementRarity = "uncommon"  // 少见
	RarityRare      AchievementRarity = "rare"      // 稀有
	RarityEpic      AchievementRarity = "epic"      // 史诗
	RarityLegendary AchievementRarity = "legendary" // 传说
)

// Achievement 成就定义表
type Achievement struct {
	gorm.Model
	Code        string              `gorm:"uniqueIndex;size:64"` // 成就唯一代码
	Name        string              `gorm:"size:128"`            // 成就名称
	Description string              `gorm:"type:text"`           // 成就描述
	Icon        string              `gorm:"size:255"`            // 图标URL
	Type        AchievementType     `gorm:"size:32;index"`       // 成就类型
	Category    AchievementCategory `gorm:"size:32;index"`       // 成就分类
	Rarity      AchievementRarity   `gorm:"size:32;index"`       // 稀有度

	// 触发条件配置(JSON格式)
	// 例如: {"counter_target": 10, "action": "comment"}
	// 或: {"time_window": "24h", "required_count": 5, "action": "comment"}
	TriggerConfig string `gorm:"type:jsonb"` // PostgreSQL JSONB类型，支持索引查询

	// 奖励
	RewardPoints int `gorm:"default:0"` // 奖励积分

	// 元数据
	IsHidden  bool `gorm:"default:false"` // 是否隐藏（隐藏成就在达成前不显示）
	IsActive  bool `gorm:"default:true"`  // 是否启用
	SortOrder int  `gorm:"default:0"`     // 排序顺序
}

// UserAchievement 用户成就表
type UserAchievement struct {
	gorm.Model
	UserID        uint        `gorm:"index;uniqueIndex:idx_user_achievement"`
	User          User        `gorm:"foreignKey:UserID"`
	AchievementID uint        `gorm:"index;uniqueIndex:idx_user_achievement"`
	Achievement   Achievement `gorm:"foreignKey:AchievementID"`

	UnlockedAt time.Time `gorm:"index"`         // 解锁时间
	Progress   int       `gorm:"default:0"`     // 进度值（用于显示进度条）
	IsNotified bool      `gorm:"default:false"` // 是否已通知用户

	// 防止重复授予：联合唯一索引 idx_user_achievement (user_id, achievement_id)
}

// AchievementProgress 成就进度追踪表（用于需要累积的成就）
type AchievementProgress struct {
	gorm.Model
	UserID          uint   `gorm:"index"`
	AchievementCode string `gorm:"size:64;index"` // 关联到Achievement.Code

	// 进度数据
	CurrentValue int       `gorm:"default:0"` // 当前进度值
	TargetValue  int       `gorm:"default:0"` // 目标值
	LastUpdated  time.Time `gorm:"index"`     // 最后更新时间

	// 时间窗口数据（用于时间相关成就）
	WindowStart *time.Time `gorm:"index"` // 时间窗口开始
	WindowEnd   *time.Time `gorm:"index"` // 时间窗口结束

	// 额外数据（JSON格式，用于存储复杂成就的中间状态）
	ExtraData string `gorm:"type:jsonb"`
}

// UserAchievementStat 用户成就统计表（用于排行榜等）
type UserAchievementStat struct {
	gorm.Model
	UserID            uint `gorm:"uniqueIndex"`
	User              User `gorm:"foreignKey:UserID"`
	TotalAchievements int  `gorm:"default:0"` // 总成就数
	TotalPoints       int  `gorm:"default:0"` // 总积分
	CommonCount       int  `gorm:"default:0"` // 普通成就数
	UncommonCount     int  `gorm:"default:0"` // 少见成就数
	RareCount         int  `gorm:"default:0"` // 稀有成就数
	EpicCount         int  `gorm:"default:0"` // 史诗成就数
	LegendaryCount    int  `gorm:"default:0"` // 传说成就数
}

func init() {
	modelRegister.Register(&Achievement{})
	modelRegister.Register(&UserAchievement{})
	modelRegister.Register(&AchievementProgress{})
	modelRegister.Register(&UserAchievementStat{})
}

// Response models

type AchievementResponse struct {
	ID           uint                `json:"id"`
	Code         string              `json:"code"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Icon         string              `json:"icon"`
	Type         AchievementType     `json:"type"`
	Category     AchievementCategory `json:"category"`
	Rarity       AchievementRarity   `json:"rarity"`
	RewardPoints int                 `json:"reward_points"`
	IsHidden     bool                `json:"is_hidden"`
	Progress     *int                `json:"progress,omitempty"`     // 当前进度
	TargetValue  *int                `json:"target_value,omitempty"` // 目标值
	UnlockedAt   *time.Time          `json:"unlocked_at,omitempty"`  // 解锁时间（如果已解锁）
	IsUnlocked   bool                `json:"is_unlocked"`
}

type UserAchievementListResponse struct {
	Unlocked   []AchievementResponse `json:"unlocked"`    // 已解锁
	InProgress []AchievementResponse `json:"in_progress"` // 进行中
	Locked     []AchievementResponse `json:"locked"`      // 未解锁（非隐藏）
}

type AchievementStatResponse struct {
	TotalAchievements int                       `json:"total_achievements"`
	UnlockedCount     int                       `json:"unlocked_count"`
	TotalPoints       int                       `json:"total_points"`
	RarityBreakdown   map[AchievementRarity]int `json:"rarity_breakdown"`
	RecentUnlocked    []AchievementResponse     `json:"recent_unlocked"` // 最近解锁的成就
}
