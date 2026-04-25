package model

type EventType string

const (
	EventApplicationSubmitted   EventType = "application.submitted"
	EventApplicationShortlisted EventType = "application.shortlisted"
	EventApplicationPitching    EventType = "application.pitching"
	EventApplicationFunded      EventType = "application.funded"
	EventApplicationRejected    EventType = "application.rejected"
)

type ApplicationEvent struct {
	EventType     EventType `json:"event_type"`
	ApplicationID string    `json:"application_id"`
	StartupID     string    `json:"startup_id"`
	StartupName   string    `json:"startup_name"`
	VCID          string    `json:"vc_id"`
	VCName        string    `json:"vc_name"`
	FounderEmail  string    `json:"founder_email"`
	VCEmail       string    `json:"vc_email"`
	Status        string    `json:"status"`
	RejectionNote string    `json:"rejection_note"`
}
