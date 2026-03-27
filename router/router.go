package router

import (
	"vesaliusm/controller/admin"
	"vesaliusm/controller/auth"
	"vesaliusm/controller/clubs"
	"vesaliusm/controller/common"
	"vesaliusm/controller/futureOrder"
	"vesaliusm/controller/myFamily"
	"vesaliusm/controller/publicBranch"
	"vesaliusm/controller/user"
	"vesaliusm/controller/userNotification"
	"vesaliusm/controller/userPackage"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, basePath string) {
	api := app.Group(basePath)
	common.SetupRoutes(api)
	auth.SetupRoutes(api)
	clubs.SetupRoutes(api)
	futureOrder.SetupRoutes(api)
	myFamily.SetupRoutes(api)
	publicBranch.SetupRoutes(api)
	user.SetupRoutes(api)
	admin.SetupRoutes(api)
	userNotification.SetupRoutes(api)
	userPackage.SetupRoutes(api)
}
