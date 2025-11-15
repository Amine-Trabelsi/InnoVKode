package bot

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/escalopa/inno-vkode/internal/config"
	"github.com/escalopa/inno-vkode/internal/domain"
	"github.com/escalopa/inno-vkode/internal/ports"
	"github.com/escalopa/inno-vkode/internal/state"
)

const (
	payloadNavPrefix  = "nav:"
	payloadActionPref = "act:"
	payloadLangPref   = "lang:"
	payloadAuthPref   = "auth:"
)

type Service struct {
	cfg       *config.Config
	log       zerolog.Logger
	backend   ports.Backend
	messenger ports.Messenger
	email     ports.EmailSender
	store     state.Store

	menus *MenuRegistry
	forms map[domain.ActionID]FormDefinition

	now       func() time.Time
	otpDigits int
	otpExpiry time.Duration
}

func New(cfg *config.Config, log zerolog.Logger, backend ports.Backend, messenger ports.Messenger, email ports.EmailSender, store state.Store) *Service {
	s := &Service{
		cfg:       cfg,
		log:       log,
		backend:   backend,
		messenger: messenger,
		email:     email,
		store:     store,
		menus:     buildMenuRegistry(),
		now:       time.Now,
		otpDigits: 6,
		otpExpiry: cfg.OTPExpiry,
	}
	s.forms = s.buildForms()
	return s
}

func (s *Service) Start(ctx context.Context) error {
	s.log.Info().Msg("starting MAX bot service")
	return s.messenger.Start(ctx, s.handleUpdate)
}

func (s *Service) handleUpdate(ctx context.Context, upd domain.Update) error {
	if upd.ChatID == 0 {
		return nil
	}
	sess := s.ensureSession(upd.ChatID)
	if upd.UserID != 0 {
		sess.UserID = upd.UserID
	}

	if handled, err := s.handleGlobalCommands(ctx, sess, upd); handled || err != nil {
		return err
	}

	switch sess.Stage {
	case domain.StageSelectLanguage:
		return s.handleLanguageSelection(ctx, sess, upd)
	case domain.StageChooseAuthMode:
		return s.handleChooseAuthMode(ctx, sess, upd)
	case domain.StageCollectEmail:
		return s.handleEmailCollection(ctx, sess, upd)
	case domain.StageAwaitOTP:
		return s.handleOTPSubmission(ctx, sess, upd)
	case domain.StageMainMenu:
		return s.handleMainMenu(ctx, sess, upd)
	default:
		return s.sendLanguagePrompt(ctx, sess, true)
	}
}

func (s *Service) handleGlobalCommands(ctx context.Context, sess *domain.Session, upd domain.Update) (bool, error) {
	text := strings.TrimSpace(strings.ToLower(upd.Text))
	switch text {
	case "/start":
		s.resetSession(sess)
		sess.Stage = domain.StageSelectLanguage
		s.saveSession(sess)
		return true, s.sendLanguagePrompt(ctx, sess, true)
	case "/language":
		sess.Stage = domain.StageSelectLanguage
		sess.PendingAction = nil
		s.saveSession(sess)
		return true, s.sendLanguagePrompt(ctx, sess, false)
	case "/help":
		helpText := s.t(sess.Language,
			"🆘 **Помощь**\n\n📋 **Команды:**\n• /start - Перезапустить бота\n• /language - Изменить язык\n• /help - Показать эту справку\n• /cancel - Отменить текущее действие\n\n❓ **Часто задаваемые вопросы:**\n• Как войти? Используйте /start и выберите 'Войти'.\n• Забыли пароль? Свяжитесь с поддержкой через меню.\n• Проблемы с ботом? Опишите в '🐞 Сообщить об ошибке'.\n\n💬 Для дополнительной помощи используйте меню 'ℹ️ Поддержка'.",
			"🆘 **Help**\n\n📋 **Commands:**\n• /start - Restart the bot\n• /language - Change language\n• /help - Show this help\n• /cancel - Cancel current action\n\n❓ **FAQs:**\n• How to login? Use /start and choose 'Login'.\n• Forgot password? Contact support via menu.\n• Bot issues? Report in '🐞 Report issue'.\n\n💬 For more help, use 'ℹ️ Support' menu.")
		return true, s.reply(ctx, sess, helpText)
	case "/cancel", "cancel", "отмена":
		if sess.PendingAction != nil {
			sess.PendingAction = nil
			s.saveSession(sess)
			return true, s.reply(ctx, sess, s.t(sess.Language, "🚫 Действие отменено.", "🚫 Action cancelled."))
		}
	}
	return false, nil
}

