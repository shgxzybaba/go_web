package data

import (
	"fmt"
	"net/http"
	"time"
)

type Student struct {
	Username string
	Id       uint16
	Password string
	Courses  []Course
}

type Course struct {
	Id          int
	Title       string
	Description string
	Today       bool
	Days        []string
}

func (c *Course) OfferedDays() (err error) {
	fmt.Println("Adding offered days to courses")

	query := `select d.name from course_days cd 
join days d on d.id = cd.day_id
              join courses c on c.id = cd.course_id
where c.title = $1 
`
	rows, err1 := DB.Query(query, c.Title)

	if err1 != nil {
		err = err1
		return
	}

	for rows.Next() {
		day := ""
		err = rows.Scan(&day)
		if err != nil {
			return
		}
		c.Days = append(c.Days, day)
	}
	for _, v := range c.Days {
		if v == time.Now().Weekday().String() {
			c.Today = true
			break
		}
	}

	return

}
func (s *Student) TodaysCourses() (courses []Course) {
	for _, c := range s.Courses {
		if c.Today {
			courses = append(courses, c)
		}
	}
	return
}

func (s *Student) AddCourses() (err error) {
	fmt.Println("Adding courses to student")

	query := `
select c.description, c.title from student_courses sc
join courses c on sc.course_id = c.id
where sc.student_id = $1;
`
	rows, err1 := DB.Query(query, s.Id)

	if err1 != nil {
		err = err1
		return
	}

	for rows.Next() {
		course := Course{}
		err = rows.Scan(&course.Description, &course.Title)
		if err != nil {
			return
		}
		_ = course.OfferedDays() //todo: handle this error properly
		s.Courses = append(s.Courses, course)
	}
	return

}

func GetStudentFromId(id int) (student Student, err error) {
	student = Student{}
	row := DB.QueryRow("SELECT code, username from student s where s.id = $1", id)
	err = row.Scan(&student.Id, &student.Username)
	return
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

func FetchAllStudents() (students []Student, err error) {
	rows, err := DB.Query("SELECT code, username FROM student;")
	if err != nil {
		return nil, err // Return nil slice and error
	}
	defer rows.Close() // Ensure rows are closed after the function exits

	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.Id, &student.Username); err != nil {
			return nil, err // Return nil slice and error
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err // Return nil slice and error if there was an error during iteration
	}

	return students, nil // Return the populated slice and no error
}
