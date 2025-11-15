package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/escalopa/inno-vkode/internal/domain"
)

type FormField struct {
	Key      string
	Prompt   map[domain.Language]string
	Optional bool
}

func (f FormField) PromptText(lang domain.Language) string {
	if txt, ok := f.Prompt[lang]; ok && txt != "" {
		return txt
	}
	return f.Prompt[domain.LanguageRU]
}

type FormDefinition struct {
	Intro    map[domain.Language]string
	Fields   []FormField
	OnSubmit func(context.Context, *Service, *domain.Session, map[string]string) (domain.OutgoingMessage, error)
}

func (s *Service) buildForms() map[domain.ActionID]FormDefinition {
	return map[domain.ActionID]FormDefinition{
		domain.ActionLibrarySearch: {
			Intro: l("Поиск по названию или автору.", "Search by title or author."),
			Fields: []FormField{
				{Key: "query", Prompt: l("Введите запрос:", "Enter your query:")},
			},
			OnSubmit: submitLibrarySearch,
		},
		domain.ActionLibraryReserve: {
			Intro: l("Укажите идентификатор книги для брони.", "Provide the book ID to reserve."),
			Fields: []FormField{
				{Key: "book_id", Prompt: l("ID книги:", "Book ID:")},
			},
			OnSubmit: submitLibraryReserve,
		},
		domain.ActionDormMaintenance: {
			Intro: l("Создание заявки на ремонт.", "Create a maintenance ticket."),
			Fields: []FormField{
				{Key: "type", Prompt: l("Тип проблемы (электрика, сантехника...):", "Issue type (electrical, plumbing...):")},
				{Key: "details", Prompt: l("Подробности:", "Details:")},
			},
			OnSubmit: submitDormMaintenance,
		},
		domain.ActionContactSupport: {
			Intro: l("Создание обращения в поддержку.", "Create a support ticket."),
			Fields: []FormField{
				{Key: "category", Prompt: l("Категория (it, hr, study...):", "Category (it, hr, study...):")},
				{Key: "subject", Prompt: l("Тема:", "Subject:")},
				{Key: "description", Prompt: l("Описание:", "Description:")},
			},
			OnSubmit: submitSupportTicket,
		},
		domain.ActionReportIssue: {
			Intro: l("Расскажите об ошибке.", "Describe the issue."),
			Fields: []FormField{
				{Key: "subject", Prompt: l("Коротко о проблеме:", "Short summary:")},
				{Key: "description", Prompt: l("Детали:", "Details:")},
			},
			OnSubmit: submitBugReport,
		},
		domain.ActionFAQ: {
			Intro: l("Задайте вопрос, и мы дадим быстрый ответ.", "Ask your question for a quick answer."),
			Fields: []FormField{
				{Key: "question", Prompt: l("Ваш вопрос:", "Your question:")},
			},
			OnSubmit: submitFAQQuery,
		},
		domain.ActionAIQuery: {
			Intro: l("RAG-поиск по базе знаний.", "RAG knowledge base search."),
			Fields: []FormField{
				{Key: "question", Prompt: l("Вопрос:", "Question:")},
			},
			OnSubmit: submitAIQuery,
		},
		domain.ActionAISummary: {
			Intro: l("Отправьте текст для краткого конспекта (до 1500 символов).", "Paste text to summarize (up to 1500 chars)."),
			Fields: []FormField{
				{Key: "text", Prompt: l("Текст:", "Text:")},
			},
			OnSubmit: submitAISummary,
		},
		domain.ActionAIQuiz: {
			Intro: l("Сгенерируем контрольные вопросы.", "Generate quick quiz questions."),
			Fields: []FormField{
				{Key: "topic", Prompt: l("Тема / курс:", "Topic / course:")},
			},
			OnSubmit: submitAIQuiz,
		},
		domain.ActionAITranscription: {
			Intro: l("Укажите ссылку или идентификатор аудио.", "Provide audio reference or link."),
			Fields: []FormField{
				{Key: "audio", Prompt: l("Audio ref:", "Audio ref:")},
			},
			OnSubmit: submitAITranscription,
		},
		domain.ActionAdvisorChat: {
			Intro: l("Диалог с академическим советником.", "Chat with the academic advisor."),
			Fields: []FormField{
				{Key: "topic", Prompt: l("Тема:", "Topic:"), Optional: true},
				{Key: "prompt", Prompt: l("Ваш вопрос:", "Your question:")},
			},
			OnSubmit: submitAdvisorChat,
		},
		domain.ActionBusinessTripRequest: {
			Intro: l("Запрос на командировку.", "Business trip request."),
			Fields: []FormField{
				{Key: "destination", Prompt: l("Направление:", "Destination:")},
				{Key: "start", Prompt: l("Дата начала (YYYY-MM-DD):", "Start date (YYYY-MM-DD):")},
				{Key: "end", Prompt: l("Дата окончания (YYYY-MM-DD):", "End date (YYYY-MM-DD):")},
				{Key: "purpose", Prompt: l("Цель / мероприятие:", "Purpose / conference:")},
			},
			OnSubmit: submitBusinessTrip,
		},
		domain.ActionVacationRequest: {
			Intro: l("Запрос отпуска.", "Vacation request."),
			Fields: []FormField{
				{Key: "start", Prompt: l("Дата начала (YYYY-MM-DD):", "Start date (YYYY-MM-DD):")},
				{Key: "end", Prompt: l("Дата окончания (YYYY-MM-DD):", "End date (YYYY-MM-DD):")},
				{Key: "type", Prompt: l("Тип (paid/unpaid):", "Type (paid/unpaid):")},
			},
			OnSubmit: submitVacationRequest,
		},
		domain.ActionCertificateRequest: {
			Intro: l("Запрос справки в HR.", "Request HR certificate."),
			Fields: []FormField{
				{Key: "type", Prompt: l("Тип (employment/income/custom):", "Type (employment/income/custom):")},
			},
			OnSubmit: submitCertificateRequest,
		},
		domain.ActionBookAdmissionEvent: {
			Intro: l("Запись на мероприятие приёмной комиссии (День открытых дверей, Экскурсия по кампусу).", "Book your seat for an admission event (Open Day, Campus Tour)."),
			Fields: []FormField{
				{Key: "event_id", Prompt: l("ID мероприятия:", "Event ID:")},
				{Key: "name", Prompt: l("Ваше имя:", "Your name:")},
				{Key: "email", Prompt: l("Email:", "Email:")},
				{Key: "phone", Prompt: l("Телефон:", "Phone:"), Optional: true},
				{Key: "note", Prompt: l("Примечание (опционально):", "Note (optional):"), Optional: true},
			},
			OnSubmit: submitAdmissionEventBooking,
		},
	}
}

