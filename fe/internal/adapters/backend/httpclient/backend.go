package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/escalopa/inno-vkode/internal/domain"
	"github.com/escalopa/inno-vkode/internal/ports"
)

type Backend struct {
	baseURL string
	client  *http.Client
	log     zerolog.Logger
}

var _ ports.Backend = (*Backend)(nil)

func New(baseURL string, timeout time.Duration, log zerolog.Logger) *Backend {
	trimmed := strings.TrimRight(baseURL, "/")
	if trimmed == "" {
		trimmed = "http://localhost:8001"
	}
	return &Backend{
		baseURL: trimmed,
		client: &http.Client{
			Timeout: timeout,
		},
		log: log,
	}
}

func (b *Backend) buildURL(p string, query url.Values) string {
	full := b.baseURL
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	full += path.Clean(p)
	if query != nil && len(query) > 0 {
		full += "?" + query.Encode()
	}
	return full
}

func (b *Backend) doRequest(ctx context.Context, method, p string, query url.Values, payload any, out any) error {
	ctx, cancel := context.WithTimeout(ctx, b.client.Timeout)
	defer cancel()

	var body io.Reader
	if payload != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return fmt.Errorf("encode payload: %w", err)
		}
		body = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, b.buildURL(p, query), body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			b.log.Error().Err(cerr).Msg("failed to close backend response body")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend %s %s returned %d: %s", method, p, resp.StatusCode, string(raw))
	}

	if out == nil {
		return nil
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	b.log.Debug().Str("response_body", string(raw)).Msg("backend response")
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (b *Backend) get(ctx context.Context, p string, query url.Values, out any) error {
	return b.doRequest(ctx, http.MethodGet, p, query, nil, out)
}

func (b *Backend) post(ctx context.Context, p string, payload any, out any) error {
	return b.doRequest(ctx, http.MethodPost, p, nil, payload, out)
}

func (b *Backend) put(ctx context.Context, p string, payload any, out any) error {
	return b.doRequest(ctx, http.MethodPut, p, nil, payload, out)
}

// region Users & Profiles

func (b *Backend) GetUserByEmail(ctx context.Context, email string) (*domain.UserProfile, error) {
	q := url.Values{}
	q.Set("email", email)
	var result domain.UserProfile
	if err := b.get(ctx, "/api/v1/users/by-email", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// endregion

// region Academic

func (b *Backend) GetSchedule(ctx context.Context, userID int64) ([]domain.ScheduleEntry, error) {
	var result []domain.ScheduleEntry
	if err := b.get(ctx, fmt.Sprintf("/api/v1/schedule/%d", userID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) GetCourses(ctx context.Context, userID int64) ([]domain.Course, error) {
	var result []domain.Course
	if err := b.get(ctx, fmt.Sprintf("/api/v1/courses/%d", userID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) GetExams(ctx context.Context, userID int64) ([]domain.ExamEntry, error) {
	var result []domain.ExamEntry
	if err := b.get(ctx, fmt.Sprintf("/api/v1/exams/%d", userID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) GetGrades(ctx context.Context, userID int64) ([]domain.GradeRecord, error) {
	var result []domain.GradeRecord
	if err := b.get(ctx, fmt.Sprintf("/api/v1/grades/%d", userID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) GetDeadlines(ctx context.Context, userID int64) ([]domain.Deadline, error) {
	var result []domain.Deadline
	if err := b.get(ctx, fmt.Sprintf("/api/v1/deadlines/%d", userID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// endregion

// region Events & Clubs

func (b *Backend) ListEvents(ctx context.Context) ([]domain.Event, error) {
	var result []domain.Event
	if err := b.get(ctx, "/api/v1/events", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) ListNews(ctx context.Context) ([]domain.NewsItem, error) {
	var result []domain.NewsItem
	if err := b.get(ctx, "/api/v1/news", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) ListClubs(ctx context.Context) ([]domain.Club, error) {
	var result []domain.Club
	if err := b.get(ctx, "/api/v1/clubs", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) RSVPEvent(ctx context.Context, eventID int64, userID int64, registrationType string, note string) (string, error) {
	payload := map[string]any{
		"user_id":           userID,
		"registration_type": registrationType,
		"note":              note,
	}
	var resp struct {
		RegistrationID int64  `json:"registration_id"`
		Status         string `json:"status"`
	}
	if err := b.post(ctx, fmt.Sprintf("/api/v1/events/%d/rsvp", eventID), payload, &resp); err != nil {
		return "", err
	}
	return resp.Status, nil
}

func (b *Backend) CancelRSVP(ctx context.Context, eventID int64, userID int64) error {
	payload := map[string]any{
		"user_id": userID,
	}
	return b.post(ctx, fmt.Sprintf("/api/v1/events/%d/cancel", eventID), payload, nil)
}

func (b *Backend) ListUserEvents(ctx context.Context, userID int64) ([]domain.Event, error) {
	var result []domain.Event
	if err := b.get(ctx, fmt.Sprintf("/api/v1/events/user/%d", userID), nil, &result); err != nil {
		// fallback to global list if endpoint not available
		b.log.Warn().Int64("user_id", userID).Err(err).Msg("ListUserEvents endpoint not available, falling back to global events feed")
		return b.ListEvents(ctx)
	}
	return result, nil
}

// endregion

// region Admissions

func (b *Backend) ListAdmissionsPrograms(ctx context.Context) ([]domain.AdmissionProgram, error) {
	var result []domain.AdmissionProgram
	if err := b.get(ctx, "/api/v1/admissions/programs", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) ListAdmissionEvents(ctx context.Context) ([]domain.AdmissionEvent, error) {
	var result []domain.AdmissionEvent
	if err := b.get(ctx, "/api/v1/admissions/events", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) BookAdmissionEvent(ctx context.Context, eventID int64, applicantName, email, phone, note string) (int64, error) {
	payload := map[string]any{
		"applicant_name": applicantName,
		"email":          email,
		"phone":          phone,
		"note":           note,
	}
	var resp struct {
		BookingID int64 `json:"booking_id"`
	}
	if err := b.post(ctx, fmt.Sprintf("/api/v1/admissions/events/%d/book", eventID), payload, &resp); err != nil {
		return 0, err
	}
	return resp.BookingID, nil
}

func (b *Backend) ListAdmissionEventBookings(ctx context.Context, eventID int64) ([]domain.AdmissionEventBooking, error) {
	var result []domain.AdmissionEventBooking
	if err := b.get(ctx, fmt.Sprintf("/api/v1/admissions/events/%d/bookings", eventID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) SubmitAdmissionApplication(ctx context.Context, applicantName, email string, programID *int64, details map[string]any) (int64, error) {
	payload := map[string]any{
		"applicant_name": applicantName,
		"email":          email,
		"program_id":     programID,
		"details":        details,
	}
	var resp struct {
		ApplicationID int64 `json:"application_id"`
	}
	if err := b.post(ctx, "/api/v1/admissions/applications", payload, &resp); err != nil {
		return 0, err
	}
	return resp.ApplicationID, nil
}

func (b *Backend) UploadAdmissionDocument(ctx context.Context, applicationID int64, fileName, fileType, storageURL string) (int64, error) {
	payload := map[string]any{
		"application_id": applicationID,
		"file_name":      fileName,
		"file_type":      fileType,
		"storage_url":    storageURL,
	}
	var resp struct {
		DocumentID int64 `json:"document_id"`
	}
	if err := b.post(ctx, "/api/v1/admissions/upload", payload, &resp); err != nil {
		return 0, err
	}
	return resp.DocumentID, nil
}

func (b *Backend) AskAdmissionQuestion(ctx context.Context, question string) (string, error) {
	payload := map[string]string{
		"question": question,
	}
	var resp struct {
		Answer string `json:"answer"`
	}
	if err := b.post(ctx, "/api/v1/admissions/faq/query", payload, &resp); err != nil {
		return "", err
	}
	return resp.Answer, nil
}

// endregion

// region Dean's Office & Tuition

func (b *Backend) CreateDeanRequest(ctx context.Context, userID int64, requestType string, payload map[string]any) (int64, error) {
	body := map[string]any{
		"user_id":      userID,
		"request_type": requestType,
		"payload":      payload,
	}
	var resp struct {
		RequestID int64 `json:"request_id"`
	}
	if err := b.post(ctx, "/api/v1/dean/requests", body, &resp); err != nil {
		return 0, err
	}
	return resp.RequestID, nil
}

// endregion

// region Dormitory

func (b *Backend) GetDormRoom(ctx context.Context, studentID int64) (*domain.DormRoom, error) {
	var result domain.DormRoom
	if err := b.get(ctx, fmt.Sprintf("/api/v1/dorms/rooms/%d", studentID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *Backend) CreateDormMaintenance(ctx context.Context, studentID int64, requestType, description string) (int64, error) {
	payload := map[string]any{
		"student_id":   studentID,
		"request_type": requestType,
		"description":  description,
	}
	var resp struct {
		RequestID int64 `json:"request_id"`
	}
	if err := b.post(ctx, "/api/v1/dorms/maintenance", payload, &resp); err != nil {
		return 0, err
	}
	return resp.RequestID, nil
}

func (b *Backend) SubmitDormPayment(ctx context.Context, studentID int64, amount float64, reference string) (int64, error) {
	payload := map[string]any{
		"student_id": studentID,
		"amount":     amount,
		"reference":  reference,
	}
	var resp struct {
		PaymentID int64 `json:"payment_id"`
	}
	if err := b.post(ctx, "/api/v1/dorms/payments", payload, &resp); err != nil {
		return 0, err
	}
	return resp.PaymentID, nil
}

// endregion

// region Library

func (b *Backend) SearchBooks(ctx context.Context, q string) ([]domain.LibraryBook, error) {
	params := url.Values{}
	params.Set("q", q)
	var result []domain.LibraryBook
	if err := b.get(ctx, "/api/v1/library/books/search", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) ReserveBook(ctx context.Context, bookID, studentID int64) (int64, error) {
	payload := map[string]any{
		"book_id":    bookID,
		"student_id": studentID,
	}
	var resp struct {
		ReservationID int64 `json:"reservation_id"`
	}
	if err := b.post(ctx, "/api/v1/library/books/reserve", payload, &resp); err != nil {
		return 0, err
	}
	return resp.ReservationID, nil
}

func (b *Backend) ListBorrowedBooks(ctx context.Context, studentID int64) ([]domain.LibraryLoan, error) {
	var result []domain.LibraryLoan
	if err := b.get(ctx, fmt.Sprintf("/api/v1/library/borrowed/%d", studentID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// endregion

// region Support & AI

func (b *Backend) SubmitSupportTicket(ctx context.Context, category, subject, description string, userID *int64) (int64, error) {
	payload := map[string]any{
		"user_id":     userID,
		"category":    category,
		"subject":     subject,
		"description": description,
	}
	var resp struct {
		TicketID int64 `json:"ticket_id"`
	}
	if err := b.post(ctx, "/api/v1/support/tickets", payload, &resp); err != nil {
		return 0, err
	}
	return resp.TicketID, nil
}

func (b *Backend) SubmitSupportQuery(ctx context.Context, userID *int64, question string) (string, error) {
	payload := map[string]any{
		"user_id":  userID,
		"question": question,
	}
	var resp struct {
		Answer string `json:"answer"`
	}
	if err := b.post(ctx, "/api/v1/support/query", payload, &resp); err != nil {
		return "", err
	}
	return resp.Answer, nil
}

func (b *Backend) AdvisorChat(ctx context.Context, userID *int64, topic, prompt string) (string, error) {
	payload := map[string]any{
		"user_id": userID,
		"topic":   topic,
		"prompt":  prompt,
	}
	var resp struct {
		Response string `json:"response"`
	}
	if err := b.post(ctx, "/api/v1/ai/chat/advisor", payload, &resp); err != nil {
		return "", err
	}
	return resp.Response, nil
}

func (b *Backend) RunAIQuery(ctx context.Context, question string, filters map[string]any) (string, error) {
	payload := map[string]any{
		"question": question,
		"filters":  filters,
	}
	var resp struct {
		Answer string `json:"answer"`
	}
	if err := b.post(ctx, "/api/v1/ai/rag/query", payload, &resp); err != nil {
		return "", err
	}
	return resp.Answer, nil
}

func (b *Backend) CreateAISummary(ctx context.Context, text string) (string, error) {
	payload := map[string]string{
		"source_text": text,
	}
	var resp struct {
		Summary string `json:"summary"`
	}
	if err := b.post(ctx, "/api/v1/ai/summary/create", payload, &resp); err != nil {
		return "", err
	}
	return resp.Summary, nil
}

func (b *Backend) GenerateAIQuiz(ctx context.Context, prompt string, courseID *int64) ([]domain.QuizQuestion, error) {
	payload := map[string]any{
		"prompt":    prompt,
		"course_id": courseID,
	}
	var resp struct {
		Questions []domain.QuizQuestion `json:"questions"`
	}
	if err := b.post(ctx, "/api/v1/ai/quiz/generate", payload, &resp); err != nil {
		return nil, err
	}
	return resp.Questions, nil
}

func (b *Backend) TranscribeAudio(ctx context.Context, audioRef string) (string, error) {
	payload := map[string]string{
		"audio_ref": audioRef,
	}
	var resp struct {
		Transcript string `json:"transcript"`
	}
	if err := b.post(ctx, "/api/v1/ai/audio/transcribe", payload, &resp); err != nil {
		return "", err
	}
	return resp.Transcript, nil
}

// endregion

// region HR

func (b *Backend) GetVacations(ctx context.Context, employeeID int64) ([]domain.VacationRequest, error) {
	var result []domain.VacationRequest
	if err := b.get(ctx, fmt.Sprintf("/api/v1/hr/vacations/%d", employeeID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) RequestVacation(ctx context.Context, employeeID int64, startISO, endISO, vacationType string) (int64, error) {
	payload := map[string]any{
		"employee_id":   employeeID,
		"vacation_type": vacationType,
		"start_date":    startISO,
		"end_date":      endISO,
	}
	var resp struct {
		VacationRequestID int64 `json:"vacation_request_id"`
	}
	if err := b.post(ctx, "/api/v1/hr/vacations/request", payload, &resp); err != nil {
		return 0, err
	}
	return resp.VacationRequestID, nil
}

func (b *Backend) GetBusinessTrips(ctx context.Context, employeeID int64) ([]domain.BusinessTrip, error) {
	var result []domain.BusinessTrip
	if err := b.get(ctx, fmt.Sprintf("/api/v1/hr/business_trips/%d", employeeID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) RequestBusinessTrip(ctx context.Context, employeeID int64, destination, startISO, endISO, purpose string) (int64, error) {
	payload := map[string]any{
		"employee_id": employeeID,
		"destination": destination,
		"start_date":  startISO,
		"end_date":    endISO,
		"purpose":     purpose,
	}
	var resp struct {
		BusinessTripID int64 `json:"business_trip_id"`
	}
	if err := b.post(ctx, "/api/v1/hr/business_trips/request", payload, &resp); err != nil {
		return 0, err
	}
	return resp.BusinessTripID, nil
}

func (b *Backend) GetCertificates(ctx context.Context, employeeID int64) ([]domain.HRLetter, error) {
	var result []domain.HRLetter
	if err := b.get(ctx, fmt.Sprintf("/api/v1/hr/certificates/%d", employeeID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) RequestCertificate(ctx context.Context, employeeID int64, certificateType string) (int64, error) {
	payload := map[string]any{
		"employee_id":      employeeID,
		"certificate_type": certificateType,
	}
	var resp struct {
		CertificateRequestID int64 `json:"certificate_request_id"`
	}
	if err := b.post(ctx, "/api/v1/hr/certificates/request", payload, &resp); err != nil {
		return 0, err
	}
	return resp.CertificateRequestID, nil
}

// endregion

// region Visa

func (b *Backend) GetVisaApplications(ctx context.Context, userID int64) ([]map[string]any, error) {
	var result []map[string]any
	if err := b.get(ctx, fmt.Sprintf("/api/v1/visa/applications/%d", userID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) CreateVisaApplication(ctx context.Context, userID int64, applicationType string) (int64, error) {
	payload := map[string]any{
		"user_id":          userID,
		"application_type": applicationType,
	}
	var resp struct {
		ApplicationID int64 `json:"application_id"`
	}
	if err := b.post(ctx, "/api/v1/visa/applications", payload, &resp); err != nil {
		return 0, err
	}
	return resp.ApplicationID, nil
}

func (b *Backend) WithdrawVisaApplication(ctx context.Context, applicationID int64) error {
	return b.post(ctx, fmt.Sprintf("/api/v1/visa/applications/%d/withdraw", applicationID), nil, nil)
}

func (b *Backend) GetVisaDocuments(ctx context.Context, applicationID int64) ([]map[string]any, error) {
	var result []map[string]any
	if err := b.get(ctx, fmt.Sprintf("/api/v1/visa/applications/%d/documents", applicationID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Backend) UploadVisaDocument(ctx context.Context, applicationID int64, fileName, fileURL string) (int64, error) {
	payload := map[string]any{
		"file_name": fileName,
		"file_url":  fileURL,
	}
	var resp struct {
		DocumentID int64 `json:"document_id"`
	}
	if err := b.post(ctx, fmt.Sprintf("/api/v1/visa/applications/%d/documents", applicationID), payload, &resp); err != nil {
		return 0, err
	}
	return resp.DocumentID, nil
}

// endregion

// region Notifications

func (b *Backend) SendNotification(ctx context.Context, subject, body string, recipientID *int64) (int64, error) {
	payload := map[string]any{
		"recipient_id": recipientID,
		"subject":      subject,
		"body":         body,
	}
	var resp struct {
		NotificationID int64 `json:"notification_id"`
	}
	if err := b.post(ctx, "/api/v1/notifications/send", payload, &resp); err != nil {
		return 0, err
	}
	return resp.NotificationID, nil
}

// endregion
