package router

import (
    "vesaliusm/router/common"
    "vesaliusm/router/auth"

    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, basePath string) {
    api := app.Group(basePath)
    common.SetupRoutes(api)
    auth.SetupRoutes(api)
}
