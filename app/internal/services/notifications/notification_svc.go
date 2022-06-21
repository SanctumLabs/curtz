package notifications

import (
	"errors"

	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
)

const (
	NotificationTypeEmail = "email"
	NotificationTypeSMS   = "sms"
	NotificationTypePush  = "push"
)

type NotificationService struct {
	emailSvc contracts.EmailService
	// smsSvc   SmsService
}

func NewNotificationSvc(emailSvc contracts.EmailService) *NotificationService {
	return &NotificationService{
		emailSvc: emailSvc,
	}
}

// SendNotification sends a notification message to given recipient with provided message
func (n *NotificationService) SendNotification(recipient, message, notifyType string) error {
	switch notifyType {
	case NotificationTypeEmail:
		return n.emailSvc.SendEmail(recipient, "", message)
	case NotificationTypeSMS:
		// return n.smsSvc.SendSms(recipient, message)
		panic("implement me")
	case NotificationTypePush:
		// return n.pushSvc.SendPush(recipient, message)
		panic("implement me")
	default:
		return errors.New("unknown notification type")
	}
}
