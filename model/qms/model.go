package qms

type QmsPatient struct {
    Prn          string `json:"prn"`
    Relationship string `json:"relationship"`
}

type QueueResult struct {
    QueueNumber        string `json:"queueNumber"`
    PatientsAheadOfYou string `json:"patientsAheadOfYou"`
    TicketStatus       string `json:"ticketStatus"`
    PatientPrn         string `json:"patientPrn"`
    Relationship       string `json:"relationship"`
    PatientName        string `json:"patientName"`
    DoctorName         string `json:"doctorName"`
    RoomNumber         string `json:"roomNumber"`
    AsAt               string `json:"asAt"`
}

type WsAuth struct {
    WsClientName string `json:"WsClientName"`
    Password     string `json:"Password"`
}

type DateRange struct {
    StartDate string `json:"StartDate"`
    EndDate   string `json:"EndDate"`
}

type Cust struct {
    MrnNumber string `json:"MrnNumber"`
}

type GetTicketInfoByDateListReq struct {
    WsAuth    WsAuth    `json:"WsAuth"`
    DateRange DateRange `json:"DateRange"`
    Cust      Cust      `json:"Cust"`
}

type RequestPayload struct {
    GetTicketInfoByDateListReq GetTicketInfoByDateListReq `json:"GetTicketInfoByDateListReq"`
}

// Response structures
type TicketInfo struct {
    TicketStr         string `json:"TicketStr"`
    TotalWaitingAhead string `json:"TotalWaitingAhead"`
    TicketStatus      string `json:"TicketStatus"`
    Cust              struct {
        Name string `json:"Name"`
    } `json:"Cust"`
    IndService struct {
        ServiceName string `json:"ServiceName"`
        ShortName   string `json:"ShortName"`
    } `json:"IndService"`
}

type TicketInfoList struct {
    TicketInfo any `json:"TicketInfo"` // Can be single object or array
}

type TicketInfoByDate struct {
    TicketInfoList TicketInfoList `json:"TicketInfoList"`
}

type TicketInfoByDateList struct {
    TicketInfoByDate any `json:"TicketInfoByDate"`
}

type GetTicketInfoByDateListResp struct {
    TicketInfoByDateList TicketInfoByDateList `json:"TicketInfoByDateList"`
    TotalRec             string               `json:"TotalRec"`
}

type ResponseData struct {
    GetTicketInfoByDateListResp GetTicketInfoByDateListResp `json:"GetTicketInfoByDateListResp"`
}