func (s *Service) startForm(ctx context.Context, sess *domain.Session, action domain.ActionID, def FormDefinition) error {
	sess.PendingAction = &domain.PendingAction{
		ID:   action,
		Step: 0,
		Data: map[string]string{},
	}
	s.saveSession(sess)
	intro := ""
	if def.Intro != nil {
		intro = def.Intro[domain.LanguageRU]
		if sess.Language == domain.LanguageEN && def.Intro[domain.LanguageEN] != "" {
			intro = def.Intro[domain.LanguageEN]
		}
	}
	prompt := def.Fields[0].PromptText(sess.Language)
	text := strings.TrimSpace(fmt.Sprintf("%s\n%s", intro, prompt))
	return s.reply(ctx, sess, text)
}

func (s *Service) handleFormInput(ctx context.Context, sess *domain.Session, input string) error {
	pa := sess.PendingAction
	if pa == nil {
		return nil
	}
	def, ok := s.forms[pa.ID]
	if !ok {
		sess.PendingAction = nil
		s.saveSession(sess)
		return s.reply(ctx, sess, s.t(sess.Language, "Форма недоступна.", "Form is no longer available."))
	}
	field := def.Fields[pa.Step]
	if input == "" && !field.Optional {
		return s.reply(ctx, sess, s.t(sess.Language, "Поле не может быть пустым.", "This field cannot be empty."))
	}
	if input == "" && field.Optional {
		pa.Data[field.Key] = ""
	} else {
		pa.Data[field.Key] = input
	}
	pa.Step++
	if pa.Step >= len(def.Fields) {
		sess.PendingAction = nil
		s.saveSession(sess)
		msg, err := def.OnSubmit(ctx, s, sess, pa.Data)
		if err != nil {
			return s.reply(ctx, sess, s.t(sess.Language, "Не удалось обработать форму. Попробуйте позже.", "Failed to submit form, please try again later."))
		}
		return s.replyMessage(ctx, sess, msg)
	}
	s.saveSession(sess)
	nextPrompt := def.Fields[pa.Step].PromptText(sess.Language)
	return s.reply(ctx, sess, nextPrompt)
}

