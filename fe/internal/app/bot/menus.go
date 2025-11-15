package bot

import (
	"github.com/escalopa/inno-vkode/internal/domain"
)

type MenuNode struct {
	ID          string
	ParentID    string
	Title       map[domain.Language]string
	Description map[domain.Language]string
	Action      domain.ActionID
	Children    []*MenuNode
}

func (n *MenuNode) TitleText(lang domain.Language) string {
	if n == nil {
		return ""
	}
	if text, ok := n.Title[lang]; ok && text != "" {
		return text
	}
	return n.Title[domain.LanguageRU]
}

func (n *MenuNode) DescriptionText(lang domain.Language) string {
	if n == nil || n.Description == nil {
		return ""
	}
	if text, ok := n.Description[lang]; ok && text != "" {
		return text
	}
	return n.Description[domain.LanguageRU]
}

type MenuRegistry struct {
	nodes map[string]*MenuNode
	roots map[domain.Role]*MenuNode
}

func buildMenuRegistry() *MenuRegistry {
	reg := &MenuRegistry{
		nodes: make(map[string]*MenuNode),
		roots: make(map[domain.Role]*MenuNode),
	}
	reg.registerRoot(domain.RoleApplicant, applicantMenu())
	reg.registerRoot(domain.RoleStudent, studentMenu())
	reg.registerRoot(domain.RoleEmployee, employeeMenu())
	reg.registerRoot(domain.RoleLeadership, leadershipMenu())
	return reg
}

func (r *MenuRegistry) registerRoot(role domain.Role, root *MenuNode) {
	if root == nil {
		return
	}
	r.walk(root, "")
	r.roots[role] = root
}

func (r *MenuRegistry) walk(node *MenuNode, parentID string) {
	if node == nil {
		return
	}
	node.ParentID = parentID
	r.nodes[node.ID] = node
	for _, child := range node.Children {
		r.walk(child, node.ID)
	}
}

func (r *MenuRegistry) Root(role domain.Role) *MenuNode {
	return r.roots[role]
}

func (r *MenuRegistry) Node(id string) *MenuNode {
	return r.nodes[id]
}

func applicantMenu() *MenuNode {
	return menuNode("applicant.root", l("🏠 Главное меню", "🏠 Main menu"), l("🎓 Гостевой режим для абитуриентов и гостей университета.", "🎓 Guest mode for applicants and university guests."), "", []*MenuNode{
		menuNode("applicant.admission", l("📚 Поступление", "📚 Admission"), l("Информация о программах и мероприятиях.", "Programs, open days and campus tours."), "", []*MenuNode{
			actionNode("applicant.admission.programs", l("ℹ️ О программах", "ℹ️ Programs & faculties"), domain.ActionViewAdmissionsPrograms),
			actionNode("applicant.admission.open_day", l("📅 День открытых дверей", "📅 Open day info"), domain.ActionBookOpenDay),
			actionNode("applicant.admission.campus_tour", l("🏛️ Тур по кампусу", "🏛️ Campus tour info"), domain.ActionBookCampusTour),
			actionNode("applicant.admission.book", l("✅ Забронировать место", "✅ Book event seat"), domain.ActionBookAdmissionEvent),
			actionNode("applicant.admission.contact", l("📞 Контакты приёмной", "📞 Contact admissions"), domain.ActionAdmissionsContact),
		}),
		menuNode("applicant.documents", l("📄 Документы", "📄 Documents"), nil, "", []*MenuNode{
			actionNode("applicant.documents.list", l("Требования", "Requirements list"), domain.ActionAdmissionDocuments),
			actionNode("applicant.documents.appointment", l("Записаться на подачу", "Book submission slot"), domain.ActionAdmissionAppointment),
		}),
		actionNode("applicant.language", l("🌐 Язык", "🌐 Language"), domain.ActionSwitchLanguage),
	})
}