func (s *Service) handleLanguageSelection(ctx context.Context, sess *domain.Session, upd domain.Update) error {
	code := ""
	if upd.Type == domain.UpdateTypeCallback && strings.HasPrefix(upd.Payload, payloadLangPref) {
		code = strings.TrimPrefix(upd.Payload, payloadLangPref)
	}
	if code == "" && upd.Text != "" {
		code = strings.ToLower(strings.TrimSpace(upd.Text))
	}
	switch code {
	case "ru", "rus", "russian":
		sess.Language = domain.LanguageRU
	case "en", "eng", "english":
		sess.Language = domain.LanguageEN
	default:
		return s.sendLanguagePrompt(ctx, sess, false)
	}

	if sess.Profile != nil {
		// Switching language while logged in
		sess.Stage = domain.StageMainMenu
		sess.PendingAction = nil
		sess.PendingEventID = 0
		s.saveSession(sess)
		greeting := s.t(sess.Language, "🌐 Язык интерфейса изменён!", "🌐 Interface language changed!")
		if err := s.reply(ctx, sess, greeting); err != nil {
			return err
		}
		return s.sendCurrentMenu(ctx, sess)
	}

	sess.Stage = domain.StageChooseAuthMode
	s.saveSession(sess)
	return s.sendAuthModePrompt(ctx, sess)
}

func (s *Service) handleChooseAuthMode(ctx context.Context, sess *domain.Session, upd domain.Update) error {
	mode := ""
	if upd.Type == domain.UpdateTypeCallback && strings.HasPrefix(upd.Payload, payloadAuthPref) {
		mode = strings.TrimPrefix(upd.Payload, payloadAuthPref)
	}
	if mode == "" && upd.Text != "" {
		mode = strings.ToLower(strings.TrimSpace(upd.Text))
	}
	switch mode {
	case "guest", "гость":
		sess.Role = domain.RoleApplicant
		sess.Stage = domain.StageMainMenu
		if root := s.menus.Root(sess.Role); root != nil {
			sess.CurrentMenu = root.ID
		}
		s.saveSession(sess)
		greeting := s.t(sess.Language, "🎉 Добро пожаловать в гостевой режим!\n\n👋 Вы можете просматривать информацию о поступлении и общих сервисах.", "🎉 Welcome to guest mode!\n\n👋 You can browse admission info and general services.")
		if err := s.reply(ctx, sess, greeting); err != nil {
			return err
		}
		return s.sendCurrentMenu(ctx, sess)
	case "login", "войти":
		sess.Stage = domain.StageCollectEmail
		s.saveSession(sess)
		return s.reply(ctx, sess, s.t(sess.Language, "📧 Введите ваш университетский email для авторизации:\n\nПример: student@university.edu", "📧 Enter your university email for login:\n\nExample: student@university.edu"))
	default:
		return s.sendAuthModePrompt(ctx, sess)
	}
}

func (s *Service) handleEmailCollection(ctx context.Context, sess *domain.Session, upd domain.Update) error {
	if upd.Text == "" {
		return nil
	}
	email := strings.TrimSpace(strings.ToLower(upd.Text))
	if !strings.Contains(email, "@") || len(email) < 5 {
		return s.reply(ctx, sess, s.t(sess.Language, "❌ Неверный формат email. Попробуйте снова.", "❌ Invalid email format. Please try again."))
	}
	sess.Email = email
	code := s.generateOTP()
	sess.PendingOTP = &domain.PendingOTP{
		Code:      code,
		ExpiresAt: s.now().Add(s.otpExpiry),
	}
	if err := s.email.SendOTP(context.Background(), email, code); err != nil {
		s.log.Error().Err(err).Str("email", email).Msg("failed to send otp")
	}
	sess.Stage = domain.StageAwaitOTP
	s.saveSession(sess)
	return s.reply(ctx, sess, s.t(sess.Language, "🔐 Мы отправили 6-значный код подтверждения на вашу почту!\n\n📨 Проверьте папку \"Входящие\" и введите код:", "🔐 We sent a 6-digit verification code to your email!\n\n📨 Check your inbox and enter the code:"))
}

