package router

import (
    "vesaliusm/router/common"
    "vesaliusm/router/auth"
    "vesaliusm/router/futureOrder"
    "vesaliusm/router/user"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, basePath string) {
    api := app.Group(basePath)
    common.SetupRoutes(api)
    auth.SetupRoutes(api)
    futureOrder.SetupRoutes(api)
    user.SetupRoutes(api)
}
