package domain

type OutgoingMessage struct {
	Text      string
	ParseMode string
	Keyboard  *Keyboard
	Reset     bool
}

type Keyboard struct {
	Rows [][]KeyboardButton
}

type KeyboardButton struct {
	Label   string
	Payload string
	Kind    ButtonKind
	Style   ButtonStyle
	URL     string
}
