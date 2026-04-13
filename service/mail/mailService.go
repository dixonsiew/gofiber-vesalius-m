package mail

import (
    "crypto/tls"
    "fmt"
    "path/filepath"
    "strconv"
    "strings"
    "vesaliusm/config"
    "vesaliusm/model"
    lgm "vesaliusm/model/logistic"
    mm "vesaliusm/model/mail"
    "vesaliusm/service/emailMaster"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    //"github.com/go-gomail/gomail"
    "github.com/nleeper/goment"
    gomail "gopkg.in/mail.v2"
)

var MailSvc *MailService = NewMailService()

type LogoConfig struct {
    File   string
    Width  string
    Height string
}

type MailConfig struct {
    host       string
    port       int
    user       string
    pass       string
    ssl        bool
    requireTLS bool
    TLSConfig  *tls.Config
}

type MailService struct {
    infoEmailSender         string
    orderEmailSender        string
    notificationEmailSender string
    emailFromName           string
    emailAppName            string
    emailToAddress          string
    emailMasterSvc          *emailMaster.EmailMasterService
    mailConfig              MailConfig
    logoConfig              LogoConfig
}

func NewMailService() *MailService {
    return &MailService{
        infoEmailSender:         config.Config("email.sender.info"),
        orderEmailSender:        config.Config("email.sender.order"),
        notificationEmailSender: config.Config("email.sender.notification"),
        emailFromName:           config.Config("email.from.name"),
        emailAppName:            config.Config("email.app.name"),
        emailToAddress:          config.Config("email.to.address"),
        emailMasterSvc:          emailMaster.EmailMasterSvc,
        mailConfig:              getConfig(),
        logoConfig:              getLogoConfig(),
    }
}

func (s *MailService) init() *gomail.Dialer {
    config := s.mailConfig
    d := gomail.NewDialer(
        config.host,
        config.port,
        config.user,
        config.pass,
    )
    d.SSL = config.ssl
    d.TLSConfig = config.TLSConfig
    return d
}

