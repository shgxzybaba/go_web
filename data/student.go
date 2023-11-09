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

type CourseNote struct {
	Id          int
	Text        string
	CourseTitle string
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

func GetStudentCourseNotes(course string, studentId int) (notes []CourseNote, err error) {
	query := `SELECT n.text,n.id, c.title from notes n
            join courses c on c.id = n.course_id
            where n.student_id = $1
            and c.title = $2`
	rows, err := DB.Query(query, studentId, course)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		note := CourseNote{}
		if err := rows.Scan(&note.Text, &note.Id, &note.CourseTitle); err != nil {
			return nil, err // Return nil slice and error
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err // Return nil slice and error if there was an error during iteration
	}

	return
}

func SaveStudentNote(course string, studentId int, note string) (string, error) {
	query := `insert into notes(text, student_id, course_id)
select $1, $2, c.id from courses c where c.title = $3`
	_, err := DB.Exec(query, note, studentId, course)
	return note, err
}

func GetStudentNote(studentId int, noteId int) (note CourseNote, err error) {
	query := `select n.id, n.text from notes n
where n.student_id = $1 and n.id = $2
`
	row := DB.QueryRow(query, studentId, noteId)
	courseNote := CourseNote{}
	err = row.Scan(&courseNote.Id, &courseNote.Text)
	return
}
