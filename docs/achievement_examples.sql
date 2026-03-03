-- 成就系统数据库示例 SQL
-- 这个文件包含一些示例成就配置，可以直接执行来添加更多成就

-- ========================================
-- 评论相关成就
-- ========================================

-- 深夜评论家
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
) ON CONFLICT (code) DO NOTHING;

-- 评论狂魔
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
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
    false,
    100,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 评论传说
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'legendary_commenter',
    '评论传说',
    '发布1000条评论',
    'counter',
    'comment',
    'legendary',
    '{"counter_target": 1000, "action": "comment"}',
    5000,
    true,
    false,
    101,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- ========================================
-- 点赞相关成就
-- ========================================

-- 人气王
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'super_popular',
    '人气王',
    '累计收到500个点赞',
    'counter',
    'like',
    'epic',
    '{"counter_target": 500, "action": "like_received"}',
    1000,
    true,
    false,
    110,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 万人迷
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'beloved',
    '万人迷',
    '累计收到1000个点赞',
    'counter',
    'like',
    'legendary',
    '{"counter_target": 1000, "action": "like_received"}',
    3000,
    true,
    false,
    111,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 慷慨的赞赏者
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'generous_liker',
    '慷慨的赞赏者',
    '给出100个点赞',
    'counter',
    'like',
    'uncommon',
    '{"counter_target": 100, "action": "like_given"}',
    100,
    true,
    false,
    120,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- ========================================
-- 社交相关成就
-- ========================================

-- 活跃交流者
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'active_replier',
    '活跃交流者',
    '发布50条回复',
    'counter',
    'social',
    'uncommon',
    '{"counter_target": 50, "action": "reply"}',
    150,
    true,
    false,
    200,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 课程探索者
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'course_explorer',
    '课程探索者',
    '在10个不同的课程组发布评论',
    'complex',
    'explore',
    'rare',
    '{"unique_course_groups": 10, "action": "comment"}',
    300,
    true,
    false,
    300,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- ========================================
-- 时间相关成就
-- ========================================

-- 周末战士
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'weekend_warrior',
    '周末战士',
    '在一个周末（周六和周日）内发布10条评论',
    'time_based',
    'time',
    'uncommon',
    '{"time_window": "48h", "required_count": 10, "action": "comment", "weekend_only": true}',
    150,
    true,
    false,
    210,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 马拉松评论家
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'marathon_commenter',
    '马拉松评论家',
    '连续7天每天至少发布1条评论',
    'time_based',
    'time',
    'rare',
    '{"consecutive_days": 7, "daily_minimum": 1, "action": "comment"}',
    500,
    true,
    false,
    211,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- ========================================
-- 特殊成就
-- ========================================

-- 早起的鸟儿
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'early_bird',
    '早起的鸟儿',
    '在早上6-8点发布评论',
    'one_time',
    'special',
    'uncommon',
    '{"time_range": ["06:00", "08:00"], "action": "comment"}',
    50,
    true,
    false,
    400,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 完美主义者
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'perfectionist',
    '完美主义者',
    '发布一条评论获得100个点赞',
    'complex',
    'special',
    'epic',
    '{"single_comment_likes": 100}',
    800,
    true,
    true,
    401,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 意见领袖
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'opinion_leader',
    '意见领袖',
    '发布的评论累计被1000人阅读（需实现阅读追踪）',
    'complex',
    'special',
    'legendary',
    '{"total_views": 1000}',
    2000,
    true,
    true,
    402,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- ========================================
-- 隐藏成就（达成前不显示）
-- ========================================

-- 幸运数字
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'lucky_number',
    '幸运数字',
    '发布第88条评论',
    'counter',
    'special',
    'rare',
    '{"counter_target": 88, "action": "comment", "exact_match": true}',
    200,
    true,
    true,
    500,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- 深夜传说
INSERT INTO achievements (
    code, name, description, type, category, rarity,
    trigger_config, reward_points, is_active, is_hidden, sort_order,
    created_at, updated_at
) VALUES (
    'midnight_legend',
    '深夜传说',
    '在凌晨0点整（前后5分钟内）发布评论',
    'one_time',
    'special',
    'legendary',
    '{"exact_time": "00:00", "tolerance_minutes": 5, "action": "comment"}',
    1000,
    true,
    true,
    501,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- ========================================
-- 查询成就相关的统计信息
-- ========================================

-- 查看所有成就
-- SELECT * FROM achievements ORDER BY sort_order, id;

-- 查看某个用户的成就
-- SELECT u.nick_name, a.name, a.rarity, ua.unlocked_at 
-- FROM user_achievements ua
-- JOIN users u ON ua.user_id = u.id
-- JOIN achievements a ON ua.achievement_id = a.id
-- WHERE ua.user_id = 1
-- ORDER BY ua.unlocked_at DESC;

-- 查看稀有成就持有者
-- SELECT a.name, a.rarity, COUNT(ua.user_id) as holders
-- FROM achievements a
-- LEFT JOIN user_achievements ua ON a.id = ua.achievement_id
-- WHERE a.rarity IN ('epic', 'legendary')
-- GROUP BY a.id, a.name, a.rarity
-- ORDER BY holders DESC;

-- 查看成就排行榜
-- SELECT u.nick_name, uas.total_achievements, uas.total_points,
--        uas.legendary_count, uas.epic_count, uas.rare_count
-- FROM user_achievement_stats uas
-- JOIN users u ON uas.user_id = u.id
-- ORDER BY uas.total_points DESC
-- LIMIT 20;
