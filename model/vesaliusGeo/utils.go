package vesaliusGeo

import (
	"bytes"
    "strings"
)

type TrimmedString string

func (ts *TrimmedString) UnmarshalText(text []byte) error {
    *ts = TrimmedString(bytes.TrimSpace(text))
    return nil
}

func GetXmlResult(body string) []byte {
    i := strings.Index(body, "<ns:return>")
    j := strings.Index(body, "</ns:return>")
    content := body[11+i : j]
    content = strings.ReplaceAll(content, "&lt;", "<")
    content = strings.ReplaceAll(content, "&gt;", ">")
    content = strings.ReplaceAll(content, `encoding="UTF-8">`, `encoding="UTF-8"?>`)
    bx := []byte(content)
    return bx
}
