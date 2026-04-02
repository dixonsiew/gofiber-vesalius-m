package vesaliusGeo

import "encoding/xml"

// 970712-14-5329

type ResultPatient struct {
    XMLName xml.Name `xml:"Result" json:"-"`
    Patient Patient  `xml:"Patient"`
    Success Success  `xml:"Success"`
    Error   Error    `xml:"Error"`
}

type ResultListPatient struct {
    XMLName  xml.Name  `xml:"Result" json:"-"`
    Patients []Patient `xml:"Patient"`
    Success  Success   `xml:"Success"`
    Error    Error     `xml:"Error"`
}

type Name struct {
    XMLName    xml.Name `xml:"Name" json:"-"`
    Title      string   `xml:"Title" json:"title"`
    FirstName  string   `xml:"FirstName" json:"firstName"`
    MiddleName string   `xml:"MiddleName" json:"middleName"`
    LastName   string   `xml:"LastName" json:"lastName"`
}

type Sex struct {
    XMLName     xml.Name `xml:"Sex" json:"-"`
    Code        string   `xml:"Code" json:"code"`
    Description string   `xml:"Description" json:"description"`
}

type HomeAddress struct {
    XMLName    xml.Name `xml:"HomeAddress" json:"-"`
    Address1   string   `xml:"Address1" json:"address1"`
    Address2   string   `xml:"Address2" json:"address2"`
    Address3   string   `xml:"Address3" json:"address3"`
    CityState  string   `xml:"CityState" json:"cityState"`
    PostalCode string   `xml:"PostalCode" json:"postalCode"`
    Country    string   `xml:"Country" json:"country"`
}

type ContactNumber struct {
    XMLName xml.Name `xml:"ContactNumber" json:"-"`
    Home    string   `xml:"Home" json:"home"`
    Email   string   `xml:"Email" json:"email"`
}

type Document struct {
    XMLName     xml.Name `xml:"Document" json:"-"`
    Code        string   `xml:"Code" json:"code"`
    Description string   `xml:"Description" json:"description"`
    Value       string   `xml:"Value" json:"value"`
    ExpireDate  string   `xml:"ExpireDate" json:"expireDate"`
}

type Nationality struct {
    XMLName     xml.Name `xml:"Nationality" json:"-"`
    Code        string   `xml:"Code" json:"code"`
    Description string   `xml:"Description" json:"description"`
}

type ChargeCategory struct {
    XMLName     xml.Name `xml:"ChargeCategory" json:"-"`
    Code        string   `xml:"Code" json:"code"`
    Description string   `xml:"Description" json:"description"`
}

type PaymentClass struct {
    XMLName     xml.Name `xml:"PaymentClass" json:"-"`
    Code        string   `xml:"Code" json:"code"`
    Description string   `xml:"Description" json:"description"`
}

type Patient struct {
    XMLName        xml.Name       `xml:"Patient" json:"-"`
    Prn            string         `xml:"PRN" json:"prn"`
    Name           Name           `xml:"Name" json:"name"`
    Resident       string         `xml:"Resident" json:"resident"`
    DOB            string         `xml:"DOB" json:"dob"`
    Sex            Sex            `xml:"Sex" json:"sex"`
    HomeAddress    HomeAddress    `xml:"HomeAddress" json:"homeAddress"`
    ContactNumber  ContactNumber  `xml:"ContactNumber" json:"contactNumber"`
    Document       []Document     `xml:"Document" json:"documents"`
    Nationality    Nationality    `xml:"Nationality" json:"nationality"`
    ChargeCategory ChargeCategory `xml:"ChargeCategory" json:"chargeCategory"`
    PaymentClass   PaymentClass   `xml:"PaymentClass" json:"paymentClass"`
}
