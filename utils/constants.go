package utils

const (
    PAGE_SIZE        = 10
    MAX_PAGE_NUMBERS = 10
    X_TOTAL_COUNT    = "x-total-count"
    X_TOTAL_PAGE     = "x-total-page"
    JWT_SECRET       = "itsSecret"
    ROLE_ADMIN       = "ADMIN"
    ROLE_SUPER_ADMIN = "SUPER ADMIN"
    ROLE_USER        = "USER"
)

const (
    MsgTypePromotion         = "PROMOTION"
    MsgTypeQueueNotification = "QUEUE_NOTIFICATION"
    MsgTypeUpcomingAppt      = "UPCOMING_APPT"
    MsgTypePatientFeedback   = "PATIENT_FEEDBACK"
    MsgTypeGeneralInfo       = "GENERAL_INFO"
)

const (
    QmsMethodNew  = "New"
    QmsMethodNear = "Near"
    QmsMethodCall = "Call"
)

const (
    NotificationTitleQueueNotification = "Attention: Queue Update for Your Service"
    NotificationTitleUpcomingAppt      = "Friendly Reminder: Your Appointment is Approaching"
    NotificationTitlePatientFeedback   = "Your Voice Matters: Help Us Enhance Your Experience"
)

const (
    PaymentStatusSubmitted = "Submitted"
    PaymentStatusPaid      = "Paid"
)

const (
    PackageStatusOrdered   = "ORDERED"
    PackageStatusPurchased = "Purchased"
    PackageStatusBooked    = "Booked"
    PackageStatusRedeemed  = "Redeemed"
    PackageStatusCancelled = "Cancelled"
)

const (
    LogisticRequestStatusConfirmed = "Confirmed"
    LogisticRequestStatusCancelled = "Cancelled"
    LogisticRequestStatusRejected  = "Rejected"
)

const (
    PaymentMethodWallex = 1
    PaymentMethodIpay88 = 2
)

const (
    ClubsDocTypeNRIC      = "NRIC"
    ClubsDocTypeBirthCert = "Birth Cert"
    ClubsDocTypePassport  = "Passport"
)
