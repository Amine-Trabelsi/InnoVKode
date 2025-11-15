package domain

import "strings"

const (
	ParseModeMarkdown = "markdown"
)

type Language string

const (
	LanguageRU Language = "ru"
	LanguageEN Language = "en"
)

func (l Language) Normalize() Language {
	switch strings.ToLower(string(l)) {
	case "en", "eng", "english":
		return LanguageEN
	default:
		return LanguageRU
	}
}

type Role string

const (
	RoleApplicant  Role = "applicant"
	RoleStudent    Role = "student"
	RoleEmployee   Role = "employee"
	RoleLeadership Role = "leadership"
)

type Stage string

const (
	StageInit           Stage = "init"
	StageSelectLanguage Stage = "select_language"
	StageChooseAuthMode Stage = "choose_auth_mode"
	StageCollectEmail   Stage = "collect_email"
	StageAwaitOTP       Stage = "await_otp"
	StageMainMenu       Stage = "main_menu"
)

type UpdateType string

const (
	UpdateTypeMessage  UpdateType = "message"
	UpdateTypeCallback UpdateType = "callback"
	UpdateTypeContact  UpdateType = "contact"
	UpdateTypeUnknown  UpdateType = "unknown"
)

type ButtonKind string

const (
	ButtonKindCallback ButtonKind = "callback"
	ButtonKindLink     ButtonKind = "link"
	ButtonKindCommand  ButtonKind = "command"
)

type ButtonStyle string

const (
	ButtonStylePrimary   ButtonStyle = "primary"
	ButtonStyleSecondary ButtonStyle = "secondary"
	ButtonStyleDanger    ButtonStyle = "danger"
	ButtonStyleInfo      ButtonStyle = "info"
)
