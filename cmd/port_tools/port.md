# 数据迁移日志

## Step 1：从 ELRC 获取所有课程信息

注意，虽然大部分的 teachers 和 teacher_names 是一一对应的，但是也有少数课程的 teachers 和 teacher_names 不是一一对应的！会出顺序问题。

## Step 2：修正课程教师关系

只以 teachers 为准，teacher_names 从 OA 平台获取

## Step 3：生成教师与课程的统一关系

