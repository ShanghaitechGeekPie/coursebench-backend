import json

# 文件路径
TEACHERS_JSON = 'teachers.json'
DATA_JSON = 'data.json'

# 读取 data.json，建立 workcode->name 映射
def load_workcode_name_map():
    with open(DATA_JSON, 'r', encoding='utf-8') as f:
        data = json.load(f)
    # 假设 data.json 是一个列表，每个元素有 workcode 和 name 字段
    workcode_name = {}
    for item in data:
        code = str(item.get('workcode'))
        name = item.get('lastname')
        if code and name:
            workcode_name[code] = name
    return workcode_name

# 读取 teachers.json，替换姓名
def correct_teacher_names():
    workcode_name = load_workcode_name_map()
    with open(TEACHERS_JSON, 'r', encoding='utf-8') as f:
        teachers = json.load(f)
    changed = False
    for course_code, teacher_map in teachers.items():
        for workcode in list(teacher_map.keys()):
            if workcode in workcode_name:
                old_name = teacher_map[workcode]
                new_name = workcode_name[workcode]
                if old_name != new_name:
                    teacher_map[workcode] = new_name
                    changed = True
    if changed:
        with open(TEACHERS_JSON, 'w', encoding='utf-8') as f:
            json.dump(teachers, f, ensure_ascii=False, indent=4)
        print('教师姓名已根据 data.json 修正并保存到 teachers.json')
    else:
        print('没有需要修正的教师姓名')

if __name__ == '__main__':
    correct_teacher_names()
