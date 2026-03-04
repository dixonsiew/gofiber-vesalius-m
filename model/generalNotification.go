package model

import (
    "github.com/guregu/null/v6"
)

type GeneralNotification struct {
    NotificationMasterID null.Int64  `json:"notification_master_id" db:"NOTIFICATION_MASTER_ID" swaggertype:"integer"`
    NotificationTitle    null.String `json:"notificationTitle" db:"NOTIFICATION_TITLE" swaggertype:"string"`
    ShortMessage         null.String `json:"shortMessage" db:"SHORT_MESSAGE" swaggertype:"string"`
    FullMessage          null.String `json:"fullMessage" db:"FULL_MESSAGE" swaggertype:"string"`
    StartDate            null.String `json:"startDate" db:"START_DATE_TIME" swaggertype:"string"`
    EndDate              null.String `json:"endDate" db:"END_DATE_TIME" swaggertype:"string"`
    TargetAgeFrom        null.Int64  `json:"targetAgeFrom" db:"TARGET_AGE_FROM" swaggertype:"integer"`
    TargetAgeTo          null.Int64  `json:"targetAgeTo" db:"TARGET_AGE_TO" swaggertype:"integer"`
    TargetGender         null.String `json:"targetGender" db:"TARGET_GENDER" swaggertype:"string"`
    TargetNationality    null.String `json:"targetNationality" db:"TARGET_NATIONALITY" swaggertype:"string"`
    TargetCity           null.String `json:"targetCity" db:"TARGET_CITY" swaggertype:"string"`
    TargetState          null.String `json:"targetState" db:"TARGET_STATE" swaggertype:"string"`
    UserCreate           null.Int64  `json:"-" db:"USER_CREATE" swaggertype:"integer"`
    DateCreate           null.String `json:"-" db:"DATE_CREATE" swaggertype:"string"`
    UserUpdate           null.Int64  `json:"-" db:"USER_UPDATE" swaggertype:"integer"`
    DateUpdate           null.String `json:"-" db:"DATE_UPDATE" swaggertype:"string"`
}
