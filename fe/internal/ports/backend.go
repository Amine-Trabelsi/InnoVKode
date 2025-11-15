package ports

import (
	"context"

	"github.com/escalopa/inno-vkode/internal/domain"
)

type Backend interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserProfile, error)

	GetSchedule(ctx context.Context, userID int64) ([]domain.ScheduleEntry, error)
	GetCourses(ctx context.Context, userID int64) ([]domain.Course, error)
	GetExams(ctx context.Context, userID int64) ([]domain.ExamEntry, error)
	GetGrades(ctx context.Context, userID int64) ([]domain.GradeRecord, error)
	GetDeadlines(ctx context.Context, userID int64) ([]domain.Deadline, error)

	ListEvents(ctx context.Context) ([]domain.Event, error)
	ListNews(ctx context.Context) ([]domain.NewsItem, error)
	ListClubs(ctx context.Context) ([]domain.Club, error)
	RSVPEvent(ctx context.Context, eventID int64, userID int64, registrationType string, note string) (string, error)
	CancelRSVP(ctx context.Context, eventID int64, userID int64) error
	ListUserEvents(ctx context.Context, userID int64) ([]domain.Event, error)

	ListAdmissionsPrograms(ctx context.Context) ([]domain.AdmissionProgram, error)
	ListAdmissionEvents(ctx context.Context) ([]domain.AdmissionEvent, error)
	BookAdmissionEvent(ctx context.Context, eventID int64, applicantName, email, phone, note string) (int64, error)
	ListAdmissionEventBookings(ctx context.Context, eventID int64) ([]domain.AdmissionEventBooking, error)
	SubmitAdmissionApplication(ctx context.Context, applicantName, email string, programID *int64, details map[string]any) (int64, error)
	UploadAdmissionDocument(ctx context.Context, applicationID int64, fileName, fileType, storageURL string) (int64, error)
	AskAdmissionQuestion(ctx context.Context, question string) (string, error)

	CreateDeanRequest(ctx context.Context, userID int64, requestType string, payload map[string]any) (int64, error)

	GetDormRoom(ctx context.Context, studentID int64) (*domain.DormRoom, error)
	CreateDormMaintenance(ctx context.Context, studentID int64, requestType, description string) (int64, error)
	SubmitDormPayment(ctx context.Context, studentID int64, amount float64, reference string) (int64, error)

	SearchBooks(ctx context.Context, q string) ([]domain.LibraryBook, error)
	ReserveBook(ctx context.Context, bookID, studentID int64) (int64, error)
	ListBorrowedBooks(ctx context.Context, studentID int64) ([]domain.LibraryLoan, error)

	SubmitSupportTicket(ctx context.Context, category, subject, description string, userID *int64) (int64, error)
	SubmitSupportQuery(ctx context.Context, userID *int64, question string) (string, error)
	AdvisorChat(ctx context.Context, userID *int64, topic, prompt string) (string, error)
	RunAIQuery(ctx context.Context, question string, filters map[string]any) (string, error)
	CreateAISummary(ctx context.Context, text string) (string, error)
	GenerateAIQuiz(ctx context.Context, prompt string, courseID *int64) ([]domain.QuizQuestion, error)
	TranscribeAudio(ctx context.Context, audioRef string) (string, error)

	GetVacations(ctx context.Context, employeeID int64) ([]domain.VacationRequest, error)
	RequestVacation(ctx context.Context, employeeID int64, startISO, endISO, vacationType string) (int64, error)
	GetBusinessTrips(ctx context.Context, employeeID int64) ([]domain.BusinessTrip, error)
	RequestBusinessTrip(ctx context.Context, employeeID int64, destination, startISO, endISO, purpose string) (int64, error)
	GetCertificates(ctx context.Context, employeeID int64) ([]domain.HRLetter, error)
	RequestCertificate(ctx context.Context, employeeID int64, certificateType string) (int64, error)

	SendNotification(ctx context.Context, subject, body string, recipientID *int64) (int64, error)
}
