# CourseBench 成就系统设计文档

## 概述

成就系统是一个游戏化机制，用于激励用户参与和贡献。本文档描述了成就系统的设计、实现和使用方法。

## 核心设计原则

### 1. 判断与触发机制

#### 事件驱动触发
- **同步检查**: 关键业务操作（如发布评论）后同步检查简单成就
- **异步检查**: 复杂成就和统计相关成就使用goroutine异步检查，避免阻塞主流程
- **定时任务**: 需要聚合统计的成就可通过定时任务批量处理

#### 双重验证
- 使用数据库唯一约束防止重复授予
- 在授予成就前检查用户是否已拥有
- 使用事务确保数据一致性

### 2. 安全性保障

#### 事务保护
```go
err := db.Transaction(func(tx *gorm.DB) error {
    // 在事务中授予成就
    // 更新用户统计
    // 更新用户积分
    return nil
})
```

#### 幂等性保证
```go
// 使用 ON CONFLICT DO NOTHING 实现幂等性
err := tx.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "user_id"}, {Name: "achievement_id"}},
    DoNothing: true,
}).Create(&userAchievement).Error
```

#### 服务端校验
- 所有成就判断在后端完成
- 不信任前端数据
- 使用实际数据库统计进行验证

### 3. 存储策略

#### 数据库表结构

1. **achievements** - 成就定义表
   - 存储成就的元数据、规则、奖励等
   - 使用JSONB存储灵活的触发条件配置

2. **user_achievements** - 用户成就表
   - 记录用户已解锁的成就
   - 使用联合唯一索引 (user_id, achievement_id) 防止重复

3. **achievement_progress** - 成就进度表
   - 追踪未完成成就的进度
   - 支持时间窗口限制

4. **user_achievement_stats** - 用户成就统计表
   - 预计算的统计数据
   - 用于排行榜等高频查询场景

#### 缓存策略
```go
// 可选：使用Redis缓存用户成就列表
// key: "user:achievements:{user_id}"
// TTL: 1小时
```

## 成就类型

### 1. 计数类型成就 (Counter)
简单的累积计数，达到目标值即解锁。

**触发配置示例**:
```json
{
    "counter_target": 10,
    "action": "comment"
}
```

**使用场景**:
- 发布X条评论
- 收到X个点赞
- 回复X次

### 2. 时间相关成就 (Time-Based)
在特定时间窗口内完成指定次数的行为。

**触发配置示例**:
```json
{
    "time_window": "24h",
    "required_count": 5,
    "action": "comment"
}
```

**使用场景**:
- 24小时内发布5条评论
- 一周内收到50个点赞

### 3. 复合条件成就 (Complex)
需要满足多个条件的成就。

**触发配置示例**:
```json
{
    "conditions": [
        {"action": "comment", "min_count": 10},
        {"action": "like_received", "min_count": 50}
    ],
    "logic": "AND"
}
```

**使用场景**:
- 同时拥有10条评论和50个点赞
- 在不同课程组发布评论

### 4. 一次性成就 (One-Time)
触发即获得的特殊成就。

**使用场景**:
- 首次登录
- 完成个人资料
- 邀请好友

## API 接口

### 用户接口

#### 获取当前用户成就列表
```http
GET /v1/achievement/list
Authorization: Bearer {token}
```

**响应**:
```json
{
    "error": false,
    "data": {
        "unlocked": [
            {
                "id": 1,
                "code": "first_comment",
                "name": "初次发声",
                "description": "发布你的第一条评论",
                "icon": "/icons/first_comment.png",
                "type": "counter",
                "category": "comment",
                "rarity": "common",
                "reward_points": 10,
                "is_unlocked": true,
                "unlocked_at": "2026-03-01T10:00:00Z",
                "progress": 100
            }
        ],
        "in_progress": [
            {
                "id": 2,
                "code": "comment_veteran_10",
                "name": "评论新手",
                "description": "发布10条评论",
                "progress": 5,
                "target_value": 10,
                "is_unlocked": false
            }
        ],
        "locked": [...]
    }
}
```

#### 获取成就统计
```http
GET /v1/achievement/stats
Authorization: Bearer {token}
```

#### 获取其他用户的公开成就
```http
GET /v1/achievement/user/:id
```

### 管理员接口

#### 初始化默认成就
```http
POST /v1/achievement/admin/initialize
Authorization: Bearer {admin_token}
```

#### 重新计算用户成就
```http
POST /v1/achievement/admin/recalculate?user_id=123
Authorization: Bearer {admin_token}
```

## 使用指南

### 1. 系统初始化

启动后端服务后，管理员需要初始化默认成就：

```bash
curl -X POST http://localhost:8080/v1/achievement/admin/initialize \
  -H "Authorization: Bearer {admin_token}"
```

### 2. 在业务逻辑中触发成就检查

#### 在控制器中同步触发
```go
// 发布评论后
err := queries.CheckAndGrantAchievements(db, userID, "comment", 1)
```

#### 异步触发（推荐）
```go
// 发布评论后，不阻塞主流程
go func() {
    db := database.GetDB()
    _ = queries.CheckAndGrantAchievements(db, userID, "comment", 1)
}()
```

### 3. 添加自定义成就