func (s *Service) handleOTPSubmission(ctx context.Context, sess *domain.Session, upd domain.Update) error {
	if upd.Text == "" || sess.PendingOTP == nil {
		return nil
	}
	if s.now().After(sess.PendingOTP.ExpiresAt) {
		sess.PendingOTP = nil
		s.saveSession(sess)
		return s.reply(ctx, sess, s.t(sess.Language, "⏰ Код подтверждения истёк. Начните процесс заново с /start", "⏰ Verification code expired. Start over with /start"))
	}
	if strings.TrimSpace(upd.Text) != sess.PendingOTP.Code {
		return s.reply(ctx, sess, s.t(sess.Language, "❌ Неверный код подтверждения. Попробуйте ещё раз.", "❌ Incorrect verification code. Try again."))
	}

	profile, err := s.backend.GetUserByEmail(ctx, sess.Email)
	if err != nil {
		s.log.Warn().Err(err).Str("email", sess.Email).Msg("user not found, fallback applicant")
		profile = &domain.UserProfile{
			Email:     sess.Email,
			NameRU:    sess.Email,
			NameEN:    sess.Email,
			Role:      domain.RoleApplicant,
			Language:  sess.Language,
			CreatedAt: s.now(),
		}
	}
	sess.Profile = profile
	sess.Role = profile.Role
	sess.Stage = domain.StageMainMenu
	sess.PendingOTP = nil
	if root := s.menus.Root(sess.Role); root != nil {
		sess.CurrentMenu = root.ID
	}
	s.saveSession(sess)

	greeting := s.t(sess.Language, fmt.Sprintf("🎊 Добро пожаловать, %s!\n\n✅ Авторизация успешна. Доступ ко всем сервисам открыт.", profile.NameRU), fmt.Sprintf("🎊 Welcome, %s!\n\n✅ Login successful. Full access to all services.", profile.NameEN))
	if err := s.reply(ctx, sess, greeting); err != nil {
		return err
	}
	return s.sendCurrentMenu(ctx, sess)
}

func (s *Service) handleMainMenu(ctx context.Context, sess *domain.Session, upd domain.Update) error {
	if sess.PendingAction != nil && upd.Type == domain.UpdateTypeMessage && strings.TrimSpace(upd.Text) != "" {
		return s.handleFormInput(ctx, sess, strings.TrimSpace(upd.Text))
	}

	if upd.Type == domain.UpdateTypeCallback {
		switch {
		case strings.HasPrefix(upd.Payload, payloadNavPrefix):
			return s.navigateTo(ctx, sess, strings.TrimPrefix(upd.Payload, payloadNavPrefix))
		case strings.HasPrefix(upd.Payload, payloadActionPref):
			return s.executeAction(ctx, sess, domain.ActionID(strings.TrimPrefix(upd.Payload, payloadActionPref)))
		case strings.HasPrefix(upd.Payload, payloadLangPref):
			sess.Stage = domain.StageSelectLanguage
			sess.PendingAction = nil
			s.saveSession(sess)
			return s.sendLanguagePrompt(ctx, sess, false)
		case strings.HasPrefix(upd.Payload, "event_select:"):
			return s.handleEventSelect(ctx, sess, strings.TrimPrefix(upd.Payload, "event_select:"))
		case strings.HasPrefix(upd.Payload, "event_mode:"):
			return s.handleEventMode(ctx, sess, strings.TrimPrefix(upd.Payload, "event_mode:"))
		case strings.HasPrefix(upd.Payload, "cancel_event:"):
			return s.handleCancelEvent(ctx, sess, strings.TrimPrefix(upd.Payload, "cancel_event:"))
		}
	}

	if strings.TrimSpace(upd.Text) != "" {
		return s.reply(ctx, sess, s.t(sess.Language, "Используйте кнопки меню или команды /start, /language, /help.", "Use the menu buttons or /start /language /help commands."))
	}
	return nil
}

func (s *Service) navigateTo(ctx context.Context, sess *domain.Session, nodeID string) error {
	node := s.menus.Node(nodeID)
	if node == nil {
		return nil
	}
	root := s.menus.Root(sess.Role)
	if root == nil {
		return nil
	}
	sess.CurrentMenu = node.ID
	s.saveSession(sess)
	return s.sendMenuNode(ctx, sess, node)
}

func (s *Service) executeAction(ctx context.Context, sess *domain.Session, action domain.ActionID) error {
	if action == domain.ActionSwitchLanguage {
		sess.Stage = domain.StageSelectLanguage
		sess.PendingAction = nil
		sess.PendingEventID = 0
		s.saveSession(sess)
		return s.sendLanguagePrompt(ctx, sess, false)
	}

	if form, ok := s.forms[action]; ok {
		return s.startForm(ctx, sess, action, form)
	}

	msg, err := s.handleAction(ctx, sess, action)
	if err != nil {
		s.log.Error().Err(err).Str("action", string(action)).Msg("action handler failed")
		return s.reply(ctx, sess, s.t(sess.Language, "Произошла ошибка. Попробуйте позже.", "Something went wrong, please try later."))
	}
	return s.replyMessage(ctx, sess, msg)
}

