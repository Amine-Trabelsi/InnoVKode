package bot

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/escalopa/inno-vkode/internal/domain"
)

func (s *Service) handleAction(ctx context.Context, sess *domain.Session, action domain.ActionID) (domain.OutgoingMessage, error) {
	switch action {
	case domain.ActionViewAdmissionsPrograms:
		return s.handleAdmissionsOverview(ctx, sess)
	case domain.ActionBookOpenDay:
		return s.handleAdmissionEvent(ctx, sess, "open_day", "📅 Open Day"), nil
	case domain.ActionBookCampusTour:
		return s.handleAdmissionEvent(ctx, sess, "campus_tour", "🏛️ Campus tour"), nil
	case domain.ActionAdmissionsContact:
		return domain.OutgoingMessage{
			Text:      fmt.Sprintf("📞 **Приёмная комиссия**\n\n📧 Email: %s\n📱 Телефон: %s\n🏢 Офис: %s", s.cfg.AdmissionsEmail, s.cfg.AdmissionsPhone, s.cfg.AdmissionsOffice),
			ParseMode: domain.ParseModeMarkdown,
		}, nil
	case domain.ActionAdmissionDocuments:
		return domain.OutgoingMessage{
			Text:      "📄 **Необходимые документы для поступления:**\n\n• Копия паспорта\n• Предыдущий диплом + транскрипт\n• Результаты ЕГЭ\n• Мотивационное письмо\n• Фото 3x4 (6 шт.)\n\n⚠️ Загрузите PDF-файлы перед записью на подачу документов.",
			ParseMode: domain.ParseModeMarkdown,
		}, nil
	case domain.ActionAdmissionAppointment:
		return domain.OutgoingMessage{
			Text: "📅 Используйте форму записи для выбора даты/времени.\nПриносите оригиналы документов в кампус.\n\n🕐 Доступные слоты: Пн–Пт 10:00-17:00.",
		}, nil
	case domain.ActionViewSchedule:
		return s.handleSchedule(ctx, sess)
	case domain.ActionViewExams:
		return s.handleExams(ctx, sess)
	case domain.ActionViewGrades:
		return s.handleGrades(ctx, sess)
	case domain.ActionViewDeadlines:
		return s.handleDeadlines(ctx, sess)
	case domain.ActionTeacherFeedback:
		return domain.OutgoingMessage{
			Text: "Teacher feedback form:\n1. Open course in LMS\n2. Rate 1-5 ⭐️\n3. Add optional comment\n\nWe're syncing ratings nightly to dean's dashboard.",
		}, nil
	case domain.ActionElectiveRegistration:
		return domain.OutgoingMessage{
			Text: "Elective enrollment opens each semester via ISU portal. Browse catalog → add to cart → confirm by advisor.",
		}, nil
	case domain.ActionSubmitProject:
		return domain.OutgoingMessage{
			Text: "Project submission checklist:\n• Title & summary\n• Team composition\n• Skills needed\nSend details via advisor or innovation centre. We'll add a form in the next update.",
		}, nil
	case domain.ActionBuildTeam:
		return domain.OutgoingMessage{
			Text: "Team management: review incoming requests in LMS > Projects. Approve members, assign roles and publish needs in the chat.",
		}, nil
	case domain.ActionBrowseProjects:
		return domain.OutgoingMessage{
			Text: "Sample active projects:\n1. AI Campus Guide – looking for ML engineer.\n2. Green Dorms – needs UX/UI designer.\n3. Smart Attendance – needs backend Go dev.\nUse /support to reach coordinators.",
		}, nil
	case domain.ActionMyProjects:
		return domain.OutgoingMessage{
			Text: "We'll sync personal project dashboards soon. Meanwhile, track tasks in Notion board shared by your supervisor.",
		}, nil
	case domain.ActionCareerConsultation:
		return domain.OutgoingMessage{
			Text: "Career Center works Mon–Thu. Choose topic (CV, interview, job search) and book via /support → consultation.",
		}, nil
	case domain.ActionBrowseJobs:
		return domain.OutgoingMessage{
			Text: "Featured roles:\n• Data Analyst Intern — FinTech Lab — apply before 25 Nov.\n• Product Manager Assistant — Innovation Hub.\n• QA Engineer — Partner company.\nUse Apply button in career portal.",
		}, nil
	case domain.ActionApplyJob:
		return domain.OutgoingMessage{
			Text: "To apply, attach CV + motivation letter and reference job ID. HR responds within 5 days.",
		}, nil
	case domain.ActionMyApplications:
		return domain.OutgoingMessage{
			Text: "Application tracker:\n• Data Analyst Intern — Interview scheduled.\n• Product Manager Assistant — Under review.\nWe'll add live updates soon.",
		}, nil
	case domain.ActionDeanCertificates:
		return domain.OutgoingMessage{
			Text: "Certificate desk issues documents in 3 working days. Use /support to specify type (enrollment, scholarship, transcript).",
		}, nil
	case domain.ActionDeanTuition:
		return domain.OutgoingMessage{
			Text: fmt.Sprintf("Balance & payments available in student portal.\nOnline payment link: %s", s.cfg.TuitionPaymentURL),
		}, nil
	case domain.ActionDeanCompensation:
		return domain.OutgoingMessage{
			Text: "Compensation request flow:\n1. Choose program (transport, medical, tech purchase).\n2. Attach receipts (PDF/JPG up to 10 MB).\n3. Submit via /support → compensation.\nFinance reviews within 7 days.",
		}, nil
	case domain.ActionDeanAppointment:
		return domain.OutgoingMessage{
			Text: "Dean's office appointments available Tue/Thu. Provide topic (documents, transfer, leave) and preferred time when contacting support.",
		}, nil
	case domain.ActionDeanApplications:
		return domain.OutgoingMessage{
			Text: "Transfer / academic leave requests:\n• Fill template form\n• Add justification\n• Upload supporting docs\nSend via /support to initiate workflow.",
		}, nil
	case domain.ActionDormPayment:
		return s.handleDormPayment(ctx, sess)
	case domain.ActionDormServices:
		return domain.OutgoingMessage{
			Text: "Available services: laundry, cleaning, linen exchange. Order via dorm desk or /support specifying room & slot.",
		}, nil
	case domain.ActionDormGuestPass:
		return domain.OutgoingMessage{
			Text: "Guest pass steps:\n1. Send guest name + passport + visit hours via /support.\n2. Duty officer confirms by SMS.\n3. Collect printed pass at lobby.",
		}, nil
	case domain.ActionEventsCalendar:
		return s.handleEvents(ctx, sess)
	case domain.ActionEventsRegister:
		return s.handleEventRegistration(ctx, sess)
	case domain.ActionEventsMine:
		return s.handlePersonalEvents(ctx, sess)
	case domain.ActionLibraryMy:
		return s.handleLibraryLoans(ctx, sess)
	case domain.ActionVisaStatus:
		return s.handleVisaStatus(ctx, sess)
	case domain.ActionVisaMakeApplication:
		return s.handleVisaMakeApplication(ctx, sess)
	case domain.ActionViewProfile:
		return s.handleProfile(sess), nil
	case domain.ActionToggleNotifications:
		sess.NotificationsEnabled = !sess.NotificationsEnabled
		s.saveSession(sess)
		if sess.NotificationsEnabled {
			return domain.OutgoingMessage{Text: s.t(sess.Language, "🔔 Уведомления включены!", "🔔 Notifications enabled!")}, nil
		}
		return domain.OutgoingMessage{Text: s.t(sess.Language, "🔕 Уведомления отключены.", "🔕 Notifications disabled.")}, nil
	case domain.ActionFAQ:
		return domain.OutgoingMessage{Text: s.t(sess.Language, "❓ Опишите вопрос, и бот подскажет из базы знаний.\nНажмите кнопку ещё раз, чтобы заполнить форму.", "❓ Describe your question, then press the button again to fill the quick form.")}, nil
	case domain.ActionReportIssue:
		return domain.OutgoingMessage{Text: s.t(sess.Language, "🐞 Кратко опишите найденную ошибку и прикрепите скриншот через форму.", "🐞 Describe the issue and attach a screenshot via the form.")}, nil
	case domain.ActionLeadershipNews:
		return s.handleNews(ctx, sess)
	case domain.ActionLeadershipAlerts:
		return domain.OutgoingMessage{
			Text: "Alerts deliver daily digest of critical mentions. Enable notifications in ⚙️ Settings to receive push updates.",
		}, nil
	case domain.ActionLeadershipEvents:
		return s.handleEvents(ctx, sess)
	case domain.ActionBusinessTripsList:
		return s.handleBusinessTrips(ctx, sess)
	case domain.ActionVacationsList:
		return s.handleVacations(ctx, sess)
	case domain.ActionCertificatesList:
		return s.handleCertificates(ctx, sess)
	case domain.ActionOfficeGuestPass:
		return domain.OutgoingMessage{
			Text: "Office guest passes available for HQ buildings. Provide guest details via /support at least 1 day before visit.",
		}, nil
	case domain.ActionHRAppointment:
		return domain.OutgoingMessage{
			Text: "HR office bookings: Mon–Thu 14:00-17:00. Specify topic (documents, onboarding, policies) when opening a support ticket.",
		}, nil
	default:
		return domain.OutgoingMessage{Text: "🚀 Функция скоро появится. Следите за обновлениями!"}, nil
	}
}

