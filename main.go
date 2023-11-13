package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/handlers"
	"github.com/shgxzybaba/go_web01/security"
	"github.com/shgxzybaba/go_web01/utils"
)

func indexHandler(c *fiber.Ctx) (err error) {
	log.Info(c.Locals("Headers"))

	return c.Render("index", utils.DefaultResponse(c), "layout")
}

func main() {

	log.Info("Hello server!")

	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layout",
		PassLocalsToViews: true,
	})

	app.Use(func(c *fiber.Ctx) error {
		sess, err := security.GetSessionData(c)
		if err != nil {
			return c.Next()
		}

		c.Locals("sessionData", sess)
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {

		if c.Path() == "/static" || c.Path() == "/logout" {
			return c.Next()
		}

		if c.Locals("sessionData") != nil {
			sessionData := (c.Locals("sessionData")).(data.SessionData)
			c.Locals("Headers", utils.LoggedInHeaders)
			c.Locals("LoggedIn", true)
			c.Locals("Username", sessionData.Email)
		} else {
			c.Locals("LoggedIn", false)
			c.Locals("Headers", utils.HeaderLinks)
		}

		return c.Next()
	})

	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	app.Static("/static", "./static")

	app.Get("/", indexHandler)
	app.Get("/login", security.GetLoginPage)
	app.Post("/login", security.LoginHandler)
	app.Get("/logout", security.LogoutHandler)
	app.Get("/dashboard", handlers.DashboardHandler)
	app.Get("/tutor", handlers.GetTutor)
	app.Get("/notes", handlers.AllCourseNotes)
	app.Get("/add_note", handlers.AddNote)
	app.Post("/note", handlers.CreateNote)
	app.Put("/note", handlers.UpdateNote)
	app.Get("/edit_note", handlers.EditNote)

	e := app.Listen(":8085")
	if e != nil {
		log.Error("An error occurred while starting the server", e)
	}
}
