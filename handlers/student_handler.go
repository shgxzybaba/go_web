package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/security"
)

func DashboardHandler(c *fiber.Ctx) error {

	sessionData, err := security.GetSessionData(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("unable to get user sessions")
	}

	// Get student from session
	student, err := data.GetStudentFromId(sessionData.UserId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Failed to get student from session")
	}

	// Add courses to the student
	err = student.AddCourses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Failed to add courses to the student")
	}

	// Render the dashboard template with the student data
	return c.Render("dashboard", fiber.Map{
		"Student":          student,
		"coursesForTheDay": student.TodaysCourses(),
	})
}
