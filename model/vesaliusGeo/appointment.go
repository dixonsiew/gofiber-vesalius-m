package vesaliusGeo

import "encoding/xml"

type ResultListAppointment struct {
    XMLName      xml.Name      `xml:"Result" json:"-"`
    Appointments []Appointment `xml:"Appointment"`
    Success      Success       `xml:"Success"`
    Error        Error         `xml:"Error"`
}

type ResultListSlot struct {
    XMLName xml.Name `xml:"Result" json:"-"`
    Slots   []Slot   `xml:"Slot"`
    Success Success  `xml:"Success"`
    Error   Error    `xml:"Error"`
}

type ResultAppointmentBookingConfirmation struct {
    XMLName                        xml.Name                       `xml:"Result" json:"-"`
    AppointmentBookingConfirmation AppointmentBookingConfirmation `xml:"AppointmentBookingConfirmation"`
    Success                        Success                        `xml:"Success"`
    Error                          Error                          `xml:"Error"`
}

type ResultAppointmentChangeConfirmation struct {
    XMLName                       xml.Name                      `xml:"Result" json:"-"`
    AppointmentChangeConfirmation AppointmentChangeConfirmation `xml:"AppointmentChangeConfirmation"`
    Success                       Success                       `xml:"Success"`
    Error                         Error                         `xml:"Error"`
}

type ResultAppointmentCancelConfirmation struct {
    XMLName                       xml.Name                      `xml:"Result" json:"-"`
    AppointmentCancelConfirmation AppointmentCancelConfirmation `xml:"AppointmentCancelConfirmation"`
    Success                       Success                       `xml:"Success"`
    Error                         Error                         `xml:"Error"`
}

type Appointment struct {
    XMLName           xml.Name `xml:"Appointment" json:"-"`
    AppointmentNumber string   `xml:"AppointmentNumber" json:"appointmentNumber"`
    Date              string   `xml:"Date" json:"date"`
    Day               string   `xml:"Day" json:"day"`
    StartTime         string   `xml:"StartTime" json:"startTime"`
    EndTime           string   `xml:"EndTime" json:"endTime"`
    DoctorMCR         string   `xml:"DoctorMCR" json:"doctorMcr"`
    DoctorName        string   `xml:"DoctorName" json:"doctorName"`
    SpecialtyCode     string   `xml:"SpecialtyCode" json:"specialtyCode"`
    Specialty         string   `xml:"Specialty" json:"specialty"`
    Clinic            string   `xml:"Clinic" json:"clinic"`
    Room              string   `xml:"Room" json:"room"`
    CaseType          string   `xml:"CaseType" json:"caseType"`
}

type Slot struct {
    XMLName     xml.Name `xml:"Slot" json:"-"`
    SlotNumber  string   `xml:"SlotNumber" json:"slotNumber"`
    Date        string   `xml:"Date" json:"date"`
    Day         string   `xml:"Day" json:"day"`
    StartTime   string   `xml:"StartTime" json:"startTime"`
    EndTime     string   `xml:"EndTime" json:"endTime"`
    DoctorName  string   `xml:"DoctorName" json:"doctorName"`
    Speciality  string   `xml:"Specialty" json:"specialty"`
    Clinic      string   `xml:"Clinic" json:"clinic"`
    Room        string   `xml:"Room" json:"room"`
    CaseType    string   `xml:"CaseType" json:"caseType"`
    SessionType string   `json:"sessionType"`
}

type AppointmentBookingConfirmation struct {
    XMLName           xml.Name `xml:"AppointmentBookingConfirmation" json:"-"`
    AppointmentStatus string   `xml:"AppointmentStatus" json:"appointmentStatus"`
    AppointmentNumber string   `xml:"AppointmentNumber" json:"appointmentNumber"`
    Date              string   `xml:"Date" json:"date"`
    Day               string   `xml:"Day" json:"day"`
    StartTime         string   `xml:"StartTime" json:"startTime"`
    EndTime           string   `xml:"EndTime" json:"endTime"`
    DoctorName        string   `xml:"DoctorName" json:"doctorName"`
    Specialty         string   `xml:"Specialty" json:"specialty"`
    Clinic            string   `xml:"Clinic" json:"clinic"`
    Room              string   `xml:"Room" json:"room"`
    CaseType          string   `xml:"CaseType" json:"caseType"`
}

type AppointmentChangeConfirmation struct {
    XMLName           xml.Name `xml:"AppointmentChangeConfirmation" json:"-"`
    AppointmentStatus string   `xml:"AppointmentStatus" json:"appointmentStatus"`
    AppointmentNumber string   `xml:"AppointmentNumber" json:"appointmentNumber"`
    Date              string   `xml:"Date" json:"date"`
    Day               string   `xml:"Day" json:"day"`
    StartTime         string   `xml:"StartTime" json:"startTime"`
    EndTime           string   `xml:"EndTime" json:"endTime"`
    DoctorName        string   `xml:"DoctorName" json:"doctorName"`
    Specialty         string   `xml:"Specialty" json:"specialty"`
    Clinic            string   `xml:"Clinic" json:"clinic"`
    Room              string   `xml:"Room" json:"room"`
    CaseType          string   `xml:"CaseType" json:"caseType"`
}

type AppointmentCancelConfirmation struct {
    XMLName           xml.Name `xml:"AppointmentCancelConfirmation" json:"-"`
    AppointmentStatus string   `xml:"AppointmentStatus" json:"appointmentStatus"`
    AppointmentNumber string   `xml:"AppointmentNumber" json:"appointmentNumber"`
    Date              string   `xml:"Date" json:"date"`
    Day               string   `xml:"Day" json:"day"`
    StartTime         string   `xml:"StartTime" json:"startTime"`
    EndTime           string   `xml:"EndTime" json:"endTime"`
    DoctorName        string   `xml:"DoctorName" json:"doctorName"`
    Specialty         string   `xml:"Specialty" json:"specialty"`
    Clinic            string   `xml:"Clinic" json:"clinic"`
    Room              string   `xml:"Room" json:"room"`
    CaseType          string   `xml:"CaseType" json:"caseType"`
}
