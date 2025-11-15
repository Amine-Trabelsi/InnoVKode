package domain

import "time"

type UserProfile struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	NameRU     string    `json:"full_name_ru"`
	NameEN     string    `json:"full_name_en"`
	Role       Role      `json:"role"`
	Language   Language  `json:"language"`
	IsForeign  bool      `json:"is_foreign"`
	DormRoom   string    `json:"dorm_room"`
	Faculty    string    `json:"faculty"`
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
}

type ScheduleEntry struct {
	SessionID   int64     `json:"session_id"`
	SessionType string    `json:"session_type"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	WeekLabel   string    `json:"week_label"`
	CourseID    int64     `json:"course_id"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
}

type Course struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Faculty     string `json:"faculty"`
}

type ExamEntry struct {
	ExamID   int64     `json:"exam_id"`
	Date     time.Time `json:"exam_date"`
	Room     string    `json:"room"`
	Format   string    `json:"exam_format"`
	CourseID int64     `json:"course_id"`
	Code     string    `json:"code"`
	Title    string    `json:"title"`
}

type GradeRecord struct {
	GradeID   int64     `json:"grade_id"`
	Grade     string    `json:"grade"`
	GPAPoints float64   `json:"gpa_points"`
	GradedOn  time.Time `json:"graded_on"`
	CourseID  int64     `json:"course_id"`
	Code      string    `json:"code"`
	Title     string    `json:"title"`
}

type Deadline struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	DueDate  time.Time `json:"due_date"`
	Category string    `json:"category"`
	Status   string    `json:"status"`
	Details  string    `json:"details"`
}

type Event struct {
	ID                   int64     `json:"id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Category             string    `json:"category"`
	DateTime             time.Time `json:"date_time"`
	Location             string    `json:"location"`
	MaxAttendees         int64     `json:"max_attendees"`
	CurrentAttendees     int64     `json:"current_attendees"`
	RegistrationType     string    `json:"registration_type"`
	UserRegistrationType string    `json:"user_registration_type"`
}

type Club struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	MeetingSchedule string `json:"meeting_schedule"`
	Contact         string `json:"contact"`
}

type EventRegistration struct {
	ID      int64  `json:"registration_id"`
	EventID int64  `json:"event_id"`
	UserID  int64  `json:"user_id"`
	Status  string `json:"status"`
	Note    string `json:"note"`
}

type NewsItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
}

type AdmissionProgram struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	DurationYears int    `json:"duration_years"`
	Tuition       string `json:"tuition"`
	Faculty       string `json:"faculty"`
	Requirements  string `json:"requirements"`
}

type AdmissionEvent struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	EventType        string    `json:"event_type"`
	Description      string    `json:"description"`
	DateTime         time.Time `json:"date_time"`
	Location         string    `json:"location"`
	MaxAttendees     int64     `json:"max_attendees"`
	CurrentAttendees int64     `json:"current_attendees"`
}

type AdmissionEventBooking struct {
	ID            int64     `json:"id"`
	EventID       int64     `json:"event_id"`
	ApplicantName string    `json:"applicant_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

type DormRoom struct {
	ID        int64   `json:"id"`
	StudentID int64   `json:"student_id"`
	Room      string  `json:"room_number"`
	Building  string  `json:"building"`
	Balance   float64 `json:"balance"`
}

type LibraryBook struct {
	ID              int64    `json:"id"`
	Title           string   `json:"title"`
	Author          string   `json:"author"`
	Keywords        []string `json:"keywords"`
	AvailableCopies int      `json:"available_copies"`
}

type LibraryReservation struct {
	ID        int64  `json:"reservation_id"`
	BookID    int64  `json:"book_id"`
	StudentID int64  `json:"student_id"`
	Status    string `json:"status"`
}

type LibraryLoan struct {
	LoanID     int64     `json:"loan_id"`
	BorrowedAt time.Time `json:"borrowed_at"`
	DueAt      time.Time `json:"due_at"`
	Status     string    `json:"status"`
	BookID     int64     `json:"book_id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
}

type VacationRequest struct {
	ID           int64     `json:"id"`
	EmployeeID   int64     `json:"employee_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	VacationType string    `json:"vacation_type"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type BusinessTrip struct {
	ID          int64     `json:"id"`
	EmployeeID  int64     `json:"employee_id"`
	Destination string    `json:"destination"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Purpose     string    `json:"purpose"`
	Status      string    `json:"status"`
}

type HRLetter struct {
	ID              int64     `json:"id"`
	EmployeeID      int64     `json:"employee_id"`
	CertificateType string    `json:"certificate_type"`
	Status          string    `json:"status"`
	DownloadURL     string    `json:"download_url"`
	RequestedAt     time.Time `json:"requested_at"`
}

type QuizQuestion struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}
