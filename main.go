package main

import (
    "errors"
    "fmt"
    "html/template"
    "os"
    "strings"
    "vesaliusm/config"
    "vesaliusm/database"
    _ "vesaliusm/docs"
    "vesaliusm/router"
    "vesaliusm/utils"

    "github.com/go-playground/validator/v10"
    // jwtware "github.com/gofiber/contrib/jwt"
    swaggo "github.com/gofiber/contrib/v3/swaggo"
    fiberzerolog "github.com/gofiber/contrib/v3/zerolog"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/compress"
    "github.com/gofiber/fiber/v3/middleware/cors"
    "github.com/gofiber/fiber/v3/middleware/healthcheck"
    "github.com/gofiber/fiber/v3/middleware/recover"
    "github.com/gofiber/fiber/v3/middleware/static"
    "github.com/gofiber/template/django/v3"
)

// convert this to go using go-ora with github.com/jmoiron/sqlx

// @title (IHP-UAT) Vesalius-m Backend API
// @version 1.0
// @description Vesalius-m Backend API Docs (IHP-UAT Version).
// @BasePath /ih/mobile_central_2_0_0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
    defer utils.CatchPanic("main")
    runLogFile, _ := os.OpenFile(
        "app.log",
        os.O_APPEND|os.O_CREATE|os.O_WRONLY,
        0664,
    )
    defer runLogFile.Close()
    utils.SetClient()
    utils.SetLogger(runLogFile)
    port := config.Config("port")
    engine := django.New("./views", ".django")
    app := fiber.New(fiber.Config{
        StructValidator: &utils.StructValidator{Xvalidate: validator.New()},
        ErrorHandler: func(c fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            var e *fiber.Error
            if errors.As(err, &e) {
                code = e.Code
            }

            ms := err.Error()
            if strings.Contains(ms,"no rows in result set") {
                ms = "No record found"
                code = fiber.StatusNotFound
            }

            return c.Status(code).JSON(fiber.Map{
                "statusCode": code,
                "message":    ms,
            })
        },
        Views: engine,
    })
    app.Use(recover.New())
    app.Use(compress.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins:  []string{"*"},
        ExposeHeaders: []string{"Authorization", "filename", utils.X_TOTAL_COUNT, utils.X_TOTAL_PAGE},
    }))
    app.Use(fiberzerolog.New(fiberzerolog.Config{
        Logger: &utils.Logger,
    }))

    // app.Use(jwtware.New(jwtware.Config{
    //     SigningKey: jwtware.SigningKey{Key: []byte(utils.JWT_SECRET)},
    // }))
    database.ConnectDB()
    database.ConnectDBRs()
    defer database.CloseDB()
    defer database.CloseDBRs()

    basePath := "ih/mobile_central_2_0_0"
    initSwagger(app, basePath)
    app.Get(healthcheck.StartupEndpoint, healthcheck.New())
    app.Get(fmt.Sprintf("/%s/healthz", basePath), healthcheck.New())
    router.SetupRoutes(app, basePath)

    // if !fiber.IsChild() {
    //     cron.Setup()
    //     defer cron.Shutdown()
    // }

    err := app.Listen(fmt.Sprintf(":%s", port), fiber.ListenConfig{
        EnablePrefork: true,
    })

    if err != nil {
        utils.Logger.Fatal().Err(err).Msg("Fiber app error")
    }
}

func initSwagger(app *fiber.App, basePath string) {
    b, _ := os.ReadFile("./public/css/theme-flattop.css")
    css := string(b)

    cfg := swaggo.Config{
        URL:          "doc.json",
        DeepLinking:  true,
        DocExpansion: "list",
        Title:        "(IHP-UAT) Vesalius-m Backend API",
        SyntaxHighlight: &swaggo.SyntaxHighlightConfig{
            Activate: true,
            Theme:    "arta",
        },
        CustomStyle:          template.CSS(css),
        PersistAuthorization: true,
    }

    app.Get(fmt.Sprintf("/%s/docs/*", basePath), swaggo.New(cfg))
    app.Get(fmt.Sprintf("/%s/static*", basePath), static.New("./public"))
}
