# CourseBench 命令行工具

该目录包含一系列用于管理 CourseBench 后端数据的命令行工具。

在运行任何命令之前，请确保已经正确配置了 `config.json` 文件，该文件可以放在 `/etc/coursebench` 或者 `.`。

## 命令用法

### 导入课程数据

此命令用于从一个目录中批量导入课程信息。它会遍历指定目录下的所有文件，并将每个文件作为 CSV 导入。

**命令:**
```bash
go run cmd/cmd_tools/main.go import_course <directory_path>
```

**CSV 文件格式:**
该工具会跳过 CSV 文件的第一行（表头）。数据应从第二行开始，并至少包含以下列：
- `第3列`: 课程名称 (string)
- `第4列`: 课程代码 (string)
- `第5列`: 学分 (float)
- `第11列`: 开课学院 (string)
- `第13列`: 教师姓名列表 (JSON 数组格式的字符串, e.g., `'["张三", "李四"]'`)
- `第14列`: 教师 EAMS ID 列表 (JSON 数组格式的字符串, e.g., `'[123, 456]'`)

**功能:**
- 如果课程代码不存在，则会创建新的课程记录。
- 如果教师姓名不存在，则会创建新的教师记录。
- 根据教师阵容创建新的课程组，并避免创建重复的教师组。

---

### 更新教师信息

此命令用于根据 CSV 文件更新已存在教师的详细信息。

**注意：** 此命令**不会**创建新教师，仅会更新数据库中已存在的教师记录。

**命令:**
```bash
go run cmd/cmd_tools/main.go import_teacher <file_path>
```

**CSV 文件格式:**
文件需要包含以下6列，无表头：
1. `姓名` (string)
2. `照片URL` (string)
3. `职称` (string)
4. `邮箱` (string)
5. `学院` (string)
6. `简介` (string)

**功能:**
- 通过 `姓名` 查找教师，并更新其 `照片URL`, `职称`, `邮箱`, `学院`, 和 `简介` 字段。

---

### 设置/取消管理员

这些命令用于授予或撤销用户的管理员权限。

**命令:**
```bash
# 设置为普通管理员
go run cmd/cmd_tools/main.go set_admin <user_id>

# 取消普通管理员
go run cmd/cmd_tools/main.go unset_admin <user_id>

# 设置为社区管理员
go run cmd/cmd_tools/main.go set_community_admin <user_id>

# 取消社区管理员
go run cmd/cmd_tools/main.go unset_community_admin <user_id>
```

**功能:**
- `set_admin` / `unset_admin`: 修改用户在 `users` 表中的 `is_admin` 字段。
- `set_community_admin` / `unset_community_admin`: 修改用户在 `users` 表中的 `is_community_admin` 字段。
- 一个用户不能同时是普通管理员和社区管理员。

---

### 清除用户数据

**警告：这是一个危险的破坏性操作！**

此命令会删除所有用户生成的内容。为防止误操作，必须在命令最后加上 `Yes_Confirm` 参数。

**命令:**
```bash
go run cmd/cmd_tools/main.go clear_userdata Yes_Confirm
```

**功能:**
- 删除 `users`, `comments`, `comment_likes` 表中的所有记录。
- 重置 `courses` 和 `course_groups` 表中所有记录的 `comment_count` 和 `scores` 字段。

## 数据库结构

以下是这些工具主要交互的数据库表的结构推断：

- **`courses`** 课程信息
  - `id`: 主键
  - `code`: 课程代码 (唯一, 不含班级代码)
  - `name`: 课程名称
  - `institute`: 开课学院
  - `credit`: 学分
  - `comment_count`: 评论数
  - `scores`: 四维评分数组

- **`teachers`** 教师信息
  - `id`: 主键
  - `name`: 教师姓名
  - `eams_id`: EAMS ID
  - `photo`: 照片URL
  - `job`: 职称
  - `email`: 邮箱
  - `institute`: 学院
  - `introduction`: 简介
  - `uni_id`: 教师工号


- **`course_groups`** 课程组(细分课程到班/学期)
  - `id`: 主键
  - `course_id`: 外键，关联 `courses.id`
  - `comment_count`: 评论数
  - `scores`: 评分数组

- **`coursegroup_teachers`** (中间表) 课程组与教师的关系
  - `course_group_id`: 外键，关联 `course_groups.id`
  - `teacher_id`: 外键，关联 `teachers.id`

- **`users`** 用户信息
  - `id`: 主键
  - `is_admin`: 布尔值，是否为普通管理员
  - `is_community_admin`: 布尔值，是否为社区管理员

- **`comments`** 用户评论
  - (存储用户发表的评论)

- **`comment_likes`** 用户点赞记录
  - (存储用户对评论的点赞记录) 

## 迁移日志

### `import_and_fix_teacher_uniid`

这是根据 ELRC 中获取到的 `teachers.json` 文件导入教师信息，并修正教师工号（`uni_id`）的命令。

这条命令会**删除所有**原先的 EAMS 教师记录，并重新导入教师信息。

但是由于 ELRC 仍然有很多课程信息不全，会有一些老师的 UniID 仍然是 null，需要手动更正（写入 teachers.json 即可），并且将 EamsID 手动清空或者届时删除。