package router

import (
    "vesaliusm/router/admin"
    "vesaliusm/router/auth"
    "vesaliusm/router/common"
    "vesaliusm/router/clubs"
    "vesaliusm/router/futureOrder"
    "vesaliusm/router/publicBranch"
    "vesaliusm/router/user"
    "vesaliusm/router/userNotification"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, basePath string) {
    api := app.Group(basePath)
    common.SetupRoutes(api)
    auth.SetupRoutes(api)
    clubs.SetupRoutes(api)
    futureOrder.SetupRoutes(api)
    publicBranch.SetupRoutes(api)
    user.SetupRoutes(api)
    admin.SetupRoutes(api)
    userNotification.SetupRoutes(api)
}
