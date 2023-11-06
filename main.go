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

	response := utils.Data{}
	students, err := data.FetchAllStudents()
	if err != nil {
		response.Err = err.Error()
		return c.Render("error", response)
	}
	response.Response, response.Err = students, ""

	return c.Render("index", response, "layout")
}

func main() {

	log.Info("Hello server!")

	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layout",
		PassLocalsToViews: true,
	})

	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	app.Static("/static", "./static")

	app.Get("/", indexHandler)
	app.Get("/login", security.GetLoginPage)
	app.Post("/login", security.LoginHandler)
	app.Get("/dashboard", handlers.DashboardHandler)
	app.Get("/tutor", handlers.GetTutor)

	e := app.Listen(":8085")
	if e != nil {
		log.Error("An error occurred while starting the server", e)
	}
}
