package data

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Tutor struct {
	FullName          string
	YearsOfExperience int
}

func GetTutor(c *fiber.Ctx) error {
	courseName := c.Query("course", "")
	tutor, err := getTutorByCourse(courseName)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Failed to add courses to the student")
	}
	return c.SendString(fmt.Sprintf("Fullname: %s\nExperience: %d", tutor.FullName, tutor.YearsOfExperience))
}

func getTutorByCourse(course string) (tutor Tutor, err error) {
	query := `SELECT t.fullname, t.experience FROM tutor
    join courses_tutors ct on tutor.id = ct.tutor_id
    join tutor t on t.id = ct.tutor_id
    join courses c on c.id = ct.course_id
    where c.title = $1` //todo use course id to shorten joins
	row := DB.QueryRow(query, course)
	err = row.Scan(&tutor.FullName, &tutor.YearsOfExperience)
	return
}
