# 成就系统测试指南

## 测试环境准备

### 1. 启动后端服务
```bash
cd backend
go run cmd/coursebench-backend/main.go
```

### 2. 初始化成就系统
```bash
# 替换YOUR_ADMIN_TOKEN为实际的管理员token
curl -X POST http://localhost:8080/v1/achievement/admin/initialize \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 3. 验证成就已创建
```sql
SELECT code, name, type, category, rarity FROM achievements ORDER BY sort_order;
```

## 功能测试

### 测试1: 发布评论触发成就

#### 测试步骤：
1. 使用一个新用户登录
2. 发布第一条评论
3. 查看成就列表

```bash
# 发布评论
curl -X POST http://localhost:8080/v1/comment/post \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "group": 1,
    "title": "测试评论",
    "content": "这是一条测试评论",
    "semester": 20241,
    "scores": [5,5,5,5],
    "student_score_ranking": 3,
    "is_anonymous": false
  }'

# 等待1-2秒（异步处理）

# 查看成就列表
curl http://localhost:8080/v1/achievement/list \
  -H "Authorization: Bearer USER_TOKEN"
```

#### 预期结果：
- 应该获得 "初次发声" (first_comment) 成就
- 返回的JSON中 `unlocked` 数组应包含该成就
- 用户reward积分应增加10分

#### 验证SQL：
```sql
SELECT * FROM user_achievements 
WHERE user_id = YOUR_USER_ID AND achievement_id = (
    SELECT id FROM achievements WHERE code = 'first_comment'
);

SELECT reward FROM users WHERE id = YOUR_USER_ID;
```

### 测试2: 累积评论触发成就

#### 测试步骤：
1. 连续发布10条评论（在不同的课程组）
2. 查看成就进度
3. 验证是否获得 "评论新手" 成就

```bash
# 循环发布10条评论（需要修改group为不同的课程组ID）
for i in {1..10}; do
  curl -X POST http://localhost:8080/v1/comment/post \
    -H "Authorization: Bearer USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"group\": $i,
      \"title\": \"测试评论 $i\",
      \"content\": \"这是第 $i 条测试评论\",
      \"semester\": 20241,
      \"scores\": [5,5,5,5],
      \"student_score_ranking\": 3
    }"
  sleep 1
done

# 查看成就列表
curl http://localhost:8080/v1/achievement/list \
  -H "Authorization: Bearer USER_TOKEN"
```

#### 预期结果：
- 应该获得 "评论新手" (comment_veteran_10) 成就
- 用户reward积分应增加50分
- 如果发布了50条，应获得 "评论达人" 成就

#### 验证进度：
```sql
SELECT ap.achievement_code, ap.current_value, ap.target_value
FROM achievement_progress ap
WHERE ap.user_id = YOUR_USER_ID;
```

### 测试3: 点赞触发成就

#### 测试步骤：
1. 用户A发布评论
2. 用户B给该评论点赞
3. 多个用户给用户A的评论点赞
4. 检查用户A是否获得点赞相关成就

```bash
# 用户A发布评论
curl -X POST http://localhost:8080/v1/comment/post \
  -H "Authorization: Bearer USER_A_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "group": 1,
    "title": "优质评论",
    "content": "这是一条很有价值的评论",
    "semester": 20241,
    "scores": [5,5,5,5],
    "student_score_ranking": 3
  }'

# 记录返回的comment_id

# 用户B点赞
curl -X POST http://localhost:8080/v1/comment/like \
  -H "Authorization: Bearer USER_B_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": COMMENT_ID,
    "status": 1
  }'

# 重复10次以上，用不同用户账号

# 查看用户A的成就
curl http://localhost:8080/v1/achievement/list \
  -H "Authorization: Bearer USER_A_TOKEN"
```

#### 预期结果：
- 收到10个点赞后，应获得 "小有名气" (popular_10) 成就
- 收到50个点赞后，应获得 "人气之星" (popular_50) 成就

### 测试4: 时间窗口成就

#### 测试步骤：
1. 在24小时内发布5条评论
2. 检查是否获得 "活跃的一天" 成就

```bash
# 快速发布5条评论
for i in {1..5}; do
  curl -X POST http://localhost:8080/v1/comment/post \
    -H "Authorization: Bearer USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"group\": $i,
      \"title\": \"快速评论 $i\",
      \"content\": \"测试时间窗口成就\",
      \"semester\": 20241,
      \"scores\": [5,5,5,5],
      \"student_score_ranking\": 3
    }"
  sleep 2
done

# 查看成就
curl http://localhost:8080/v1/achievement/list \
  -H "Authorization: Bearer USER_TOKEN"
