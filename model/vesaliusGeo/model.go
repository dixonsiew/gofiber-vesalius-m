package vesaliusGeo

import "encoding/xml"

type XmlResponse struct {
    Return string `xml:"return"`
}

type VesaliusWSException struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type Error struct {
    XMLName      xml.Name `xml:"Error" json:"-"`
    ErrorCode    string   `xml:"ErrorCode" json:"errorCode"`
    ErrorMessage string   `xml:"ErrorMessage" json:"errorMessage"`
}

type Success struct {
    XMLName xml.Name `xml:"Success" json:"-"`
    Code    string   `xml:"Code" json:"code"`
    Message string   `xml:"Message" json:"message"`
}

type ResultList struct {
    XMLName xml.Name `xml:"Result" json:"-"`
    Success Success  `xml:"Success"`
    Error   Error    `xml:"Error"`
}