func (s *Service) sendCurrentMenu(ctx context.Context, sess *domain.Session) error {
	node := s.menus.Node(sess.CurrentMenu)
	if node == nil {
		if root := s.menus.Root(sess.Role); root != nil {
			sess.CurrentMenu = root.ID
			s.saveSession(sess)
			node = root
		} else {
			return nil
		}
	}
	return s.sendMenuNode(ctx, sess, node)
}

func (s *Service) sendMenuNode(ctx context.Context, sess *domain.Session, node *MenuNode) error {
	title := node.TitleText(sess.Language)
	desc := node.DescriptionText(sess.Language)
	text := title
	if desc != "" {
		text = fmt.Sprintf("%s\n\n%s", title, desc)
	}
	msg := domain.OutgoingMessage{
		Text:     text,
		Keyboard: s.buildMenuKeyboard(sess, node),
	}
	return s.replyMessage(ctx, sess, msg)
}

func (s *Service) buildMenuKeyboard(sess *domain.Session, node *MenuNode) *domain.Keyboard {
	if node == nil || len(node.Children) == 0 {
		return nil
	}
	kb := &domain.Keyboard{}
	for _, child := range node.Children {
		btn := domain.KeyboardButton{
			Label: child.TitleText(sess.Language),
			Style: domain.ButtonStylePrimary,
			Kind:  domain.ButtonKindCallback,
		}
		if child.Action != "" {
			btn.Payload = payloadActionPref + string(child.Action)
		} else {
			btn.Payload = payloadNavPrefix + child.ID
		}
		kb.Rows = append(kb.Rows, []domain.KeyboardButton{btn})
	}

	backRow := []domain.KeyboardButton{}
	if node.ParentID != "" {
		backRow = append(backRow, domain.KeyboardButton{
			Label:   s.t(sess.Language, "⬅ Назад", "⬅ Back"),
			Style:   domain.ButtonStyleSecondary,
			Kind:    domain.ButtonKindCallback,
			Payload: payloadNavPrefix + node.ParentID,
		})
	}
	if root := s.menus.Root(sess.Role); root != nil && node.ID != root.ID {
		backRow = append(backRow, domain.KeyboardButton{
			Label:   s.t(sess.Language, "🏠 В меню", "🏠 Main menu"),
			Style:   domain.ButtonStyleSecondary,
			Kind:    domain.ButtonKindCallback,
			Payload: payloadNavPrefix + root.ID,
		})
	}
	if len(backRow) > 0 {
		kb.Rows = append(kb.Rows, backRow)
	}
	return kb
}

func (s *Service) sendLanguagePrompt(ctx context.Context, sess *domain.Session, reset bool) error {
	msg := domain.OutgoingMessage{
		Text: s.t(sess.Language, "🌐 Выберите язык интерфейса\n\n🇷🇺 Русский - Полная локализация\n🇬🇧 English - Full localization", "🌐 Choose interface language\n\n🇷🇺 Русский - Full localization\n🇬🇧 English - Полная локализация"),
		Keyboard: &domain.Keyboard{
			Rows: [][]domain.KeyboardButton{
				{
					{Label: "🇷🇺 Русский", Kind: domain.ButtonKindCallback, Payload: payloadLangPref + "ru", Style: domain.ButtonStylePrimary},
					{Label: "🇬🇧 English", Kind: domain.ButtonKindCallback, Payload: payloadLangPref + "en", Style: domain.ButtonStylePrimary},
				},
			},
		},
		Reset: reset,
	}
	return s.replyMessage(ctx, sess, msg)
}

func (s *Service) sendAuthModePrompt(ctx context.Context, sess *domain.Session) error {
	msg := domain.OutgoingMessage{
		Text:      s.t(sess.Language, "🎭 Выберите режим доступа:\n\n👤 **Гость** - Просмотр информации без авторизации\n📧 **Войти** - Полный доступ к личным сервисам", "🎭 Choose access mode:\n\n👤 **Guest** - Browse info without login\n📧 **Login** - Full access to personal services"),
		ParseMode: domain.ParseModeMarkdown,
		Keyboard: &domain.Keyboard{
			Rows: [][]domain.KeyboardButton{
				{
					{Label: "👤 " + s.t(sess.Language, "Гость", "Guest"), Kind: domain.ButtonKindCallback, Payload: payloadAuthPref + "guest", Style: domain.ButtonStyleSecondary},
					{Label: "📧 " + s.t(sess.Language, "Войти", "Login"), Kind: domain.ButtonKindCallback, Payload: payloadAuthPref + "login", Style: domain.ButtonStylePrimary},
				},
			},
		},
	}
	return s.replyMessage(ctx, sess, msg)
}

