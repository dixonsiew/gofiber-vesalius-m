package utils

import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/go-playground/validator/v10"
    "github.com/go-resty/resty/v2"
    "github.com/gofiber/fiber/v2"
    "github.com/rs/zerolog"
    "github.com/ztrue/tracerr"
)

type (
    ErrorResponse struct {
        Error       bool
        FailedField string
        Tag         string
        Value       any
        Param       string
    }

    XValidator struct {
        validator *validator.Validate
    }
)

var (
    validate = validator.New()
    Logger zerolog.Logger
    iLogger zerolog.Logger
    appValidator *XValidator
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

func SetValidator() {
    v := &XValidator{
        validator: validate,
    }
    appValidator = v
}

func GetValidator() *XValidator {
    return appValidator
}

func ValidatePayload(data any, c *fiber.Ctx) error {
    errs := GetValidator().Validate(data)
    if len(errs) > 0 && errs[0].Error {
        errMsgs := make([]string, 0)
        for _, err := range errs {
            switch err.Tag {
            case "required":
                ex := fmt.Sprintf("[%s] is %s", err.FailedField, err.Tag)
                errMsgs = append(errMsgs, ex)
            case "max":
                ex := fmt.Sprintf("[%s] max length is %s", err.FailedField, err.Param)
                errMsgs = append(errMsgs, ex)
            case "min":
                ex := fmt.Sprintf("[%s] min length is %s", err.FailedField, err.Param)
                errMsgs = append(errMsgs, ex)
            default:
                errMsgs = append(errMsgs, fmt.Sprintf(
                    "[%s]: '%v' | Needs to implement '%s' '%s'",
                    err.FailedField,
                    err.Value,
                    err.Tag,
                    err.Param,
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

func (v XValidator) Validate(data any) []ErrorResponse {
    validationErrors := []ErrorResponse{}

    errs := validate.Struct(data)
    if errs != nil {
        for _, err := range errs.(validator.ValidationErrors) {
            // In this case data object is actually holding the User struct
            var elem ErrorResponse

            elem.FailedField = err.Field() // Export struct field name
            elem.Tag = err.Tag()           // Export struct tag
            elem.Value = err.Value()       // Export field value
            elem.Error = true
            elem.Param = err.Param()

            validationErrors = append(validationErrors, elem)
        }
    }

    return validationErrors
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