func studentMenu() *MenuNode {
	return menuNode("student.root", l("🏠 Главное меню", "🏠 Main menu"), l("🎓 Персонализированные сервисы для студентов университета.", "🎓 Personalized services for university students."), "", []*MenuNode{
		menuNode("student.education", l("📚 Обучение", "📚 Education"), nil, "", []*MenuNode{
			actionNode("student.education.schedule", l("📅 Расписание", "📅 Schedule"), domain.ActionViewSchedule),
			actionNode("student.education.exams", l("🧪 Экзамены", "🧪 Exams"), domain.ActionViewExams),
			actionNode("student.education.grades", l("📊 Оценки", "📊 Grades"), domain.ActionViewGrades),
			actionNode("student.education.deadlines", l("⏰ Дедлайны", "⏰ Deadlines"), domain.ActionViewDeadlines),
			actionNode("student.education.feedback", l("💬 Отзывы преподавателям", "💬 Teacher feedback"), domain.ActionTeacherFeedback),
			actionNode("student.education.electives", l("➕ Запись на элективы", "➕ Elective registration"), domain.ActionElectiveRegistration),
		}),
		menuNode("student.projects", l("🚀 Проекты", "🚀 Projects"), nil, "", []*MenuNode{
			actionNode("student.projects.submit", l("💡 Подать проект", "💡 Submit project"), domain.ActionSubmitProject),
			actionNode("student.projects.team", l("👥 Команда", "👥 Build team"), domain.ActionBuildTeam),
			actionNode("student.projects.browse", l("🔍 Найти проект", "🔍 Browse projects"), domain.ActionBrowseProjects),
			actionNode("student.projects.mine", l("📋 Мои проекты", "📋 My projects"), domain.ActionMyProjects),
		}),
		menuNode("student.career", l("💼 Карьера", "💼 Career"), nil, "", []*MenuNode{
			actionNode("student.career.consult", l("📞 Консультация", "📞 Career consultation"), domain.ActionCareerConsultation),
			actionNode("student.career.jobs", l("💼 Вакансии", "💼 Job board"), domain.ActionBrowseJobs),
			actionNode("student.career.apply", l("✅ Откликнуться", "✅ Apply for job"), domain.ActionApplyJob),
			actionNode("student.career.my", l("📄 Мои заявки", "📄 My applications"), domain.ActionMyApplications),
		}),
		menuNode("student.dean", l("🏛️ Деканат", "🏛️ Dean's office"), nil, "", []*MenuNode{
			actionNode("student.dean.cert", l("📄 Справки", "📄 Certificates"), domain.ActionDeanCertificates),
			actionNode("student.dean.tuition", l("💳 Оплата обучения", "💳 Tuition payment"), domain.ActionDeanTuition),
			actionNode("student.dean.compensation", l("💵 Компенсации", "💵 Compensation"), domain.ActionDeanCompensation),
			actionNode("student.dean.appointment", l("📅 Приём", "📅 Appointments"), domain.ActionDeanAppointment),
			actionNode("student.dean.applications", l("📝 Заявления", "📝 Applications"), domain.ActionDeanApplications),
		}),
		menuNode("student.dorm", l("🏠 Общежитие", "🏠 Dormitory"), nil, "", []*MenuNode{
			actionNode("student.dorm.payment", l("💰 Оплата", "💰 Payment"), domain.ActionDormPayment),
			actionNode("student.dorm.services", l("🛎️ Сервисы", "🛎️ Services"), domain.ActionDormServices),
			actionNode("student.dorm.guests", l("🎫 Пропуск гостя", "🎫 Guest pass"), domain.ActionDormGuestPass),
			actionNode("student.dorm.maintenance", l("🔧 Заявка", "🔧 Maintenance"), domain.ActionDormMaintenance),
		}),
		menuNode("student.events", l("🎭 События", "🎭 Events"), nil, "", []*MenuNode{
			actionNode("student.events.calendar", l("📅 Календарь", "📅 Calendar"), domain.ActionEventsCalendar),
			actionNode("student.events.register", l("✅ Регистрация", "✅ Register"), domain.ActionEventsRegister),
			actionNode("student.events.my", l("📋 Мои события", "📋 My events"), domain.ActionEventsMine),
		}),
		menuNode("student.library", l("📚 Библиотека", "📚 Library"), nil, "", []*MenuNode{
			actionNode("student.library.search", l("🔍 Поиск книг", "🔍 Search books"), domain.ActionLibrarySearch),
			actionNode("student.library.reserve", l("📖 Резерв", "📖 Reserve book"), domain.ActionLibraryReserve),
			actionNode("student.library.my", l("📋 Мои книги", "📋 My library"), domain.ActionLibraryMy),
		}),
		menuNode("student.visa", l("🛂 Виза", "🛂 Visa services"), nil, "", []*MenuNode{
			actionNode("student.visa.status", l("📋 Статус", "📋 Status"), domain.ActionVisaStatus),
			actionNode("student.visa.renewal", l("🔄 Продление", "🔄 Renewal"), domain.ActionVisaRenewal),
			actionNode("student.visa.appointment", l("📅 Запись", "📅 Appointment"), domain.ActionVisaAppointment),
		}),
		menuNode("student.ai", l("🤖 Учебный ассистент", "🤖 AI Assistant"), nil, "", []*MenuNode{
			actionNode("student.ai.query", l("🔎 Вопрос RAG", "🔎 Knowledge query"), domain.ActionAIQuery),
			actionNode("student.ai.summary", l("📝 Конспект", "📝 Summarize text"), domain.ActionAISummary),
			actionNode("student.ai.quiz", l("❓ Генерация квиза", "❓ Generate quiz"), domain.ActionAIQuiz),
			actionNode("student.ai.transcribe", l("🎧 Транскрибация", "🎧 Transcription"), domain.ActionAITranscription),
			actionNode("student.ai.advisor", l("🧑‍🏫 Эдвайзер", "🧑‍🏫 Advisor chat"), domain.ActionAdvisorChat),
		}),
		menuNode("student.settings", l("⚙️ Настройки", "⚙️ Settings"), nil, "", []*MenuNode{
			actionNode("student.settings.profile", l("👤 Профиль", "👤 Profile"), domain.ActionViewProfile),
			actionNode("student.settings.language", l("🌐 Язык", "🌐 Language"), domain.ActionSwitchLanguage),
			actionNode("student.settings.notifications", l("🔔 Уведомления", "🔔 Notifications"), domain.ActionToggleNotifications),
		}),
		menuNode("student.support", l("ℹ️ Поддержка", "ℹ️ Support"), nil, "", []*MenuNode{
			actionNode("student.support.faq", l("❓ FAQ / AI", "❓ FAQ / AI"), domain.ActionFAQ),
			actionNode("student.support.contact", l("📨 Обратиться", "📨 Contact support"), domain.ActionContactSupport),
			actionNode("student.support.report", l("🐞 Сообщить об ошибке", "🐞 Report issue"), domain.ActionReportIssue),
		}),
	})
}