func (s *MailService) SendUserResetPassword(o *model.ApplicationUser) error {
    logo := s.logoConfig
    cid := "unique@logo"

    html := fmt.Sprintf(`
    <div style='text-align:center'><img src='cid:%s' style='width:%s;height:%s' /></div>
    <div style='font-family:Arial;font-size:14px'>
    Please sign in to Mobile Application with <br/><br/>
    Username: <strong>%s</strong><br/><br/>
    New Password: <strong>%s</strong><br/><br/>
    Thank you.<br/>
    %s Team
    </div>
    `, cid, logo.Width, logo.Height, o.Username.String, o.Password.String, s.emailAppName)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", s.emailFromName, s.infoEmailSender))
    m.SetHeader("To", o.Email.String)
    m.SetHeader("Subject", "Password Reset")
    m.SetBody("text/html", html)
    m.Embed(filepath.Join("ref", "img", logo.File), gomail.SetHeader(map[string][]string{
        "Content-ID":          {"<unique@logo>"},
        "Content-Disposition": {"inline; filename=\"" + logo.File + "\""},
    }))

    err := s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendAdminResetPassword(o *model.AdminUser) error {
    html := fmt.Sprintf(`
    <div style='font-family:Arial;font-size:14px'>
    Please sign in to Healthcare Web Admin Portal with: <br/><br/>
    Username: <strong>%s</strong><br/><br/>
    New Password: <strong>%s</strong><br/><br/>
    Thank you.<br/>
    %s Team
    </div>
    `, o.Username.String, o.Password.String, s.emailAppName)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", s.emailFromName, s.infoEmailSender))
    m.SetHeader("To", o.Email.String)
    m.SetHeader("Subject", "Password Reset")
    m.SetBody("text/html", html)

    err := s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendResetPassword(o *model.ApplicationUser) error {
    html := strings.NewReplacer(
        "{{first_name}}", o.FirstName.String,
        "{{email_appname}}", s.emailAppName,
        "{{verification_code}}", o.VerificationCode.String,
    ).Replace(constants.EmailTemplateConstantResetPw)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", s.emailFromName, s.infoEmailSender))
    m.SetHeader("To", o.Username.String)
    m.SetHeader("Subject", fmt.Sprintf("%s - Reset Password", s.emailAppName))
    m.SetBody("text/html", html)

    err := s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendAdminSignUp(o *model.AdminUser) error {
    html := fmt.Sprintf(`
    <div style='font-family:Arial;font-size:14px'>
    Please sign in to Healthcare Web Admin Portal with: <br/><br/>
    Username: <strong>%s</strong><br/><br/>
    Password: <strong>%s</strong><br/><br/>
    Thank you.<br />
    %s Team
    </div>
    `, o.Username.String, o.Password.String, s.emailAppName)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", s.emailFromName, s.infoEmailSender))
    m.SetHeader("To", o.Email.String)
    m.SetHeader("Subject", fmt.Sprintf("Welcome %s", o.FirstName.String))
    m.SetBody("text/html", html)

    err := s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendSignUp(o *model.ApplicationUser, newEmail string) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendSignUp")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{first_name}}", o.FirstName.String,
        "{{email_appname}}", s.emailAppName,
        "{{verification_code}}", o.VerificationCode.String,
    ).Replace(em.EmailTemplate.String)

    to := o.Username.String
    if newEmail != "" {
        to = newEmail
    }

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", to)
    m.SetHeader("Subject", em.EmailSubject.String)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendSignUpSuccess(o mm.MailSignUpSuccess) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendSignUpSuccess")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{patient_name}}", o.PatientName,
        "{{username}}", o.Username,
    ).Replace(em.EmailTemplate.String)

    to := o.Email
    if em.EmailRecipient.Valid {
        to = em.EmailRecipient.String
    }

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", to)
    m.SetHeader("Subject", em.EmailSubject.String)
    m.SetBody("text/html", html)
    
    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendGuestAppointmentConfirmationToPatient(o mm.MailGuestAppointmentConfirmation) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendGuestAppointmentConfirmationToPatient")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{guest_name}}", o.GuestName,
        "{{doctor_name}}", o.DoctorName,
        "{{appointment_date}}", o.AppointmentDate,
        "{{appointment_time}}", o.AppointmentTime,
        "{{clinic_location}}", o.ClinicLocation,
    ).Replace(em.EmailTemplate.String)

    subject := strings.NewReplacer(
        "{{doctor_name}}", o.DoctorName,
    ).Replace(em.EmailSubject.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", o.Email)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendGuestAppointmentConfirmationToIH(o mm.MailGuestAppointmentConfirmation) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendGuestAppointmentConfirmationToPatient")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{guest_name}}", o.GuestName,
        "{{doctor_name}}", o.DoctorName,
        "{{appointment_date}}", o.AppointmentDate,
        "{{appointment_time}}", o.AppointmentTime,
        "{{clinic_location}}", o.ClinicLocation,
    ).Replace(em.EmailTemplate.String)

    to := o.Email
    if em.EmailRecipient.Valid {
        to = em.EmailRecipient.String
    }

    subject := strings.NewReplacer(
        "{{doctor_name}}", o.DoctorName,
    ).Replace(em.EmailSubject.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendPatientFeedbackSubmitted() error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendPatientFeedbackSubmitted")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := em.EmailTemplate.String

    subject := em.EmailSubject.String

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", em.EmailRecipient.String)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendTest(senderEmail string) error {
    html := "Test Email from Nova"

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", s.emailFromName, s.infoEmailSender))
    m.SetHeader("To", senderEmail)
    m.SetHeader("Subject", "Test Email from Nova")
    m.SetBody("text/html", html)

    err := s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendLittleKids(o mm.MailLittleKids) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendLittleKids")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{name}}", o.KidsName,
    ).Replace(em.EmailTemplate.String)

    to := o.Email
    if em.EmailRecipient.Valid {
        to = em.EmailRecipient.String
    }

    subject := em.EmailSubject.String

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendGoldenPearl(o mm.MailGoldenPearl) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendGoldenPearl")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{name}}", o.GoldenName,
    ).Replace(em.EmailTemplate.String)

    to := o.Email
    if em.EmailRecipient.Valid {
        to = em.EmailRecipient.String
    }

    subject := em.EmailSubject.String

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendClubsEventRegistrationToMember(o mm.MailClubsEventRegistrationToMember) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendClubsEventRegistrationToMember")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{member_name}}", o.MemberName,
        "{{event_name}}", o.ActivityName,
    ).Replace(em.EmailTemplate.String)

    subject := strings.NewReplacer(
        "{{activity_name}}", o.ActivityName,
    ).Replace(em.EmailSubject.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", o.Email)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendClubsEventRegistrationToIH(o mm.MailClubsEventRegistrationToMember) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendClubsEventRegistrationToIH")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{member_name}}", o.MemberName,
        "{{event_name}}", o.ActivityName,
    ).Replace(em.EmailTemplate.String)

    to := o.Email
    if em.EmailRecipient.Valid {
        to = em.EmailRecipient.String
    }

    subject := strings.NewReplacer(
        "{{activity_name}}", o.ActivityName,
    ).Replace(em.EmailSubject.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendLogisticConfirmation(o *lgm.LogisticRequest) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendLogisticConfirmation")
    if err != nil {
        utils.LogError(err)
        return err
    }

    withCompanion := "No"
    if o.VisitWithCompanion.String == "Y" {
        withCompanion = "Yes"
    }

    companionName := "-"
    if o.VisitWithCompanion.String == "Y" {
        companionName = o.CompanionName.String
    }

    html := strings.NewReplacer(
        "{{requester_name}}", o.RequesterName.String,
        "{{requester_prn}}", o.RequesterPrn.String,
        "{{with_companion}}", withCompanion,
        "{{companion_name}}", companionName,
        "{{pickup_datetime}}", fmt.Sprintf("%s %s", o.RequestedPickupDate.String, o.RequestedPickupTime.String),
        "{{logistic_number}}", o.LogisticRequestNumber.String,
    ).Replace(em.EmailTemplate.String)

    subject := strings.NewReplacer(
        "{{logistic_request_number}}", o.LogisticRequestNumber.String,
    ).Replace(em.EmailSubject.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", em.EmailRecipient.String)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendLogisticCancellation(o *lgm.LogisticRequest) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendLogisticCancellation")
    if err != nil {
        utils.LogError(err)
        return err
    }

    if o.RequestedPickupDate.Valid && o.RequestedPickupDate.String != "" {
        g, _ := goment.New(o.RequestedPickupDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.RequestedPickupDate = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }

    withCompanion := "No"
    if o.VisitWithCompanion.String == "Y" {
        withCompanion = "Yes"
    }

    companionName := "-"
    if o.VisitWithCompanion.String == "Y" {
        companionName = o.CompanionName.String
    }

    html := strings.NewReplacer(
        "{{requester_name}}", o.RequesterName.String,
        "{{requester_prn}}", o.RequesterPrn.String,
        "{{with_companion}}", withCompanion,
        "{{companion_name}}", companionName,
        "{{pickup_datetime}}", fmt.Sprintf("%s %s", o.RequestedPickupDate.String, o.RequestedPickupTime.String),
        "{{logistic_number}}", o.LogisticRequestNumber.String,
    ).Replace(constants.EmailTemplateConstantLogisticRequestCancellation)

    subject := strings.NewReplacer(
        "{{logistic_request_number}}", o.LogisticRequestNumber.String,
    ).Replace(em.EmailSubject.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", em.EmailRecipient.String)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendSuccessOutstandingBillPayment(o mm.MailSuccessOutstandingBillPayment) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendSuccessOutstandingBillPayment")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{amount}}", fmt.Sprintf("RM%s", o.Amount),
        "{{payment_method}}", o.PaymentMethod,
        "{{bill_number}}", o.BillNumber,
        "{{invoice_number}}", o.InvoiceNumber,
    ).Replace(em.EmailTemplate.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", o.Email)
    m.SetHeader("Subject", em.EmailSubject.String)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func (s *MailService) SendSuccessPackagePayment(o mm.MailSuccessPackagePayment) error {
    em, err := s.emailMasterSvc.FindByEmailFunctionName("sendSuccessPackagePayment")
    if err != nil {
        utils.LogError(err)
        return err
    }

    html := strings.NewReplacer(
        "{{patient_name}}", o.PatientName,
        "{{order_number}}", o.OrderNumber,
        "{{date_of_purchase}}", o.DateOfPurchase,
        "{{product_name}}", o.ProductName,
        "{{product_quantity}}", o.ProductQuantity,
        "{{product_price}}", fmt.Sprintf("RM%s", o.ProductPrice),
        "{{subtotal_price}}", fmt.Sprintf("RM%s", o.SubtotalPrice),
        "{{payment_method}}", o.PaymentMethod,
        "{{total_price}}", fmt.Sprintf("RM%s", o.TotalPrice),
        "{{package_expiry_date}}", o.PackageExpiryDate,
        "{{billing_address}}", o.BillingAddress,
    ).Replace(em.EmailTemplate.String)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", em.EmailSenderName.String, em.EmailSender.String))
    m.SetHeader("To", o.Email)
    m.SetHeader("Subject", em.EmailSubject.String)
    m.SetBody("text/html", html)

    err = s.init().DialAndSend(m)
    if err != nil {
        utils.LogError(err)
    }

    return err
}

func getConfig() MailConfig {
    port := config.Config("mail.port")
    iport, _ := strconv.Atoi(port)
    requireTLS := config.Config("mail.requireTLS") == "true"
    return MailConfig{
        host:       config.Config("mail.host"),
        port:       iport,
        user:       config.Config("mail.username"),
        pass:       config.Config("mail.password"),
        ssl:        iport == 465,
        requireTLS: requireTLS,
        TLSConfig:  &tls.Config{InsecureSkipVerify: true},
    }
}

func getLogoConfig() LogoConfig {
    return LogoConfig{
        File:   config.Config("logo.file"),
        Width:  config.Config("logo.width"),
        Height: config.Config("logo.height"),
    }
}
