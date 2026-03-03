# 成就系统快速开始

## 已实现的文件

### 数据模型
- [backend/pkg/models/achievement.go](backend/pkg/models/achievement.go) - 成就系统数据模型

### 业务逻辑
- [backend/pkg/queries/achievement.go](backend/pkg/queries/achievement.go) - 成就检查、授予、查询逻辑

### API控制器
- [backend/internal/controllers/achievement/achievement.go](backend/internal/controllers/achievement/achievement.go) - 成就API
- [backend/internal/controllers/achievements.go](backend/internal/controllers/achievements.go) - 路由定义

### 数据库迁移
- [backend/pkg/database/upgrade/upgrade.go](backend/pkg/database/upgrade/upgrade.go) - 添加了v4到v5的升级

### 集成点
- [backend/internal/controllers/comments/post.go](backend/internal/controllers/comments/post.go) - 评论后触发成就
- [backend/internal/controllers/comments/like.go](backend/internal/controllers/comments/like.go) - 点赞后触发成就
- [backend/internal/fiber/route.go](backend/internal/fiber/route.go) - 添加成就路由

## 快速开始

### 1. 启动服务
```bash
cd backend
go run cmd/coursebench-backend/main.go
```

数据库会自动升级到v5并创建成就相关表。

### 2. 初始化默认成就
```bash
curl -X POST http://localhost:8080/v1/achievement/admin/initialize \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 3. 测试成就系统

#### 发布评论触发成就
```bash
# 发布评论后，会自动检查"初次发声"、"评论新手"等成就
curl -X POST http://localhost:8080/v1/comment/post \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "group": 1,
    "title": "test",
    "content": "test content",
    "semester": 20231,
    "scores": [5,5,5,5],
    "student_score_ranking": 3
  }'
```

#### 查看成就列表
```bash
curl http://localhost:8080/v1/achievement/list \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 查看成就统计
```bash
curl http://localhost:8080/v1/achievement/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 默认成就列表

| 代码 | 名称 | 描述 | 类型 | 稀有度 | 奖励积分 |
|------|------|------|------|--------|---------|
| first_comment | 初次发声 | 发布你的第一条评论 | 计数 | 普通 | 10 |
| comment_veteran_10 | 评论新手 | 发布10条评论 | 计数 | 普通 | 50 |
| comment_veteran_50 | 评论达人 | 发布50条评论 | 计数 | 少见 | 200 |
| comment_master_100 | 评论大师 | 发布100条评论 | 计数 | 稀有 | 500 |
| popular_10 | 小有名气 | 累计收到10个点赞 | 计数 | 普通 | 30 |
| popular_50 | 人气之星 | 累计收到50个点赞 | 计数 | 少见 | 100 |
| popular_100 | 万众瞩目 | 累计收到100个点赞 | 计数 | 稀有 | 300 |
| active_day | 活跃的一天 | 24小时内发布5条评论 | 时间 | 少见 | 100 |
| active_week | 活跃的一周 | 7天内发布20条评论 | 时间 | 稀有 | 300 |

## 如何添加新成就

### 方法1：直接插入数据库
```sql
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'night_owl',
    '夜猫子',
    '在凌晨2-5点发布评论',
    'one_time',
    'special',
    'uncommon',
    '{"time_range": ["02:00", "05:00"], "action": "comment"}',
    50,
    true,
    false,
    30,
    NOW(),
    NOW()
);
```

### 方法2：在代码中创建
编辑 [backend/pkg/queries/achievement.go](backend/pkg/queries/achievement.go) 的 `InitializeDefaultAchievements` 函数，添加新成就定义。

## 在业务逻辑中触发成就

### 异步触发（推荐）
```go
go func() {
    db := database.GetDB()
    _ = queries.CheckAndGrantAchievements(db, userID, "action_name", value)
}()
```

### 同步触发
```go
err := queries.CheckAndGrantAchievements(db, userID, "action_name", value)
if err != nil {
    log.Error(err)
}
```

### 支持的动作类型
- `comment` - 发布评论
- `like_received` - 收到点赞
- `like_given` - 给出点赞（需要添加集成）
- `reply` - 发布回复（需要添加集成）

## 数据修复

如果成就数据不准确，可以使用管理员接口重新计算：

```bash
curl -X POST "http://localhost:8080/v1/achievement/admin/recalculate?user_id=123" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## 架构设计要点

### ✅ 安全性
- **幂等性保证**: 使用数据库唯一约束 + ON CONFLICT DO NOTHING
- **事务保护**: 所有成就授予操作在事务中完成
- **服务端验证**: 不信任前端数据，所有判断在后端完成

### ✅ 性能
- **异步处理**: 使用goroutine异步检查成就，不阻塞主流程
- **索引优化**: 为高频查询字段添加索引
- **可扩展缓存**: 预留Redis缓存接口

### ✅ 可扩展性
- **灵活配置**: 触发条件使用JSONB存储，支持复杂规则
- **多种类型**: 支持计数、时间、复合、一次性等多种成就类型
- **易于添加**: 新增成就只需插入数据，无需修改代码

## 详细文档

查看完整设计文档: [backend/docs/ACHIEVEMENT_SYSTEM.md](backend/docs/ACHIEVEMENT_SYSTEM.md)

## 后续优化建议

1. **通知系统**: 成就解锁时实时通知用户（WebSocket）
2. **成就展示**: 用户可以选择展示的成就
3. **统计面板**: 成就排行榜、稀有成就持有者等
4. **复杂成就**: 实现复合条件成就的检查逻辑
5. **定时任务**: 使用cron定时检查时间相关成就
6. **缓存优化**: 集成Redis缓存高频查询数据