func (s *Service) reply(ctx context.Context, sess *domain.Session, text string) error {
	return s.replyMessage(ctx, sess, domain.OutgoingMessage{Text: text})
}

func (s *Service) replyMessage(ctx context.Context, sess *domain.Session, msg domain.OutgoingMessage) error {
	s.saveSession(sess)
	return s.messenger.Send(ctx, sess.ChatID, sess.UserID, msg)
}

func (s *Service) ensureSession(chatID int64) *domain.Session {
	if sess, ok := s.store.Get(chatID); ok {
		return sess
	}
	sess := &domain.Session{
		ChatID:   chatID,
		Language: domain.LanguageRU,
		Stage:    domain.StageInit,
		Role:     domain.RoleApplicant,
	}
	s.store.Save(sess)
	return sess
}

func (s *Service) resetSession(sess *domain.Session) {
	sess.Stage = domain.StageInit
	sess.PendingAction = nil
	sess.PendingOTP = nil
	sess.PendingEventID = 0
	sess.Profile = nil
	sess.Email = ""
	sess.Role = domain.RoleApplicant
	sess.CurrentMenu = ""
	s.saveSession(sess)
}

func (s *Service) saveSession(sess *domain.Session) {
	s.store.Save(sess)
}

func (s *Service) generateOTP() string {
	const digits = "0123456789"
	b := make([]byte, s.otpDigits)
	if _, err := rand.Read(b); err != nil {
		for i := range b {
			b[i] = digits[i%len(digits)]
		}
	} else {
		for i := range b {
			b[i] = digits[int(b[i])%len(digits)]
		}
	}
	return string(b)
}

func (s *Service) handleEventSelect(ctx context.Context, sess *domain.Session, eventIDStr string) error {
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		return s.reply(ctx, sess, "Invalid event ID.")
	}
	sess.PendingEventID = eventID
	s.saveSession(sess)
	kb := &domain.Keyboard{
		Rows: [][]domain.KeyboardButton{
			{
				{Label: "Attendee", Kind: domain.ButtonKindCallback, Payload: "event_mode:attendee", Style: domain.ButtonStylePrimary},
				{Label: "Participant", Kind: domain.ButtonKindCallback, Payload: "event_mode:participant", Style: domain.ButtonStylePrimary},
			},
		},
	}
	msg := domain.OutgoingMessage{
		Text:    "Choose your registration type:",
		Keyboard: kb,
	}
	return s.replyMessage(ctx, sess, msg)
}

func (s *Service) handleEventMode(ctx context.Context, sess *domain.Session, mode string) error {
	if sess.PendingEventID == 0 {
		return s.reply(ctx, sess, "No event selected.")
	}
	if sess.Profile == nil || sess.Profile.ID == 0 {
		sess.PendingEventID = 0
		s.saveSession(sess)
		return s.reply(ctx, sess, "Please login first.")
	}
	status, err := s.backend.RSVPEvent(ctx, sess.PendingEventID, sess.Profile.ID, mode, "")
	if err != nil {
		sess.PendingEventID = 0
		s.saveSession(sess)
		return s.reply(ctx, sess, "Registration failed: "+err.Error())
	}
	sess.PendingEventID = 0
	s.saveSession(sess)
	switch status {
	case "registered", "updated":
		return s.reply(ctx, sess, "Registration successful as "+mode+"!")
	case "already_registered":
		return s.reply(ctx, sess, "You are already registered for this event.")
	default:
		return s.reply(ctx, sess, "Registration status: "+status)
	}
}

func (s *Service) handleCancelEvent(ctx context.Context, sess *domain.Session, eventIDStr string) error {
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		return s.reply(ctx, sess, "Invalid event ID.")
	}
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return s.reply(ctx, sess, "Please login first.")
	}
	err = s.backend.CancelRSVP(ctx, eventID, sess.Profile.ID)
	if err != nil {
		return s.reply(ctx, sess, "Cancellation failed: "+err.Error())
	}
	return s.reply(ctx, sess, "Registration cancelled successfully!")
}

func (s *Service) t(lang domain.Language, ru, en string) string {
	if lang == domain.LanguageEN {
		return en
	}
	return ru
}
