package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/utils"
)

type NotesDto struct {
	Notes []data.CourseNote
	Title string
}

func DashboardHandler(c *fiber.Ctx) error {

	sess := c.Locals("sessionData")
	if sess == nil {
		return c.Redirect("/")
	}

	sessionData := (sess).(data.SessionData)

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

func AllCourseNotes(c *fiber.Ctx) error {

	sess := c.Locals("sessionData")
	if sess == nil {
		return c.Redirect("/")
	}

	sessionData := (sess).(data.SessionData)
	course := c.Query("course", "")
	notes, err := data.GetStudentCourseNotes(course, sessionData.UserId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("unable to get user sessions")
	}
	noteDto := NotesDto{notes, course} //todo: should validate/escape the course in the request
	html, err := utils.GenerateHTML(noteDto, "notes")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("unable to get user sessions")
	}
	return c.SendString(html)

}

func AddNote(c *fiber.Ctx) error {

	sess := c.Locals("sessionData")
	if sess == nil {
		return c.Redirect("/")
	}

	course := c.Query("course", "") //todo: validate param
	if course == "" {
		return c.Status(fiber.StatusNotAcceptable).SendString("no course was selected")
	}
	html, err := utils.GenerateHTML(course, "add-note")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("unable to get user sessions")
	}
	return c.SendString(html)
}

func CreateNote(c *fiber.Ctx) error {
	sess := c.Locals("sessionData")
	if sess == nil {
		return c.Redirect("/")
	}

	sessionData := (sess).(data.SessionData)
	course := c.Query("course", "")
	note := c.FormValue("create-note") //todo: validate this
	if note == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Note text cannot be empty")
	}
	_, err := data.SaveStudentNote(course, sessionData.UserId, note)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("An error occurred while saving note")
	}

	return AllCourseNotes(c)
}

func EditNote(c *fiber.Ctx) error {
	sess := c.Locals("sessionData")

	if sess == nil {
		return c.Redirect("/")
	}
	sessionData := (sess).(data.SessionData)
	noteId := c.QueryInt("note_id", 0)
	if noteId == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid note id")
	}

	note, err := data.GetStudentNote(sessionData.UserId, noteId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid note id")
	}
	html, err := utils.GenerateHTML(note, "edit-note")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("unable to get user sessions")
	}
	return c.SendString(html)

}
