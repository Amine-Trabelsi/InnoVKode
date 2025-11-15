package domain

import "time"

type Session struct {
	ChatID               int64
	UserID               int64
	Language             Language
	Role                 Role
	Stage                Stage
	Email                string
	Profile              *UserProfile
	CurrentMenu          string
	PendingAction        *PendingAction
	PendingOTP           *PendingOTP
	PendingEventID       int64
	NotificationsEnabled bool
	LastActivity         time.Time
}

type PendingOTP struct {
	Code      string
	ExpiresAt time.Time
}

type PendingAction struct {
	ID        ActionID
	Step      int
	Data      map[string]string
	StartedAt time.Time
}