func (s *Service) handleAdmissionsOverview(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	programs, err := s.backend.ListAdmissionsPrograms(ctx)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	events, err := s.backend.ListAdmissionEvents(ctx)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	lines := []string{"Programs:"}
	for i, p := range programs {
		if i >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("• %s (%s) — %d years, tuition %s₽", p.Title, p.Faculty, p.DurationYears, p.Tuition))
	}
	lines = append(lines, "", "Upcoming events:")
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateTime.Before(events[j].DateTime)
	})
	for i, ev := range events {
		if i >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("• %s — %s (%s)", ev.Title, ev.DateTime.Format("02 Jan 15:04"), ev.Location))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleAdmissionEvent(ctx context.Context, sess *domain.Session, eventType, title string) domain.OutgoingMessage {
	events, err := s.backend.ListAdmissionEvents(ctx)
	if err != nil {
		return domain.OutgoingMessage{Text: "Unable to load events, please try later."}
	}
	lines := []string{title + ":"}
	for _, ev := range events {
		if eventType != "" && !strings.Contains(strings.ToLower(ev.EventType), strings.ToLower(eventType)) {
			continue
		}
		capacity := ""
		if ev.MaxAttendees > 0 {
			available := ev.MaxAttendees - ev.CurrentAttendees
			if available <= 0 {
				capacity = " [ПОЛНОСТЬЮ ЗАБРОНИРОВАНО / FULLY BOOKED]"
			} else {
				capacity = fmt.Sprintf(" [Свободно мест / Available: %d/%d]", available, ev.MaxAttendees)
			}
		}
		lines = append(lines, fmt.Sprintf("• ID %d: %s — %s — %s%s", ev.ID, ev.Title, ev.DateTime.Format("02 Jan 15:04"), ev.Location, capacity))
	}
	if len(lines) == 1 {
		lines = append(lines, "Slots will be published soon. Stay tuned!")
	} else {
		lines = append(lines, "", s.t(sess.Language,
			"Для записи используйте меню 'Поступление' → 'Забронировать место'.",
			"To book, use 'Admission' menu → 'Book event seat'."))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}
}

