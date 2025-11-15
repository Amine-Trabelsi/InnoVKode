from sqlalchemy import (
    Boolean,
    Column,
    DateTime,
    Float,
    ForeignKey,
    Integer,
    JSON,
    Numeric,
    String,
    Table,
    Text,
)
from sqlalchemy.sql import func

from .db import metadata

users_table = Table(
    "users",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("email", String(255), nullable=False, unique=True),
    Column("full_name_ru", String(255)),
    Column("full_name_en", String(255)),
    Column("role", String(50), nullable=False),  # student, employee, leadership, applicant
    Column("language", String(5), default="ru"),
    Column("is_foreign", Boolean, default=False),
    Column("dorm_room", String(50)),
    Column("faculty", String(120)),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

courses_table = Table(
    "courses",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("code", String(50), nullable=False, unique=True),
    Column("title", String(255), nullable=False),
    Column("description", Text),
    Column("faculty", String(120)),
    Column("ects", Float),
    Column("teacher_id", ForeignKey("users.id", ondelete="SET NULL")),
)

course_enrollments = Table(
    "course_enrollments",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("enrolled_at", DateTime(timezone=True), server_default=func.now()),
)

course_sessions = Table(
    "course_sessions",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("session_type", String(50), nullable=False),
    Column("start_time", DateTime(timezone=True), nullable=False),
    Column("end_time", DateTime(timezone=True), nullable=False),
    Column("location", String(120)),
    Column("week_label", String(20)),
)

exam_schedules = Table(
    "exam_schedules",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("exam_date", DateTime(timezone=True), nullable=False),
    Column("room", String(50)),
    Column("exam_format", String(50)),
)

grade_records = Table(
    "grade_records",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("grade", String(10)),
    Column("gpa_points", Numeric(3, 2)),
    Column("graded_on", DateTime(timezone=True), server_default=func.now()),
)

deadlines_table = Table(
    "deadlines",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("title", String(255), nullable=False),
    Column("due_date", DateTime(timezone=True), nullable=False),
    Column("category", String(80)),
    Column("status", String(40), default="open"),
    Column("details", Text),
)

notifications_table = Table(
    "notifications",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("recipient_id", ForeignKey("users.id")),
    Column("channel", String(40), default="in_app"),
    Column("subject", String(255)),
    Column("body", Text),
    Column("status", String(40), default="pending"),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

rooms_table = Table(
    "rooms",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("name", String(120), nullable=False),
    Column("location", String(255)),
    Column("capacity", Integer, nullable=False),
    Column("equipment", JSON, default=list),
)

room_bookings = Table(
    "room_bookings",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("room_id", ForeignKey("rooms.id"), nullable=False),
    Column("user_id", ForeignKey("users.id"), nullable=False),
    Column("start_time", DateTime(timezone=True), nullable=False),
    Column("end_time", DateTime(timezone=True), nullable=False),
    Column("purpose", String(255)),
    Column("status", String(30), default="confirmed"),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

events_table = Table(
    "events",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("title", String(255), nullable=False),
    Column("description", Text),
    Column("category", String(80)),
    Column("date_time", DateTime(timezone=True), nullable=False),
    Column("location", String(255)),
    Column("max_attendees", Integer),
    Column("current_attendees", Integer, default=0),
    Column("registration_type", String(30), default="attendee"),
)

event_registrations = Table(
    "event_registrations",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("event_id", ForeignKey("events.id"), nullable=False),
    Column("user_id", ForeignKey("users.id"), nullable=False),
    Column("registration_type", String(30), default="attendee"),
    Column("status", String(30), default="registered"),
    Column("note", Text),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

news_table = Table(
    "news_items",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("title", String(255), nullable=False),
    Column("category", String(80)),
    Column("body", Text),
    Column("published_at", DateTime(timezone=True), server_default=func.now()),
)

clubs_table = Table(
    "clubs",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("name", String(255), nullable=False),
    Column("description", Text),
    Column("meeting_schedule", String(255)),
    Column("contact", String(255)),
)

ai_sources = Table(
    "ai_sources",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("source_type", String(50), nullable=False),
    Column("reference", String(255), nullable=False),
    Column("title", String(255)),
    Column("metadata", JSON),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

ai_queries = Table(
    "ai_queries",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("query_text", Text, nullable=False),
    Column("response_text", Text),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

ai_quizzes = Table(
    "ai_quizzes",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id")),
    Column("prompt", Text, nullable=False),
    Column("questions", JSON, nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

ai_summaries = Table(
    "ai_summaries",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("source", Text, nullable=False),
    Column("summary", Text, nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

ai_transcriptions = Table(
    "ai_transcriptions",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("audio_ref", String(255), nullable=False),
    Column("transcript", Text),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

dean_requests = Table(
    "dean_requests",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("user_id", ForeignKey("users.id"), nullable=False),
    Column("request_type", String(80), nullable=False),
    Column("payload", JSON, nullable=False),
    Column("status", String(40), default="submitted"),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
    Column("updated_at", DateTime(timezone=True), server_default=func.now(), onupdate=func.now()),
)

admission_programs = Table(
    "admission_programs",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("title", String(255), nullable=False),
    Column("description", Text),
    Column("duration_years", Integer),
    Column("tuition", Numeric(10, 2)),
    Column("faculty", String(120)),
    Column("requirements", Text),
)

admission_events = Table(
    "admission_events",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("title", String(255), nullable=False),
    Column("event_type", String(80)),
    Column("description", Text),
    Column("date_time", DateTime(timezone=True)),
    Column("location", String(255)),
    Column("max_attendees", Integer),
    Column("current_attendees", Integer, default=0),
)

admission_applications = Table(
    "admission_applications",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("applicant_name", String(255), nullable=False),
    Column("email", String(255), nullable=False),
    Column("program_id", ForeignKey("admission_programs.id")),
    Column("status", String(40), default="received"),
    Column("submitted_at", DateTime(timezone=True), server_default=func.now()),
    Column("details", JSON),
)

admission_documents = Table(
    "admission_documents",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("application_id", ForeignKey("admission_applications.id"), nullable=False),
    Column("file_name", String(255)),
    Column("file_type", String(40)),
    Column("storage_url", String(255)),
    Column("uploaded_at", DateTime(timezone=True), server_default=func.now()),
)

admission_faq_queries = Table(
    "admission_faq_queries",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("question", Text, nullable=False),
    Column("response", Text),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

admission_event_bookings = Table(
    "admission_event_bookings",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("event_id", ForeignKey("admission_events.id"), nullable=False),
    Column("applicant_name", String(255), nullable=False),
    Column("email", String(255), nullable=False),
    Column("phone", String(50)),
    Column("status", String(30), default="confirmed"),
    Column("note", Text),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

teaching_attendance = Table(
    "teaching_attendance",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("professor_id", ForeignKey("users.id"), nullable=False),
    Column("session_date", DateTime(timezone=True), nullable=False),
    Column("attendance", JSON, nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

teaching_grade_uploads = Table(
    "teaching_grade_uploads",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("professor_id", ForeignKey("users.id"), nullable=False),
    Column("payload", JSON, nullable=False),
    Column("uploaded_at", DateTime(timezone=True), server_default=func.now()),
)

teaching_submissions = Table(
    "teaching_submissions",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("title", String(255), nullable=False),
    Column("submitted_at", DateTime(timezone=True), server_default=func.now()),
    Column("status", String(40), default="pending"),
    Column("grade", String(10)),
)

teaching_announcements = Table(
    "teaching_announcements",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("professor_id", ForeignKey("users.id"), nullable=False),
    Column("message", Text, nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

teaching_feedback = Table(
    "teaching_feedback",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("course_id", ForeignKey("courses.id"), nullable=False),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("rating", Integer, nullable=False),
    Column("comment", Text),
    Column("submitted_at", DateTime(timezone=True), server_default=func.now()),
)

vacation_requests = Table(
    "vacation_requests",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("employee_id", ForeignKey("users.id"), nullable=False),
    Column("start_date", DateTime(timezone=True), nullable=False),
    Column("end_date", DateTime(timezone=True), nullable=False),
    Column("vacation_type", String(40), default="paid"),
    Column("status", String(40), default="pending"),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

business_trip_requests = Table(
    "business_trip_requests",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("employee_id", ForeignKey("users.id"), nullable=False),
    Column("destination", String(255), nullable=False),
    Column("start_date", DateTime(timezone=True), nullable=False),
    Column("end_date", DateTime(timezone=True), nullable=False),
    Column("purpose", Text),
    Column("status", String(40), default="pending"),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

hr_certificates = Table(
    "hr_certificates",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("employee_id", ForeignKey("users.id"), nullable=False),
    Column("certificate_type", String(80), nullable=False),
    Column("status", String(40), default="processing"),
    Column("download_url", String(255)),
    Column("requested_at", DateTime(timezone=True), server_default=func.now()),
)

hr_notifications = Table(
    "hr_notifications",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("title", String(255), nullable=False),
    Column("body", Text),
    Column("published_at", DateTime(timezone=True), server_default=func.now()),
)

dorm_rooms = Table(
    "dorm_rooms",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("room_number", String(50), nullable=False),
    Column("building", String(80)),
    Column("balance", Numeric(10, 2), default=0),
)

dorm_requests = Table(
    "dorm_requests",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("request_type", String(80), nullable=False),
    Column("description", Text),
    Column("status", String(40), default="open"),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

dorm_payments = Table(
    "dorm_payments",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("amount", Numeric(10, 2), nullable=False),
    Column("paid_at", DateTime(timezone=True), server_default=func.now()),
    Column("reference", String(120)),
)

library_books = Table(
    "library_books",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("title", String(255), nullable=False),
    Column("author", String(255)),
    Column("keywords", JSON, default=list),
    Column("available_copies", Integer, default=1),
)

library_reservations = Table(
    "library_reservations",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("book_id", ForeignKey("library_books.id"), nullable=False),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("reserved_at", DateTime(timezone=True), server_default=func.now()),
    Column("status", String(40), default="pending"),
)

library_loans = Table(
    "library_loans",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("book_id", ForeignKey("library_books.id"), nullable=False),
    Column("student_id", ForeignKey("users.id"), nullable=False),
    Column("borrowed_at", DateTime(timezone=True), server_default=func.now()),
    Column("due_at", DateTime(timezone=True)),
    Column("status", String(40), default="borrowed"),
)

library_digital_assets = Table(
    "library_digital_assets",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("book_id", ForeignKey("library_books.id")),
    Column("format", String(40)),
    Column("access_url", String(255)),
    Column("metadata", JSON),
)

support_tickets = Table(
    "support_tickets",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("user_id", ForeignKey("users.id")),
    Column("category", String(80), nullable=False),
    Column("subject", String(255), nullable=False),
    Column("description", Text),
    Column("status", String(40), default="open"),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

support_queries = Table(
    "support_queries",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("user_id", ForeignKey("users.id")),
    Column("question", Text, nullable=False),
    Column("answer", Text),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

ai_advisor_sessions = Table(
    "ai_advisor_sessions",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("user_id", ForeignKey("users.id")),
    Column("topic", String(255)),
    Column("prompt", Text, nullable=False),
    Column("response", Text),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
)

visa_applications = Table(
    "visa_applications",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("user_id", ForeignKey("users.id"), nullable=False),
    Column("application_type", String(80), nullable=False),  # visa_renewal, registration_renewal
    Column("status", String(40), default="pending"),  # pending, withdrawn, approved, rejected
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
    Column("updated_at", DateTime(timezone=True), server_default=func.now(), onupdate=func.now()),
)

visa_documents = Table(
    "visa_documents",
    metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("application_id", ForeignKey("visa_applications.id"), nullable=False),
    Column("file_name", String(255), nullable=False),
    Column("file_url", String(255), nullable=False),
    Column("uploaded_at", DateTime(timezone=True), server_default=func.now()),
)
