package domain

type Update struct {
	Type       UpdateType
	ChatID     int64
	UserID     int64
	Text       string
	Payload    string
	MessageID  string
	Language   Language
	Raw        any
}