func (s *Service) handleSchedule(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.GetSchedule(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: s.t(sess.Language, "Расписание занятий отсутствует на этой неделе.", "No sessions scheduled this week.")}, nil
	}
	// Build keyboard with unique days
	dayMap := make(map[string]string) // date -> display
	for _, item := range items {
		date := item.StartTime.Format("2006-01-02")
		display := item.StartTime.Format("Mon 02 Jan")
		dayMap[date] = display
	}
	kb := &domain.Keyboard{}
	for date, display := range dayMap {
		btn := domain.KeyboardButton{
			Label:   display,
			Style:   domain.ButtonStylePrimary,
			Kind:    domain.ButtonKindCallback,
			Payload: "schedule:" + date,
		}
		kb.Rows = append(kb.Rows, []domain.KeyboardButton{btn})
	}
	// Add "All" button
	kb.Rows = append(kb.Rows, []domain.KeyboardButton{{
		Label:   s.t(sess.Language, "Все", "All"),
		Style:   domain.ButtonStyleSecondary,
		Kind:    domain.ButtonKindCallback,
		Payload: "schedule:all",
	}})
	return domain.OutgoingMessage{
		Text:     s.t(sess.Language, "Выберите день для просмотра расписания:", "Select a day to view schedule:"),
		Keyboard: kb,
	}, nil
}