// --- Form submit helpers ---

func submitEventRegistration(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	eventID, err := strconv.ParseInt(strings.TrimSpace(data["event_id"]), 10, 64)
	if err != nil {
		return messageError(sess.Language, "Неверный ID события.", "Invalid event ID."), nil
	}
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Авторизуйтесь, чтобы регистрироваться на события.", "Please authenticate to register for events."), nil
	}
	mode := data["mode"]
	if mode == "" {
		mode = "attendee"
	}
	status, err := s.backend.RSVPEvent(ctx, eventID, sess.Profile.ID, mode, data["note"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	switch status {
	case "registered", "updated":
		return messageSuccess(sess.Language, "Вы записаны на событие!", "Registration confirmed!"), nil
	case "already_registered":
		return messageError(sess.Language, "Вы уже зарегистрированы на это событие.", "You are already registered for this event."), nil
	default:
		return messageSuccess(sess.Language, "Статус регистрации: "+status, "Registration status: "+status), nil
	}
}

func submitLibrarySearch(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	results, err := s.backend.SearchBooks(ctx, data["query"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(results) == 0 {
		return messageSuccess(sess.Language, "Ничего не найдено.", "No books found."), nil
	}
	lines := []string{s.t(sess.Language, "Найдено:", "Results:")}
	for i, book := range results {
		if i >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("%d. %s — %s (ID: %d)", i+1, book.Title, book.Author, book.ID))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func submitLibraryReserve(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Сначала войдите в систему.", "Please login first."), nil
	}
	bookID, err := strconv.ParseInt(strings.TrimSpace(data["book_id"]), 10, 64)
	if err != nil {
		return messageError(sess.Language, "Неверный формат ID книги.", "Invalid book ID."), nil
	}
	if _, err := s.backend.ReserveBook(ctx, bookID, sess.Profile.ID); err != nil {
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language, "Запрос на резерв передан библиотеке.", "Reservation submitted to the library."), nil
}

func submitDormMaintenance(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Авторизуйтесь как студент.", "Please login as a student."), nil
	}
	reqID, err := s.backend.CreateDormMaintenance(ctx, sess.Profile.ID, data["type"], data["details"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language, fmt.Sprintf("Заявка #%d создана.", reqID), fmt.Sprintf("Request #%d created.", reqID)), nil
}

func submitSupportTicket(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	var userID *int64
	if sess.Profile != nil && sess.Profile.ID != 0 {
		userID = &sess.Profile.ID
	}
	id, err := s.backend.SubmitSupportTicket(ctx, data["category"], data["subject"], data["description"], userID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language, fmt.Sprintf("Заявка #%d зарегистрирована.", id), fmt.Sprintf("Ticket #%d created.", id)), nil
}

func submitBugReport(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	var userID *int64
	if sess.Profile != nil && sess.Profile.ID != 0 {
		userID = &sess.Profile.ID
	}
	id, err := s.backend.SubmitSupportTicket(ctx, "bug", data["subject"], data["description"], userID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language, fmt.Sprintf("Спасибо! Тикет #%d открыт.", id), fmt.Sprintf("Thanks! Ticket #%d opened.", id)), nil
}

func submitFAQQuery(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	var userID *int64
	if sess.Profile != nil && sess.Profile.ID != 0 {
		userID = &sess.Profile.ID
	}
	answer, err := s.backend.SubmitSupportQuery(ctx, userID, data["question"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return domain.OutgoingMessage{Text: answer}, nil
}

func submitAIQuery(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	answer, err := s.backend.RunAIQuery(ctx, data["question"], nil)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return domain.OutgoingMessage{Text: answer}, nil
}

func submitAISummary(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	summary, err := s.backend.CreateAISummary(ctx, data["text"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return domain.OutgoingMessage{Text: summary}, nil
}

func submitAIQuiz(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	questions, err := s.backend.GenerateAIQuiz(ctx, data["topic"], nil)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(questions) == 0 {
		return messageError(sess.Language, "Не удалось построить вопросы.", "Could not generate questions."), nil
	}
	lines := []string{s.t(sess.Language, "Ваш мини-квиз:", "Your mini quiz:")}
	for i, q := range questions {
		if i >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("%d) %s", i+1, q.Question))
		if len(q.Options) > 0 {
			for idx, opt := range q.Options {
				lines = append(lines, fmt.Sprintf("   %c) %s", 'A'+idx, opt))
			}
		}
		lines = append(lines, "")
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func submitAITranscription(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	transcript, err := s.backend.TranscribeAudio(ctx, data["audio"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return domain.OutgoingMessage{Text: transcript}, nil
}

func submitAdvisorChat(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	var userID *int64
	if sess.Profile != nil && sess.Profile.ID != 0 {
		userID = &sess.Profile.ID
	}
	response, err := s.backend.AdvisorChat(ctx, userID, data["topic"], data["prompt"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return domain.OutgoingMessage{Text: response}, nil
}

func submitBusinessTrip(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Авторизация обязательна.", "Authentication required."), nil
	}
	startISO, err := normalizeDate(data["start"])
	if err != nil {
		return messageError(sess.Language, "Введите дату в формате YYYY-MM-DD.", "Use YYYY-MM-DD date format."), nil
	}
	endISO, err := normalizeDate(data["end"])
	if err != nil {
		return messageError(sess.Language, "Введите дату в формате YYYY-MM-DD.", "Use YYYY-MM-DD date format."), nil
	}
	id, err := s.backend.RequestBusinessTrip(ctx, sess.Profile.ID, data["destination"], startISO, endISO, data["purpose"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language, fmt.Sprintf("Заявка #%d создана.", id), fmt.Sprintf("Request #%d submitted.", id)), nil
}

func submitVacationRequest(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Авторизация обязательна.", "Authentication required."), nil
	}
	startISO, err := normalizeDate(data["start"])
	if err != nil {
		return messageError(sess.Language, "Неверный формат даты.", "Invalid date format."), nil
	}
	endISO, err := normalizeDate(data["end"])
	if err != nil {
		return messageError(sess.Language, "Неверный формат даты.", "Invalid date format."), nil
	}
	id, err := s.backend.RequestVacation(ctx, sess.Profile.ID, startISO, endISO, data["type"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language, fmt.Sprintf("Заявка #%d создана.", id), fmt.Sprintf("Request #%d submitted.", id)), nil
}

func submitCertificateRequest(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Авторизуйтесь, чтобы запросить справку.", "Login to request a certificate."), nil
	}
	id, err := s.backend.RequestCertificate(ctx, sess.Profile.ID, data["type"])
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language, fmt.Sprintf("Запрос #%d отправлен в HR.", id), fmt.Sprintf("Request #%d sent to HR.", id)), nil
}

func submitAdmissionEventBooking(ctx context.Context, s *Service, sess *domain.Session, data map[string]string) (domain.OutgoingMessage, error) {
	eventID, err := strconv.ParseInt(strings.TrimSpace(data["event_id"]), 10, 64)
	if err != nil {
		return messageError(sess.Language, "Неверный ID мероприятия.", "Invalid event ID."), nil
	}
	name := strings.TrimSpace(data["name"])
	email := strings.TrimSpace(data["email"])
	if name == "" || email == "" {
		return messageError(sess.Language, "Имя и email обязательны.", "Name and email are required."), nil
	}
	phone := strings.TrimSpace(data["phone"])
	note := strings.TrimSpace(data["note"])

	bookingID, err := s.backend.BookAdmissionEvent(ctx, eventID, name, email, phone, note)
	if err != nil {
		if strings.Contains(err.Error(), "fully booked") {
			return messageError(sess.Language, "Мероприятие полностью забронировано.", "Event is fully booked."), nil
		}
		if strings.Contains(err.Error(), "already booked") {
			return messageError(sess.Language, "Вы уже забронировали это мероприятие.", "You have already booked this event."), nil
		}
		return domain.OutgoingMessage{}, err
	}
	return messageSuccess(sess.Language,
		fmt.Sprintf("✅ Ваше место забронировано! ID брони: %d\n\nМы отправили подтверждение на %s", bookingID, email),
		fmt.Sprintf("✅ Your seat is booked! Booking ID: %d\n\nConfirmation sent to %s", bookingID, email)), nil
}

// --- helpers ---

func normalizeDate(val string) (string, error) {
	val = strings.TrimSpace(val)
	if val == "" {
		return "", fmt.Errorf("empty date")
	}
	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC3339), nil
}

func messageSuccess(lang domain.Language, ru, en string) domain.OutgoingMessage {
	if lang == domain.LanguageEN {
		return domain.OutgoingMessage{Text: en}
	}
	return domain.OutgoingMessage{Text: ru}
}

func messageError(lang domain.Language, ru, en string) domain.OutgoingMessage {
	return messageSuccess(lang, ru, en)
}
