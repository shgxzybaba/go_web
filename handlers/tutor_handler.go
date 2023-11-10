package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/security"
	"github.com/shgxzybaba/go_web01/utils"
)

func GetTutor(c *fiber.Ctx) error {
	sessionData, err := security.GetSessionData(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("could not get session from request")
	}

	courseName := c.Query("course", "")
	tutor, err := data.GetTutorByCourse(courseName, sessionData.UserId)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Failed to add courses to the student")
	}
	//return c.SendString(fmt.Sprintf("Fullname: %s\nExperience: %d", tutor.FullName, tutor.YearsOfExperience))
	html, err := utils.GenerateHTML(tutor, "tutor")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Failed to add courses to the student")
	}
	return c.SendString(html)
}
