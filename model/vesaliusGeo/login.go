package vesaliusGeo

import "encoding/xml"

type ResultToken struct {
    XMLName xml.Name `xml:"Result" json:"-"`
    Token   Token    `xml:"Token"`
    Error   Error    `xml:"Error"`
}

type ResultLogout struct {
    XMLName xml.Name `xml:"Result" json:"-"`
    Success Success  `xml:"Success"`
    Error   Error    `xml:"Error"`
}

type Token struct {
    XMLName     xml.Name `xml:"Token" json:"-"`
    TokenNumber string   `xml:"Token_number" json:"tokenNumber"`
}
