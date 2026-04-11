package mail

type MailSignUpSuccess struct {
    PatientName string `json:"patientName"`
    Username    string `json:"username"`
    Email       string `json:"email"`
}

type MailGuestAppointmentConfirmation struct {
    GuestName       string `json:"guestName"`
    DoctorName      string `json:"doctorName"`
    AppointmentDate string `json:"appointmentDate"`
    AppointmentTime string `json:"appointmentTime"`
    ClinicLocation  string `json:"clinicLocation"`
    Email           string `json:"email"`
}

type MailLittleKids struct {
    KidsName string `json:"kidsName"`
    Email    string `json:"email"`
}

type MailGoldenPearl struct {
    GoldenName string `json:"goldenName"`
    Email      string `json:"email"`
}

type MailClubsEventRegistrationToMember struct {
    ActivityName string `json:"activityName"`
    MemberName   string `json:"memberName"`
    Email        string `json:"email"`
}

type MailSuccessOutstandingBillPayment struct {
    Amount        string `json:"amount"`
    PaymentMethod string `json:"paymentMethod"`
    BillNumber    string `json:"billNumber"`
    InvoiceNumber string `json:"invoiceNumber"`
    Email         string `json:"email"`
}

type MailSuccessPackagePayment struct {
    PatientName       string `json:"patientName"`
    OrderNumber       string `json:"orderNumber"`
    DateOfPurchase    string `json:"dateOfPurchase"`
    ProductName       string `json:"productName"`
    ProductQuantity   string `json:"productQuantity"`
    ProductPrice      string `json:"productPrice"`
    SubtotalPrice     string `json:"subtotalPrice"`
    PaymentMethod     string `json:"paymentMethod"`
    TotalPrice        string `json:"totalPrice"`
    PackageExpiryDate string `json:"packageExpiryDate"`
    BillingAddress    string `json:"billingAddress"`
    Email             string `json:"email"`
}