```

#### 预期结果：
- 应该获得 "活跃的一天" (active_day) 成就

#### 验证时间窗口：
```sql
SELECT ap.achievement_code, ap.current_value, ap.target_value,
       ap.window_start, ap.window_end
FROM achievement_progress ap
WHERE ap.user_id = YOUR_USER_ID 
  AND ap.achievement_code = 'active_day';
```

### 测试5: 成就统计

#### 测试步骤：
```bash
# 获取成就统计
curl http://localhost:8080/v1/achievement/stats \
  -H "Authorization: Bearer USER_TOKEN"
```

#### 预期结果：
```json
{
  "error": false,
  "data": {
    "total_achievements": 3,
    "unlocked_count": 3,
    "total_points": 110,
    "rarity_breakdown": {
      "common": 2,
      "uncommon": 1,
      "rare": 0,
      "epic": 0,
      "legendary": 0
    },
    "recent_unlocked": [...]
  }
}
```

### 测试6: 查看其他用户成就

#### 测试步骤：
```bash
# 查看用户ID为123的成就
curl http://localhost:8080/v1/achievement/user/123
```

#### 预期结果：
- 如果用户是匿名的，只显示已解锁的成就（不显示进度）
- 如果用户不是匿名的，显示所有公开信息

### 测试7: 幂等性测试（防止重复授予）

#### 测试步骤：
1. 手动触发成就授予两次
2. 检查数据库确保只有一条记录

```sql
-- 尝试插入重复的用户成就
INSERT INTO user_achievements (user_id, achievement_id, unlocked_at, created_at, updated_at)
VALUES (1, 1, NOW(), NOW(), NOW())
ON CONFLICT (user_id, achievement_id) DO NOTHING;

-- 验证只有一条记录
SELECT COUNT(*) FROM user_achievements 
WHERE user_id = 1 AND achievement_id = 1;
-- 应该返回 1
```

### 测试8: 管理员重新计算成就

#### 测试步骤：
```bash
# 重新计算用户123的成就
curl -X POST "http://localhost:8080/v1/achievement/admin/recalculate?user_id=123" \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

#### 预期结果：
```json
{
  "error": false,
  "data": {
    "recalculated": true,
    "stats": {
      "total_comments": 10,
      "total_likes_received": 25,
      "total_likes_given": 15,
      "total_replies": 5
    }
  }
}
```

## 压力测试

### 测试9: 并发触发成就

#### 测试目的：
验证在高并发情况下，成就系统是否能正确处理，不会出现重复授予或数据不一致。

#### 测试脚本：
```bash
#!/bin/bash
# concurrent_test.sh

for i in {1..100}; do
  (
    curl -X POST http://localhost:8080/v1/comment/post \
      -H "Authorization: Bearer USER_TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"group\": $((i % 10 + 1)),
        \"title\": \"并发测试 $i\",
        \"content\": \"测试并发场景\",
        \"semester\": 20241,
        \"scores\": [5,5,5,5],
        \"student_score_ranking\": 3
      }" &
  )
done
wait

echo "所有请求已完成"
```

#### 验证：
```sql
-- 检查是否有重复授予的成就
SELECT user_id, achievement_id, COUNT(*) as count
FROM user_achievements
WHERE deleted_at IS NULL
GROUP BY user_id, achievement_id
HAVING COUNT(*) > 1;
-- 应该返回空结果

-- 检查统计数据是否一致
SELECT u.id, 
       COUNT(c.id) as actual_comments,
       COALESCE(ap.current_value, 0) as recorded_progress
FROM users u
LEFT JOIN comments c ON u.id = c.user_id
LEFT JOIN achievement_progress ap ON u.id = ap.user_id 
    AND ap.achievement_code LIKE '%comment%'
WHERE u.id = YOUR_USER_ID
GROUP BY u.id, ap.current_value;
```

## 性能测试

### 测试10: 响应时间测试

#### 测试步骤：
使用Apache Bench或类似工具测试API响应时间

```bash
# 测试成就列表API的性能
ab -n 1000 -c 10 \
  -H "Authorization: Bearer USER_TOKEN" \
  http://localhost:8080/v1/achievement/list

# 测试成就统计API的性能
ab -n 1000 -c 10 \
  -H "Authorization: Bearer USER_TOKEN" \
  http://localhost:8080/v1/achievement/stats
```

#### 预期结果：
- 平均响应时间 < 100ms
- 99%请求的响应时间 < 200ms
- 无错误响应

## 数据一致性验证

### 测试11: 验证统计数据准确性

