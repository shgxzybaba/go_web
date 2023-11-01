package data

import (
	"errors"
	"github.com/shgxzybaba/go_web01/utils"
	"net/http"
)

type Student struct {
	Username string
	Id       uint16
	Password string
	Courses  []Course
}

type Course struct {
	Title       string
	Description string
}

func (s *Student) addCourses() (err error) {
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
		s.Courses = append(s.Courses, course)
	}
	return

}

func (s *Student) CreateSession() (session Session, err error) {
	uuid := utils.GenerateUUID()
	session = Session{Uuid: uuid, UserId: int(s.Id)}
	_, err = DB.Exec("INSERT INTO sessions(user_id, uuid) values ($1, $2)", session.UserId, session.Uuid)
	return
}

func getStudentFromSession(uuid string) (student Student, err error) {
	student = Student{}
	row := DB.QueryRow("SELECT code, username from student join sessions s on student.code = s.user_id where s.uuid = $1", uuid)
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

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	var student Student
	var response = utils.Data{}

	cookie, err := r.Cookie("session-id")
	if err != nil {
		response.ErrorResponse(errors.New("session ID cookie not found"))
		utils.GenerateHTML(w, response, "layout", "navbar", "error") // Assuming you have an error page template
		return
	}

	student, err = getStudentFromSession(cookie.Value)
	if err != nil {
		response.ErrorResponse(errors.New("failed to get student from session"))
		utils.GenerateHTML(w, response, "layout", "navbar", "error") // Assuming you have an error page template
		return
	}

	err = student.addCourses()
	if err != nil {
		response.ErrorResponse(errors.New("failed to add courses to the student"))
		utils.GenerateHTML(w, response, "layout", "navbar", "error") // Assuming you have an error page template
		return
	}
	response.DataResponse(student)
	utils.GenerateHTML(w, response, "layout", "navbar", "dashboard")
}
