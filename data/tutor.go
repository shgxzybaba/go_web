package data

type Tutor struct {
	FullName          string
	YearsOfExperience int
}

func GetTutorByCourse(course string, student int) (tutor Tutor, err error) {
	query := `SELECT t.fullname, t.experience FROM tutor
    join courses_tutors ct on tutor.id = ct.tutor_id
    join tutor t on t.id = ct.tutor_id
    join courses c on c.id = ct.course_id
    join student_courses sc on c.id = sc.course_id
    where c.title = $1 and sc.student_id = $2` //todo use course id to shorten joins
	row := DB.QueryRow(query, course, student)
	err = row.Scan(&tutor.FullName, &tutor.YearsOfExperience)
	return
}
