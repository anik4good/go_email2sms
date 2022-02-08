package routes

import (
	Controller "github.com/anik4good/go_email2sms/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	//Controller "github.com/anik4good/go_email/api"
)

func SetUpRoutes(app *fiber.App) {

	app.Get("/dashboard", monitor.New())
	api := app.Group("/api")
	//apiv1 := api.Group("/v1")

	users := api.Group("/users")

	users.Get("/hello", Controller.Hello)
	users.Post("/create", Controller.AddUser)
	users.Post("/create-random", Controller.AddUserRandom)
	users.Get("/get/:id", Controller.GetSingleUser)
	users.Get("/getall", Controller.GetAllUsers)
	// app.Get("/book/:id", Controller.GetBook)
	// app.Post("/book", Controller.AddBook)
	// app.Put("/book", Controller.Update)
	// app.Delete("/book", Controller.Delete)

	smsc := api.Group("/smsc")

	smsc.Get("/sendmsg", Controller.SmsApi)
}

// func RegisterAPI(api fiber.Router) {
// 	// registerRoles(api, db)
// 	registerUsers(api)
// }

// func registerUsers(api fiber.Router) {
// 	users := api.Group("/users")

// 	users.Get("/", Controller.hello)
// 	// users.Get("/:id", Controller.GetUser(db))
// 	// users.Post("/", Controller.AddUser(db))
// 	// users.Put("/:id", Controller.EditUser(db))
// 	// users.Delete("/:id", Controller.DeleteUser(db))
// }
