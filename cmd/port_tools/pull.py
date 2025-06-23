# Step 1: 获取所有课程

import requests
import json
import math

def pull_courses():
    """
    Fetches all courses from the API and saves them to a JSON file.
    """
    base_url = "https://elrc.shanghaitech.edu.cn/learn/shanghai/tech/get/course"
    total_items = 9200
    page_size = 200
    total_pages = math.ceil(total_items / page_size)

    all_courses = {}

    for page in range(1, total_pages + 1):
        print(f"Fetching page {page}/{total_pages}...")
        params = {'page': page, 'size': page_size, 'courseType': 2}
        
        try:
            response = requests.get(base_url, params=params)
            response.raise_for_status()  # Raises an HTTPError for bad responses (4xx or 5xx)
        except requests.exceptions.RequestException as e:
            print(f"An error occurred while fetching page {page}: {e}")
            continue

        data = response.json()
        results = data.get("data", {}).get("results", [])

        if not results:
            print(f"No results found on page {page}.")
            continue

        for course in results:
            course_number = course.get("courseNumber")
            if not course_number:
                continue

            teachers = dict(zip(course.get("teacher", []), course.get("teacher_names", [])))

            # If the course number is not in the dictionary, add it with a new list.
            if course_number not in all_courses:
                all_courses[course_number] = []
            
            # Append the new course details to the list for that course number.
            all_courses[course_number].append({
                "code": course.get("courseNumber"),
                "name": course.get("name_"),
                "teacher": teachers,
                # also storing other fields to distinguish between courses
                "semester": course.get("semester_show_name"),
                "serialNumber": course.get("serialNumber"),
                "contentId": course.get("contentId_")
            })
    
    output_filename = 'courses.json'
    with open(output_filename, 'w', encoding='utf-8') as f:
        json.dump(all_courses, f, ensure_ascii=False, indent=4)

    print(f"Successfully fetched all courses and saved to {output_filename}")

if __name__ == "__main__":
    pull_courses()
