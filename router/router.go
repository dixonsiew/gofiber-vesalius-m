package router

import (
    "vesaliusm/controller/admin"
    "vesaliusm/controller/auth"
    "vesaliusm/controller/clubs"
    "vesaliusm/controller/common"
    "vesaliusm/controller/feedback"
    "vesaliusm/controller/futureOrder"
    "vesaliusm/controller/guest"
    "vesaliusm/controller/hpackage"
    "vesaliusm/controller/ipay"
    "vesaliusm/controller/logistic"
    "vesaliusm/controller/maintenance"
    "vesaliusm/controller/myFamily"
    "vesaliusm/controller/publicBranch"
    "vesaliusm/controller/publicVesalius"
    "vesaliusm/controller/qms"
    "vesaliusm/controller/user"
    "vesaliusm/controller/userBilling"
    "vesaliusm/controller/userNotification"
    "vesaliusm/controller/userPackage"
    "vesaliusm/controller/vesalius"
    "vesaliusm/controller/vesaliusGeo"
    "vesaliusm/controller/wallex"
    "vesaliusm/controller/wayFinding"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, basePath string) {
    api := app.Group(basePath)
    admin.SetupRoutes(api)
    auth.SetupRoutes(api)
    clubs.SetupRoutes(api)
    common.SetupRoutes(api)
    feedback.SetupRoutes(api)
    futureOrder.SetupRoutes(api)
    guest.SetupRoutes(api)
    hpackage.SetupRoutes(api)
    ipay.SetupRoutes(api)
    logistic.SetupRoutes(api)
    maintenance.SetupRoutes(api)
    myFamily.SetupRoutes(api)
    publicBranch.SetupRoutes(api)
    publicVesalius.SetupRoutes(api)
    qms.SetupRoutes(api)
    user.SetupRoutes(api)
    userBilling.SetupRoutes(api)
    userNotification.SetupRoutes(api)
    userPackage.SetupRoutes(api)
    vesalius.SetupRoutes(api)
    vesaliusGeo.SetupRoutes(api)
    wallex.SetupRoutes(api)
    wayFinding.SetupRoutes(api)
}
