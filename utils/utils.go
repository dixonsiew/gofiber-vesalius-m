package utils

import (
    "fmt"
    "math/rand"
    "os"
    "reflect"
    "strings"
    "time"

    "github.com/shopspring/decimal"
    "github.com/go-playground/validator/v10"
    "github.com/go-resty/resty/v2"
    "github.com/gofiber/fiber/v3"
    "github.com/guregu/null/v6"
    "github.com/rs/zerolog"
    "github.com/ztrue/tracerr"
)

type Map map[string]any

type StructValidator struct {
    Xvalidate *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
    return v.Xvalidate.Struct(out)
}

var (
    Logger  zerolog.Logger
    iLogger zerolog.Logger
    client  *resty.Client
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

func GetR() *resty.Request {
    return client.R()
}

func GetRWs(action string) *resty.Request {
    return client.R().
        SetHeader("Content-Type", "text/xml; charset=utf-8").
        SetHeader("SOAPAction", fmt.Sprintf("urn:%s", action))
}

func GetDbColsWithReplace(s interface{}, prefix string, m map[string]string) string {
    r := ""
    // Get the type of the struct
    t := reflect.TypeOf(s)
    if t.Kind() == reflect.Ptr {
        t = t.Elem() // Dereference the pointer
    }

    // Check if the input is a struct (or pointer to a struct)
    if t.Kind() != reflect.Struct {
        return r
    }

    var columns []string
    for i := 0; i < t.NumField(); i++ {
        // Get the field information
        field := t.Field(i)
        
        // Access the "db" tag value using the Get method
        tagValue := field.Tag.Get("db")
        
        // If a "db" tag exists and is not empty, add it to the list
        if tagValue != "" {
            x := fmt.Sprintf("%s%s", prefix, tagValue)
            if val, ok := m[x]; ok {
                if val == "" {
                    continue
                } else {
                    x = val
                }
            }
            columns = append(columns, x)
        }
    }
    
    return strings.Join(columns, ", ")
}

func GetDbCols(s interface{}, prefix string) string {
    r := ""
    // Get the type of the struct
    t := reflect.TypeOf(s)
    if t.Kind() == reflect.Ptr {
        t = t.Elem() // Dereference the pointer
    }

    // Check if the input is a struct (or pointer to a struct)
    if t.Kind() != reflect.Struct {
        return r
    }

    var columns []string
    for i := 0; i < t.NumField(); i++ {
        // Get the field information
        field := t.Field(i)
        
        // Access the "db" tag value using the Get method
        tagValue := field.Tag.Get("db")
        
        // If a "db" tag exists and is not empty, add it to the list
        if tagValue != "" {
            columns = append(columns, fmt.Sprintf("%s%s", prefix, tagValue))
        }
    }
    
    return strings.Join(columns, ", ")
}

func NewNullString(s string) null.String {
    if s == "" {
        return null.NewString(s, false)
    }
    return null.NewString(s, true)
}

func NewInt32(i int32) null.Int32 {
    return null.NewInt32(i, true)
}

func NewInt64(i int64) null.Int64 {
    return null.NewInt(i, true)
}

func NewFloat(f float64) null.Float {
    return null.NewFloat(f, true)
}

func GetAmount(s any) string {
    r := fmt.Sprintf("%v", s)
    d, _ := decimal.NewFromString(r)
    return d.StringFixed(2)
}

func BindNValidate(c fiber.Ctx, out any) error {
    if err := c.Bind().Body(out); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := GetValidationErrors(validationErrors)
            if errs != nil {
                return errs
            }
        }

        return err
    }

    return nil
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

func (m Map) GetString(key string) string {
    s := ""
    val, ok := m[key]
    if !ok {
        return s
    }

    r, ok := val.(string)
    if !ok {
        return s
    }

    return r

}

func GetRandomStr(length int) string {
    const charset = "0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
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