```sql
-- 验证评论数统计
SELECT u.id, u.nick_name,
       COUNT(c.id) as actual_comments,
       COALESCE(ap.current_value, 0) as recorded_progress
FROM users u
LEFT JOIN comments c ON u.id = c.user_id AND c.deleted_at IS NULL
LEFT JOIN achievement_progress ap ON u.id = ap.user_id 
    AND ap.achievement_code = 'comment_veteran_10'
GROUP BY u.id, u.nick_name, ap.current_value
HAVING COUNT(c.id) != COALESCE(ap.current_value, 0);

-- 验证点赞数统计
SELECT u.id, u.nick_name,
       SUM(c.like) as actual_likes,
       COALESCE(ap.current_value, 0) as recorded_progress
FROM users u
LEFT JOIN comments c ON u.id = c.user_id AND c.deleted_at IS NULL
LEFT JOIN achievement_progress ap ON u.id = ap.user_id 
    AND ap.achievement_code = 'popular_10'
GROUP BY u.id, u.nick_name, ap.current_value
HAVING SUM(c.like) != COALESCE(ap.current_value, 0);

-- 验证成就统计表
SELECT u.id,
       COUNT(ua.id) as actual_achievements,
       COALESCE(uas.total_achievements, 0) as recorded_achievements
FROM users u
LEFT JOIN user_achievements ua ON u.id = ua.user_id AND ua.deleted_at IS NULL
LEFT JOIN user_achievement_stats uas ON u.id = uas.user_id
GROUP BY u.id, uas.total_achievements
HAVING COUNT(ua.id) != COALESCE(uas.total_achievements, 0);
```

## 边界条件测试

### 测试12: 边界值测试

1. **测试目标值边界**：
   - 发布刚好10条评论，验证是否触发成就
   - 发布9条评论，验证是否不触发

2. **测试时间窗口边界**：
   - 在时间窗口边界发布评论
   - 超过时间窗口后验证是否重置

3. **测试数值溢出**：
   - 发布大量评论（>1000）
   - 验证计数器是否正常

## 错误处理测试

### 测试13: 异常情况处理

1. **数据库连接失败**：
   - 模拟数据库断开
   - 验证错误是否被正确捕获

2. **无效的成就配置**：
   ```sql
   -- 插入无效配置
   INSERT INTO achievements (code, name, type, trigger_config, created_at, updated_at)
   VALUES ('invalid', 'Invalid', 'counter', '{invalid json', NOW(), NOW());
   ```
   - 验证系统是否能正常处理

3. **用户不存在**：
   ```bash
   curl http://localhost:8080/v1/achievement/user/99999
   ```
   - 应返回404或相应错误

## 测试检查清单

- [ ] 初始化成就系统成功
- [ ] 发布评论自动触发成就检查
- [ ] 点赞自动触发成就检查
- [ ] 计数类型成就正确授予
- [ ] 时间窗口成就正确授予
- [ ] 成就不会重复授予（幂等性）
- [ ] 成就统计数据准确
- [ ] 成就积分正确累加到用户
- [ ] 隐藏成就在达成前不显示
- [ ] 匿名用户的成就隐私保护
- [ ] 管理员接口权限检查
- [ ] 管理员重新计算功能正常
- [ ] API响应时间符合预期
- [ ] 并发场景下数据一致
- [ ] 数据库约束正常工作
- [ ] 错误情况正确处理

## 回归测试

在每次代码变更后，至少运行以下核心测试：
1. 测试1 (发布评论触发成就)
2. 测试3 (点赞触发成就)
3. 测试7 (幂等性测试)
4. 测试11 (数据一致性验证)

## 自动化测试建议

建议创建自动化测试脚本：

```bash
# run_achievement_tests.sh
#!/bin/bash

echo "=== 成就系统自动化测试 ==="

# 1. 初始化
echo "1. 初始化成就系统..."
# ... 测试代码

# 2. 功能测试
echo "2. 运行功能测试..."
# ... 测试代码

# 3. 性能测试
echo "3. 运行性能测试..."
# ... 测试代码

# 4. 一致性验证
echo "4. 验证数据一致性..."
# ... 测试代码

echo "=== 测试完成 ==="
```

## 问题报告模板

如果发现问题，请按以下格式报告：

```
**问题描述**：
简要描述问题

**复现步骤**：
1. ...
2. ...
3. ...

**预期结果**：
应该发生什么

**实际结果**：
实际发生了什么

**环境信息**：
- Go版本：
- 数据库版本：
- 相关日志：

**相关数据**：
```sql
-- 相关的SQL查询和结果
```
