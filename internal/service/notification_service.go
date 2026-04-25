package service

import (
	"fmt"
	"log"

	"github.com/r200a/notification-service/internal/model"
	"github.com/r200a/notification-service/pkg/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type NotificationService struct {
	cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (s *NotificationService) HandleEvent(event model.ApplicationEvent) {
	switch event.EventType {
	case model.EventApplicationSubmitted:
		s.sendEmail(
			event.VCEmail,
			"New Application Received",
			fmt.Sprintf(
				"Hi %s,\n\n%s has applied for funding.\n\nLogin to review their application.",
				event.VCName, event.StartupName,
			),
		)

	case model.EventApplicationShortlisted:
		s.sendEmail(
			event.FounderEmail,
			"You've been shortlisted!",
			fmt.Sprintf(
				"Hi,\n\nGreat news! %s has shortlisted %s for funding consideration.\n\nLogin to check your application status.",
				event.VCName, event.StartupName,
			),
		)

	case model.EventApplicationPitching:
		s.sendEmail(
			event.FounderEmail,
			"Pitch Session Scheduled",
			fmt.Sprintf(
				"Hi,\n\n%s wants to schedule a pitch session with %s.\n\nLogin to confirm details.",
				event.VCName, event.StartupName,
			),
		)

	case model.EventApplicationFunded:
		s.sendEmail(
			event.FounderEmail,
			"Congratulations — You're Funded!",
			fmt.Sprintf(
				"Hi,\n\nExciting news! %s has decided to fund %s.\n\nLogin to view your deal details.",
				event.VCName, event.StartupName,
			),
		)
		s.sendEmail(
			event.VCEmail,
			"Deal Confirmed",
			fmt.Sprintf(
				"Hi %s,\n\nYour investment in %s has been confirmed.\n\nLogin to view deal details.",
				event.VCName, event.StartupName,
			),
		)

	case model.EventApplicationRejected:
		note := "No additional notes provided."
		if event.RejectionNote != "" {
			note = event.RejectionNote
		}
		s.sendEmail(
			event.FounderEmail,
			"Application Update",
			fmt.Sprintf(
				"Hi,\n\nAfter careful consideration, %s has decided not to proceed with %s at this time.\n\nFeedback: %s\n\nKeep building — the right investor is out there.",
				event.VCName, event.StartupName, note,
			),
		)

	default:
		log.Printf("unknown event type: %s", event.EventType)
	}
}

func (s *NotificationService) sendEmail(to, subject, body string) {
	if s.cfg.SendGridAPIKey == "" {
		// no SendGrid key — just log for now
		log.Printf("EMAIL TO: %s | SUBJECT: %s | BODY: %s", to, subject, body)
		return
	}

	from := mail.NewEmail(s.cfg.FromName, s.cfg.FromEmail)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toEmail, body, "")

	client := sendgrid.NewSendClient(s.cfg.SendGridAPIKey)
	resp, err := client.Send(message)
	if err != nil {
		log.Printf("failed to send email to %s: %v", to, err)
		return
	}

	if resp.StatusCode >= 400 {
		log.Printf("sendgrid error %d for %s: %s", resp.StatusCode, to, resp.Body)
		return
	}

	log.Printf("email sent to %s — subject: %s", to, subject)
}
