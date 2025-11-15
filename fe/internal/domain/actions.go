package domain

type ActionID string

const (
	ActionSwitchLanguage         ActionID = "switch_language"
	ActionViewAdmissionsPrograms ActionID = "admissions_programs"
	ActionBookOpenDay            ActionID = "book_open_day"
	ActionBookCampusTour         ActionID = "book_campus_tour"
	ActionBookAdmissionEvent     ActionID = "book_admission_event"
	ActionAdmissionsContact      ActionID = "admissions_contact"
	ActionAdmissionDocuments     ActionID = "admissions_documents"
	ActionAdmissionAppointment   ActionID = "admissions_appointment"

	ActionViewSchedule          ActionID = "view_schedule"
	ActionViewExams             ActionID = "view_exams"
	ActionViewGrades            ActionID = "view_grades"
	ActionViewDeadlines         ActionID = "view_deadlines"
	ActionTeacherFeedback       ActionID = "teacher_feedback"
	ActionElectiveRegistration  ActionID = "elective_registration"

	ActionSubmitProject         ActionID = "submit_project"
	ActionBuildTeam             ActionID = "build_team"
	ActionBrowseProjects        ActionID = "browse_projects"
	ActionMyProjects            ActionID = "my_projects"

	ActionCareerConsultation    ActionID = "career_consultation"
	ActionBrowseJobs            ActionID = "browse_jobs"
	ActionApplyJob              ActionID = "apply_job"
	ActionMyApplications        ActionID = "my_applications"

	ActionDeanCertificates      ActionID = "dean_certificates"
	ActionDeanTuition           ActionID = "dean_tuition"
	ActionDeanCompensation      ActionID = "dean_compensation"
	ActionDeanAppointment       ActionID = "dean_appointment"
	ActionDeanApplications      ActionID = "dean_applications"

	ActionDormPayment           ActionID = "dorm_payment"
	ActionDormServices          ActionID = "dorm_services"
	ActionDormGuestPass         ActionID = "dorm_guest_pass"
	ActionDormMaintenance       ActionID = "dorm_maintenance"

	ActionEventsCalendar        ActionID = "events_calendar"
	ActionEventsRegister        ActionID = "events_register"
	ActionEventsMine            ActionID = "events_mine"

	ActionLibrarySearch         ActionID = "library_search"
	ActionLibraryReserve        ActionID = "library_reserve"
	ActionLibraryMy             ActionID = "library_my"

	ActionVisaStatus            ActionID = "visa_status"
	ActionVisaRenewal           ActionID = "visa_renewal"
	ActionVisaAppointment       ActionID = "visa_appointment"

	ActionViewProfile           ActionID = "view_profile"
	ActionToggleNotifications   ActionID = "toggle_notifications"
	ActionContactSupport        ActionID = "contact_support"
	ActionFAQ                   ActionID = "faq"
	ActionReportIssue           ActionID = "report_issue"

	ActionAIQuery               ActionID = "ai_query"
	ActionAISummary             ActionID = "ai_summary"
	ActionAIQuiz                ActionID = "ai_quiz"
	ActionAITranscription       ActionID = "ai_transcription"
	ActionAdvisorChat           ActionID = "advisor_chat"

	ActionBusinessTripsList     ActionID = "business_trips_list"
	ActionBusinessTripRequest   ActionID = "business_trip_request"
	ActionVacationsList         ActionID = "vacations_list"
	ActionVacationRequest       ActionID = "vacation_request"
	ActionCertificatesList      ActionID = "certificates_list"
	ActionCertificateRequest    ActionID = "certificate_request"
	ActionOfficeGuestPass       ActionID = "office_guest_pass"
	ActionHRAppointment         ActionID = "hr_appointment"

	ActionLeadershipNews        ActionID = "leadership_news"
	ActionLeadershipAlerts      ActionID = "leadership_alerts"
	ActionLeadershipEvents      ActionID = "leadership_events"
)