#### 直接插入数据库
```sql
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, sort_order
) VALUES (
    'super_commenter',
    '评论狂魔',
    '发布500条评论',
    'counter',
    'comment',
    'epic',
    '{"counter_target": 500, "action": "comment"}',
    1000,
    true,
    100
);
```

#### 使用代码创建
```go
achievement := &models.Achievement{
    Code:        "super_commenter",
    Name:        "评论狂魔",
    Description: "发布500条评论",
    Type:        models.AchievementTypeCounter,
    Category:    models.CategoryComment,
    Rarity:      models.RarityEpic,
    TriggerConfig: `{"counter_target": 500, "action": "comment"}`,
    RewardPoints: 1000,
    IsActive:    true,
    SortOrder:   100,
}
db.Create(achievement)
```

### 4. 支持的动作类型 (actions)

- `comment` - 发布评论
- `like_received` - 收到点赞
- `like_given` - 给出点赞
- `reply` - 发布回复
- `course_explored` - 探索课程

可以根据需要添加更多动作类型。

## 最佳实践

### 1. 性能优化

#### 使用异步处理
```go
// 避免阻塞主请求
go func() {
    db := database.GetDB()
    _ = queries.CheckAndGrantAchievements(db, userID, action, value)
}()
```

#### 批量处理
```go
// 对于需要统计的成就，使用定时任务批量处理
func BatchCheckAchievements() {
    users := getAllActiveUsers()
    for _, user := range users {
        stats := queries.GetActualUserStatistics(db, user.ID)
        // 检查各类成就
    }
}
```

#### 使用缓存
```go
// 缓存用户成就列表
func GetUserAchievementsWithCache(userID uint) (*models.UserAchievementListResponse, error) {
    cacheKey := fmt.Sprintf("user:achievements:%d", userID)
    
    // 尝试从Redis获取
    if cached := redis.Get(cacheKey); cached != nil {
        return cached, nil
    }
    
    // 从数据库获取
    achievements := queries.GetUserAchievements(db, userID)
    
    // 写入缓存
    redis.Set(cacheKey, achievements, 1*time.Hour)
    
    return achievements, nil
}
```

### 2. 安全性

#### 防止作弊
```go
// 1. 使用速率限制
if !checkRateLimit(userID, action) {
    return errors.New("操作过于频繁")
}

// 2. 验证数据真实性
stats := queries.GetActualUserStatistics(db, userID)
if stats["total_comments"] < progress.CurrentValue {
    // 数据异常，可能存在作弊
    log.Warn("Achievement data mismatch for user", userID)
}
```

#### 权限控制
```go
// 管理员接口必须检查权限
if !user.IsAdmin {
    return errors.New(errors.PermissionDenied)
}
```

### 3. 监控与维护

#### 日志记录
```go
log.Printf("Achievement granted: user=%d, achievement=%s", userID, achievement.Code)
```

#### 数据修复
```go
// 定期检查数据一致性
func VerifyAchievementData() {
    users := getAllUsers()
    for _, user := range users {
        actualStats := queries.GetActualUserStatistics(db, user.ID)
        recordedStats := getUserAchievementProgress(db, user.ID)
        
        if actualStats != recordedStats {
            log.Warn("Data mismatch for user", user.ID)
            // 触发重新计算
            recalculateUserAchievements(user.ID)
        }
    }
}
```

## 扩展建议

### 1. 通知系统
```go
// 成就解锁后发送通知
func notifyAchievementUnlocked(userID uint, achievement models.Achievement) {
    // WebSocket实时通知
    ws.Send(userID, "achievement_unlocked", achievement)
    
    // 邮件通知（可选）
    if achievement.Rarity == models.RarityLegendary {
        mail.SendAchievementEmail(userID, achievement)
    }
}
```

### 2. 成就展示墙
```go
// 用户可以选择展示的成就
type UserAchievementDisplay struct {
    UserID        uint
    AchievementID uint
    DisplayOrder  int
    IsShowcase    bool
}
```

### 3. 成就组合
```go
// 完成一组成就后解锁特殊成就
type AchievementSet struct {
    Code              string
    Name              string
    RequiredAchievements []string
    RewardAchievement string
}
```

### 4. 季节性成就
```go
// 限时成就
type SeasonalAchievement struct {
    Achievement
    StartTime time.Time
    EndTime   time.Time
}
```

## 故障排查

### 问题：成就没有自动授予
**解决方案**:
1. 检查成就是否启用 (`is_active = true`)
2. 检查触发配置是否正确
3. 查看日志确认触发函数是否被调用
4. 使用管理员接口重新计算

### 问题：重复授予成就
**解决方案**:
1. 检查唯一索引是否存在
2. 确认使用了 `ON CONFLICT DO NOTHING`
3. 检查是否有并发问题

### 问题：进度不准确
**解决方案**:
1. 使用 `GetActualUserStatistics` 获取真实数据
2. 使用管理员接口重新计算
3. 检查是否有数据回滚导致不一致

## 总结

这个成就系统提供了：
- ✅ 灵活的成就类型定义
- ✅ 安全的授予机制（幂等性、事务保护）
- ✅ 高性能的异步处理
- ✅ 完善的进度追踪
- ✅ 易于扩展的架构

通过合理使用本系统，可以有效提升用户参与度和平台活跃度。
