# Step 3: 获取统一关系

import json

def get_teachers():
    """
    Processes courses.json to create a dictionary of teachers for each course.
    The output is a dictionary where keys are course numbers and values are
    a dictionary of teacher IDs to teacher names for that course.
    """
    input_filename = 'courses.json'
    output_filename = 'teachers.json'
    
    try:
        with open(input_filename, 'r', encoding='utf-8') as f:
            all_courses = json.load(f)
    except FileNotFoundError:
        print(f"Error: {input_filename} not found. Please run pull.py first.")
        return
    except json.JSONDecodeError:
        print(f"Error: Could not decode JSON from {input_filename}.")
        return

    all_teachers = {}

    for course_number, courses_list in all_courses.items():
        course_teachers = {}
        for course in courses_list:
            teacher_info = course.get("teacher", {})
            if teacher_info:
                course_teachers.update(teacher_info)
        
        if course_teachers:
            all_teachers[course_number] = course_teachers

    with open(output_filename, 'w', encoding='utf-8') as f:
        json.dump(all_teachers, f, ensure_ascii=False, indent=4)

    print(f"Successfully processed teachers and saved to {output_filename}")

if __name__ == "__main__":
    get_teachers() 