func employeeMenu() *MenuNode {
	return menuNode("employee.root", l("🏠 Главное меню", "🏠 Main menu"), l("💼 Профессиональные сервисы для сотрудников университета.", "💼 Professional services for university employees."), "", []*MenuNode{
		menuNode("employee.trips", l("✈️ Командировки", "✈️ Business trips"), nil, "", []*MenuNode{
			actionNode("employee.trips.list", l("📋 Мои заявки", "📋 My requests"), domain.ActionBusinessTripsList),
			actionNode("employee.trips.request", l("➕ Новая заявка", "➕ Submit request"), domain.ActionBusinessTripRequest),
		}),
		menuNode("employee.vacations", l("🏖️ Отпуска", "🏖️ Vacations"), nil, "", []*MenuNode{
			actionNode("employee.vacations.list", l("📊 Баланс и статусы", "📊 Balance & status"), domain.ActionVacationsList),
			actionNode("employee.vacations.request", l("➕ Новый запрос", "➕ Request vacation"), domain.ActionVacationRequest),
		}),
		menuNode("employee.office", l("🏢 Офисные сервисы", "🏢 Office services"), nil, "", []*MenuNode{
			actionNode("employee.office.certificates", l("📄 Справки", "📄 Certificates"), domain.ActionCertificatesList),
			actionNode("employee.office.certificate_request", l("📝 Запрос справки", "📝 Request certificate"), domain.ActionCertificateRequest),
			actionNode("employee.office.guest", l("🎫 Пропуск гостя", "🎫 Guest pass"), domain.ActionOfficeGuestPass),
			actionNode("employee.office.hr", l("📅 Визит в HR", "📅 HR appointment"), domain.ActionHRAppointment),
		}),
		menuNode("employee.events", l("🎭 События", "🎭 Events"), nil, "", []*MenuNode{
			actionNode("employee.events.calendar", l("📅 Календарь", "📅 Calendar"), domain.ActionEventsCalendar),
			actionNode("employee.events.register", l("✅ Регистрация", "✅ Register"), domain.ActionEventsRegister),
			actionNode("employee.events.my", l("📋 Мои события", "📋 My events"), domain.ActionEventsMine),
		}),
		menuNode("employee.visa", l("🛂 Виза", "🛂 Visa services"), nil, "", []*MenuNode{
			actionNode("employee.visa.status", l("📋 Статус", "📋 Status"), domain.ActionVisaStatus),
			actionNode("employee.visa.renewal", l("🔄 Продление", "🔄 Renewal"), domain.ActionVisaRenewal),
			actionNode("employee.visa.appointment", l("📅 Запись", "📅 Appointment"), domain.ActionVisaAppointment),
		}),
		menuNode("employee.ai", l("🤖 Ассистент", "🤖 AI assistant"), nil, "", []*MenuNode{
			actionNode("employee.ai.query", l("🔎 Вопрос", "🔎 Knowledge query"), domain.ActionAIQuery),
			actionNode("employee.ai.summary", l("📝 Конспект", "📝 Summary"), domain.ActionAISummary),
			actionNode("employee.ai.quiz", l("❓ Квиз", "❓ Quiz"), domain.ActionAIQuiz),
			actionNode("employee.ai.transcribe", l("🎧 Транскрибация", "🎧 Transcription"), domain.ActionAITranscription),
			actionNode("employee.ai.advisor", l("🧑‍🏫 Советник", "🧑‍🏫 Advisor"), domain.ActionAdvisorChat),
		}),
		menuNode("employee.settings", l("⚙️ Настройки", "⚙️ Settings"), nil, "", []*MenuNode{
			actionNode("employee.settings.profile", l("👤 Профиль", "👤 Profile"), domain.ActionViewProfile),
			actionNode("employee.settings.language", l("🌐 Язык", "🌐 Language"), domain.ActionSwitchLanguage),
			actionNode("employee.settings.notifications", l("🔔 Уведомления", "🔔 Notifications"), domain.ActionToggleNotifications),
		}),
		menuNode("employee.support", l("ℹ️ Поддержка", "ℹ️ Support"), nil, "", []*MenuNode{
			actionNode("employee.support.faq", l("❓ FAQ / AI", "❓ FAQ / AI"), domain.ActionFAQ),
			actionNode("employee.support.contact", l("📨 Обратиться", "📨 Contact support"), domain.ActionContactSupport),
			actionNode("employee.support.report", l("🐞 Сообщить об ошибке", "🐞 Report issue"), domain.ActionReportIssue),
		}),
	})
}

