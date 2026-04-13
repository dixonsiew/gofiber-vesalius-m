package middleware

import (
    "encoding/json"
    "fmt"
    "reflect"
    "regexp"
    "strings"

    "github.com/gofiber/fiber/v3"
)

func TrimMiddleware(c fiber.Ctx) error {
    // For body data
    if len(c.Body()) > 0 {
        c.Request().SetBody(TrimBytes(c.Body()))
    }

    route := c.Route()

    // For params (URL parameters)
    for _, paramName := range route.Params {
        c.Params(paramName, TrimString(c.Params(paramName)))
    }

    return c.Next()
}

func TrimString(value string) string {
    return strings.TrimSpace(value)
}

func TrimBytes(data []byte) []byte {
    var obj any
    json.Unmarshal(data, &obj)
    trimmed := TrimValue(obj)
    result, _ := json.Marshal(trimmed)
    return result
}

func TrimValue(value any) any {
    switch v := value.(type) {
    case string:
        return strings.TrimSpace(v)
    case map[string]any:
        trimmedMap := make(map[string]any)
        for key, val := range v {
            trimmedMap[key] = TrimValue(val)
        }
        return trimmedMap
    case []any:
        trimmedArr := make([]any, len(v))
        for i, item := range v {
            trimmedArr[i] = TrimValue(item)
        }
        return trimmedArr
    default:
        return v
    }
}

func TrimStructFieldsRecursive(obj any) {
    val := reflect.ValueOf(obj)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    
    if val.Kind() != reflect.Struct {
        return
    }
    
    for _, field := range val.Fields() {        
        switch field.Kind() {
        case reflect.String:
            trimmed := strings.TrimSpace(field.String())
            field.SetString(trimmed)
        case reflect.Struct:
            // Recursively trim nested structs
            if field.CanAddr() {
                TrimStructFieldsRecursive(field.Addr().Interface())
            }
        case reflect.Ptr:
            if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
                TrimStructFieldsRecursive(field.Interface())
            }
        }
    }
}

func TrimCompletely(s any) string {
    r := fmt.Sprintf("%v", s)
    switch v := s.(type) {
    case string:
        r = v
        re := regexp.MustCompile(`\s+`)
        r = re.ReplaceAllString(v, "")
    }

    return r
}