func (s *Service) handleExams(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.GetExams(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "No exams scheduled."}, nil
	}
	lines := []string{"Upcoming exams:"}
	for _, exam := range items {
		lines = append(lines, fmt.Sprintf("• %s — %s (%s)", exam.Date.Format("02 Jan 15:00"), exam.Title, exam.Room))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleGrades(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.GetGrades(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "No grades yet."}, nil
	}
	lines := []string{"Recent grades:"}
	for i, g := range items {
		if i >= 6 {
			break
		}
		lines = append(lines, fmt.Sprintf("• %s — %s (GPA %.2f)", g.Title, g.Grade, g.GPAPoints))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleDeadlines(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.GetDeadlines(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "No active deadlines."}, nil
	}
	lines := []string{"Deadlines:"}
	for _, d := range items {
		lines = append(lines, fmt.Sprintf("• %s — %s (%s)", d.DueDate.Format("02 Jan"), d.Title, d.Status))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleDormPayment(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	room, err := s.backend.GetDormRoom(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	return domain.OutgoingMessage{
		Text: fmt.Sprintf("Dorm room: %s (%s)\nBalance: %.2f₽\nOnline payment: %s", room.Room, room.Building, room.Balance, s.cfg.DormPaymentURL),
	}, nil
}

func (s *Service) handleEvents(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	events, err := s.backend.ListEvents(ctx)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(events) == 0 {
		return domain.OutgoingMessage{Text: "No events planned right now."}, nil
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateTime.Before(events[j].DateTime)
	})
	lines := []string{"Upcoming events:"}
	for i, ev := range events {
		if i >= 6 {
			break
		}
		free := ""
		if ev.MaxAttendees > 0 {
			free = fmt.Sprintf(" (%d/%d)", ev.CurrentAttendees, ev.MaxAttendees)
		}
		lines = append(lines, fmt.Sprintf("• %s%s — %s @ %s", ev.Title, free, ev.DateTime.Format("02 Jan 15:00"), ev.Location))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handlePersonalEvents(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	events, err := s.backend.ListUserEvents(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(events) == 0 {
		return domain.OutgoingMessage{Text: "You have no active registrations."}, nil
	}
	kb := &domain.Keyboard{}
	for _, ev := range events {
		mode := ev.UserRegistrationType
		if mode == "" {
			mode = "attendee"
		}
		label := fmt.Sprintf("❌ %s — %s (%s)", ev.Title, ev.DateTime.Format("02 Jan 15:00"), mode)
		btn := domain.KeyboardButton{
			Label:   label,
			Style:   domain.ButtonStyleSecondary,
			Kind:    domain.ButtonKindCallback,
			Payload: "cancel_event:" + strconv.FormatInt(ev.ID, 10),
		}
		kb.Rows = append(kb.Rows, []domain.KeyboardButton{btn})
	}
	return domain.OutgoingMessage{
		Text:    "Your registrations (click to cancel):",
		Keyboard: kb,
	}, nil
}

func (s *Service) handleEventRegistration(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	events, err := s.backend.ListEvents(ctx)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(events) == 0 {
		return domain.OutgoingMessage{Text: "No events available for registration."}, nil
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateTime.Before(events[j].DateTime)
	})
	kb := &domain.Keyboard{}
	for _, ev := range events {
		free := ""
		if ev.MaxAttendees > 0 {
			free = fmt.Sprintf(" (%d/%d)", ev.CurrentAttendees, ev.MaxAttendees)
		}
		label := fmt.Sprintf("%s%s — %s @ %s", ev.Title, free, ev.DateTime.Format("02 Jan 15:00"), ev.Location)
		btn := domain.KeyboardButton{
			Label:   label,
			Style:   domain.ButtonStylePrimary,
			Kind:    domain.ButtonKindCallback,
			Payload: "event_select:" + strconv.FormatInt(ev.ID, 10),
		}
		kb.Rows = append(kb.Rows, []domain.KeyboardButton{btn})
	}
	return domain.OutgoingMessage{
		Text:    "Select an event to register:",
		Keyboard: kb,
	}, nil
}

func (s *Service) handleLibraryLoans(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.ListBorrowedBooks(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "You have no borrowed books."}, nil
	}
	lines := []string{"Borrowed books:"}
	for _, loan := range items {
		lines = append(lines, fmt.Sprintf("• %s — due %s (%s)", loan.Title, loan.DueAt.Format("02 Jan"), loan.Status))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleNews(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	items, err := s.backend.ListNews(ctx)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "News feed is empty for now."}, nil
	}
	lines := []string{"Latest mentions:"}
	for i, item := range items {
		if i >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("• %s — %s", item.Title, item.PublishedAt.Format("02 Jan")))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleBusinessTrips(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.GetBusinessTrips(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "No trip requests yet."}, nil
	}
	lines := []string{"Trip requests:"}
	for _, trip := range items {
		lines = append(lines, fmt.Sprintf("• %s — %s to %s (%s)", trip.Purpose, trip.StartDate.Format("02 Jan"), trip.EndDate.Format("02 Jan"), strings.ToUpper(trip.Status)))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleVacations(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.GetVacations(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "No vacation requests yet."}, nil
	}
	lines := []string{"Vacation requests:"}
	for _, v := range items {
		lines = append(lines, fmt.Sprintf("• %s → %s (%s)", v.StartDate.Format("02 Jan"), v.EndDate.Format("02 Jan"), strings.ToUpper(v.Status)))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleCertificates(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	items, err := s.backend.GetCertificates(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(items) == 0 {
		return domain.OutgoingMessage{Text: "No certificate requests yet."}, nil
	}
	lines := []string{"Certificates:"}
	for _, c := range items {
		lines = append(lines, fmt.Sprintf("• %s — %s", c.CertificateType, strings.ToUpper(c.Status)))
	}
	return domain.OutgoingMessage{Text: strings.Join(lines, "\n")}, nil
}

func (s *Service) handleProfile(sess *domain.Session) domain.OutgoingMessage {
	if sess.Profile == nil {
		return domain.OutgoingMessage{Text: "👤 Гостевая сессия. Используйте /start для авторизации."}
	}
	p := sess.Profile
	text := fmt.Sprintf("👤 **Профиль пользователя**\n\n📝 Имя (RU): %s\n📝 Имя (EN): %s\n📧 Email: %s\n🏷️ Роль: %s\n🏫 Факультет: %s\n🏠 Общежитие: %s",
		emptyFallback(p.NameRU, "—"),
		emptyFallback(p.NameEN, "—"),
		p.Email,
		p.Role,
		emptyFallback(p.Faculty, "—"),
		emptyFallback(p.DormRoom, "—"),
	)
	return domain.OutgoingMessage{
		Text:      text,
		ParseMode: domain.ParseModeMarkdown,
	}
}

func (s *Service) handleVisaStatus(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	if sess.Profile == nil || sess.Profile.ID == 0 {
		return messageError(sess.Language, "Нужна авторизация.", "Please login first."), nil
	}
	apps, err := s.backend.GetVisaApplications(ctx, sess.Profile.ID)
	if err != nil {
		return domain.OutgoingMessage{}, err
	}
	if len(apps) == 0 {
		return domain.OutgoingMessage{Text: s.t(sess.Language, "У вас нет заявок на визу.", "You have no visa applications.")}, nil
	}
	kb := &domain.Keyboard{}
	for _, app := range apps {
		label := fmt.Sprintf("%s (%s)", app["application_type"], app["status"])
		btn := domain.KeyboardButton{
			Label:   label,
			Style:   domain.ButtonStylePrimary,
			Kind:    domain.ButtonKindCallback,
			Payload: "visa_app:" + fmt.Sprintf("%v", app["id"]),
		}
		kb.Rows = append(kb.Rows, []domain.KeyboardButton{btn})
	}
	return domain.OutgoingMessage{
		Text:     s.t(sess.Language, "Выберите заявку:", "Select an application:"),
		Keyboard: kb,
	}, nil
}

func (s *Service) handleVisaMakeApplication(ctx context.Context, sess *domain.Session) (domain.OutgoingMessage, error) {
	kb := &domain.Keyboard{
		Rows: [][]domain.KeyboardButton{
			{
				{Label: s.t(sess.Language, "Продление визы", "Visa renewal"), Kind: domain.ButtonKindCallback, Payload: "visa_type:visa_renewal", Style: domain.ButtonStylePrimary},
				{Label: s.t(sess.Language, "Продление регистрации", "Registration renewal"), Kind: domain.ButtonKindCallback, Payload: "visa_type:registration_renewal", Style: domain.ButtonStylePrimary},
			},
		},
	}
	return domain.OutgoingMessage{
		Text:     s.t(sess.Language, "Выберите тип заявки:", "Choose application type:"),
		Keyboard: kb,
	}, nil
}

func emptyFallback(val, fallback string) string {
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	return val
}
