package model

import (
    "github.com/guregu/null/v6"
)

type OneSignalNotification struct {
    NotificationID          null.Int64  `json:"notification_id" db:"NOTIFICATION_ID" swaggertype:"integer"`
    UserID                  null.Int64  `json:"user_id" db:"USER_ID" swaggertype:"integer"`
    VisitType               null.String `json:"visitType" db:"VISIT_TYPE" swaggertype:"string"`
    AccountNo               null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    NotificationTitle       null.String `json:"notificationTitle" db:"NOTIFICATION_TITLE" swaggertype:"string"`
    MsgType                 null.String `json:"msgType" db:"MSG_TYPE" swaggertype:"string"`
    ShortMessage            null.String `json:"shortMessage" db:"SHORT_MESSAGE" swaggertype:"string"`
    FullMessage             null.String `json:"fullMessage" db:"FULL_MESSAGE" swaggertype:"string"`
    MasterID                null.Int64  `json:"master_id" db:"MASTER_ID" swaggertype:"integer"`
    IsSeenV                 null.String `json:"-" db:"IS_SEEN" swaggertype:"string"`
    IsSeen                  bool        `json:"isSeen"`
    DateCreate              null.String `json:"dateCreate" db:"DATE_SENT" swaggertype:"string"`
    DateSent                null.String `json:"dateSent" swaggertype:"string"`
    GuestPlayerID           null.String `json:"guestPlayerId" db:"GUEST_PLAYER_ID" swaggertype:"string"`
    DateCreate2             null.String `json:"-" db:"DATE_CREATE" swaggertype:"string"`
    DateSeen                null.String `json:"-" db:"DATE_SEEN" swaggertype:"string"`
    OneSignalMsg            null.String `json:"-" db:"ONESIGNAL_MSG" swaggertype:"string"`
    OneSignalNotificationID null.String `json:"-" db:"ONESIGNAL_NOTIFICATION_ID" swaggertype:"string"`
}

func (o *OneSignalNotification) Set() {
    if o.IsSeenV.String == "Y" {
        o.IsSeen = true
    } else {
        o.IsSeen = false
    }
}
