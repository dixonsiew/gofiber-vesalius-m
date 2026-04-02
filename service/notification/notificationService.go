package notification

import (
	"context"
	"vesaliusm/config"

	"github.com/OneSignal/onesignal-go-api"
)

var NotificationSvc *NotificationService = NewNotificationService()

type NotificationService struct {
    client *onesignal.APIClient
    ctx    context.Context
    appId  string
}

func NewNotificationService() *NotificationService {
    conf := onesignal.NewConfiguration()
    client := onesignal.NewAPIClient(conf)
    ctx := context.WithValue(
        context.Background(),
        onesignal.AppAuth,
        config.Config("onesignal.apikey"),
    )
    return &NotificationService{
        client: client,
        ctx:    ctx,
        appId:  config.Config("onesignal.appid"),
    }
}

func (s *NotificationService) CreateNotification(notification *onesignal.Notification) (*onesignal.CreateNotificationSuccessResponse, error) {
    response, _, err := s.client.DefaultApi.CreateNotification(s.ctx).Notification(*notification).Execute()
    if err != nil {
        return nil, err
    }
    return response, nil
}

func (s *NotificationService) GetAppId() string {
    return s.appId
}
