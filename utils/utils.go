package utils

import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/go-playground/validator/v10"
    "github.com/go-resty/resty/v2"
    "github.com/gofiber/fiber/v3"
    "github.com/rs/zerolog"
    "github.com/ztrue/tracerr"
)

type StructValidator struct {
    Xvalidate *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
    return v.Xvalidate.Struct(out)
}

var (
    Logger zerolog.Logger
    iLogger zerolog.Logger
    client       *resty.Client
)

func SetClient() {
    client = resty.New()
    client.SetTimeout(time.Minute * 5)
}

func SetLogger(runLogFile *os.File) {
    multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
    Logger = zerolog.New(multi).Level(zerolog.ErrorLevel).With().Timestamp().Caller().Logger()

    iLogger = zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}

func GetValidationErrors(errs validator.ValidationErrors) error {
    if len(errs) > 0 {
        errMsgs := make([]string, 0)
        for _, err := range errs {
            switch err.Tag() {
                case "required":
                    ex := fmt.Sprintf("[%s] is %s", err.Field(), err.Tag())
                    errMsgs = append(errMsgs, ex)
                case "max":
                    ex := fmt.Sprintf("[%s] max length is %s", err.Field(), err.Param())
                    errMsgs = append(errMsgs, ex)
                case "min":
                    ex := fmt.Sprintf("[%s] min length is %s", err.Field(), err.Param())
                    errMsgs = append(errMsgs, ex)
                default:
                    errMsgs = append(errMsgs, fmt.Sprintf(
                        "[%s]: '%v' | Needs to implement '%s' '%s'",
                        err.Field(),
                        err.Value(),
                        err.Tag(),
                        err.Param(),
                    ))
            }
        }

        return &fiber.Error{
            Code:    fiber.ErrBadRequest.Code,
            Message: strings.Join(errMsgs, " and "),
        }
    }

    return nil
}

func GetErrors(errs []error) string {
    ls := []string{}
    for _, err := range errs {
        ls = append(ls, err.Error())
    }

    return strings.Join(ls, "|")
}

func CatchPanic(funcName string) {
    if err := recover(); err != nil {
        //LogInfo(fmt.Sprintf("recovered from panic -%s:%v", funcName, err))
        LogError(fmt.Errorf("recovered from panic -%s:%v", funcName, err))
    }
}

func LogError(err error) {
    ex := tracerr.Wrap(err)
    Logger.Err(err).Msg(tracerr.Sprint(ex))
}

func LogInfo(s string) {
    iLogger.Info().Msg(s)
}