func leadershipMenu() *MenuNode {
	return menuNode("leadership.root", l("🏠 Главное меню", "🏠 Main menu"), l("👔 Инструменты и аналитика для руководителей университета.", "👔 Tools and analytics for university leadership."), "", []*MenuNode{
		menuNode("leadership.news", l("📰 Новости", "📰 News feed"), nil, "", []*MenuNode{
			actionNode("leadership.news.feed", l("📊 Лента упоминаний", "📊 Mentions feed"), domain.ActionLeadershipNews),
			actionNode("leadership.news.alerts", l("🔔 Оповещения", "🔔 Alerts"), domain.ActionLeadershipAlerts),
		}),
		menuNode("leadership.events", l("🎭 Мероприятия", "🎭 Events"), nil, "", []*MenuNode{
			actionNode("leadership.events.calendar", l("📅 Календарь", "📅 Calendar"), domain.ActionLeadershipEvents),
			actionNode("leadership.events.register", l("✅ Регистрация", "✅ Register"), domain.ActionEventsRegister),
		}),
		menuNode("leadership.ai", l("🤖 Аналитика", "🤖 AI insights"), nil, "", []*MenuNode{
			actionNode("leadership.ai.query", l("🔎 Запрос RAG", "🔎 Knowledge query"), domain.ActionAIQuery),
			actionNode("leadership.ai.summary", l("📝 Executive summary", "📝 Executive summary"), domain.ActionAISummary),
			actionNode("leadership.ai.transcribe", l("🎧 Транскрибация", "🎧 Transcription"), domain.ActionAITranscription),
		}),
		menuNode("leadership.settings", l("⚙️ Настройки", "⚙️ Settings"), nil, "", []*MenuNode{
			actionNode("leadership.settings.profile", l("👤 Профиль", "👤 Profile"), domain.ActionViewProfile),
			actionNode("leadership.settings.language", l("🌐 Язык", "🌐 Language"), domain.ActionSwitchLanguage),
			actionNode("leadership.settings.notifications", l("🔔 Уведомления", "🔔 Notifications"), domain.ActionToggleNotifications),
		}),
		menuNode("leadership.support", l("ℹ️ Поддержка", "ℹ️ Support"), nil, "", []*MenuNode{
			actionNode("leadership.support.contact", l("📨 Обратиться", "📨 Contact support"), domain.ActionContactSupport),
			actionNode("leadership.support.report", l("🐞 Сообщить об ошибке", "🐞 Report issue"), domain.ActionReportIssue),
		}),
	})
}

func menuNode(id string, title, desc map[domain.Language]string, action domain.ActionID, children []*MenuNode) *MenuNode {
	return &MenuNode{
		ID:          id,
		Title:       title,
		Description: desc,
		Action:      action,
		Children:    children,
	}
}

func actionNode(id string, title map[domain.Language]string, action domain.ActionID) *MenuNode {
	return menuNode(id, title, nil, action, nil)
}

func l(ru, en string) map[domain.Language]string {
	return map[domain.Language]string{
		domain.LanguageRU: ru,
		domain.LanguageEN: en,
	}
}
