from __future__ import annotations

from datetime import datetime, timedelta, timezone
from decimal import Decimal

from sqlalchemy import insert, select

from .db import AsyncSessionLocal
from .tables import (
    admission_applications,
    admission_documents,
    admission_event_bookings,
    admission_events,
    admission_faq_queries,
    admission_programs,
    ai_advisor_sessions,
    ai_queries,
    ai_quizzes,
    ai_sources,
    ai_summaries,
    ai_transcriptions,
    business_trip_requests,
    clubs_table,
    course_enrollments,
    course_sessions,
    courses_table,
    deadlines_table,
    dean_requests,
    dorm_payments,
    dorm_requests,
    dorm_rooms,
    event_registrations,
    events_table,
    exam_schedules,
    grade_records,
    hr_certificates,
    hr_notifications,
    library_books,
    library_digital_assets,
    library_loans,
    library_reservations,
    news_table,
    notifications_table,
    room_bookings,
    rooms_table,
    support_queries,
    support_tickets,
    teaching_announcements,
    teaching_attendance,
    teaching_feedback,
    teaching_grade_uploads,
    teaching_submissions,
    users_table,
    vacation_requests,
)

UTC = timezone.utc
BASE_DATETIME = datetime(2025, 1, 13, 8, 0, tzinfo=UTC)


def dt(days: int = 0, hours: int = 0) -> datetime:
    return BASE_DATETIME + timedelta(days=days, hours=hours)
USERS = [
    (
        "anna",
        {
            "email": "anna.petrov@univ.ru",
            "full_name_ru": "Анна Петрова",
            "full_name_en": "Anna Petrova",
            "role": "student",
            "language": "ru",
            "is_foreign": False,
            "dorm_room": "A-201",
            "faculty": "Computer Science",
        },
    ),
    (
        "boris",
        {
            "email": "boris.ivanov@univ.ru",
            "full_name_ru": "Борис Иванов",
            "full_name_en": "Boris Ivanov",
            "role": "student",
            "language": "ru",
            "is_foreign": False,
            "dorm_room": "B-120",
            "faculty": "Business",
        },
    ),
    (
        "chen",
        {
            "email": "chen.li@univ.ru",
            "full_name_ru": "Чэнь Ли",
            "full_name_en": "Chen Li",
            "role": "student",
            "language": "en",
            "is_foreign": True,
            "dorm_room": "C-310",
            "faculty": "Computer Science",
        },
    ),
    (
        "ilya",
        {
            "email": "ilya.smirnov@univ.ru",
            "full_name_ru": "Илья Смирнов",
            "full_name_en": "Ilya Smirnov",
            "role": "employee",
            "language": "ru",
            "faculty": "Computer Science",
            "is_foreign": False,
            "dorm_room": None,
        },
    ),
    (
        "sofia",
        {
            "email": "sofia.morozova@univ.ru",
            "full_name_ru": "София Морозова",
            "full_name_en": "Sofia Morozova",
            "role": "employee",
            "language": "ru",
            "faculty": "HR",
            "is_foreign": False,
            "dorm_room": None,
        },
    ),
    (
        "rector",
        {
            "email": "rector.office@univ.ru",
            "full_name_ru": "Ректорский офис",
            "full_name_en": "Rector Office",
            "role": "leadership",
            "language": "ru",
            "faculty": "Administration",
            "is_foreign": False,
            "dorm_room": None,
        },
    ),
    (
        "daria",
        {
            "email": "daria.kozlova@univ.ru",
            "full_name_ru": "Дарья Козлова",
            "full_name_en": "Daria Kozlova",
            "role": "student",
            "language": "ru",
            "dorm_room": "A-305",
            "faculty": "Design",
            "is_foreign": False,
        },
    ),
    (
        "olga",
        {
            "email": "librarian@univ.ru",
            "full_name_ru": "Ольга Лебедева",
            "full_name_en": "Olga Lebedeva",
            "role": "employee",
            "language": "ru",
            "faculty": "Library",
            "is_foreign": False,
            "dorm_room": None,
        },
    ),
]


COURSES = [
    (
        "cs101",
        {
            "code": "CS101",
            "title": "Intro to Programming",
            "description": "Python foundations for engineers.",
            "faculty": "Computer Science",
            "ects": 4.0,
            "teacher": "ilya",
        },
    ),
    (
        "cs240",
        {
            "code": "CS240",
            "title": "Algorithms & Data Structures",
            "description": "Core algorithms with practical labs.",
            "faculty": "Computer Science",
            "ects": 5.0,
            "teacher": "ilya",
        },
    ),
    (
        "bus310",
        {
            "code": "BUS310",
            "title": "Project Management",
            "description": "Agile delivery for cross-functional teams.",
            "faculty": "Business",
            "ects": 3.0,
            "teacher": "sofia",
        },
    ),
]

COURSE_ENROLLMENTS = [
    {"student": "anna", "course": "cs101"},
    {"student": "anna", "course": "cs240"},
    {"student": "boris", "course": "cs101"},
    {"student": "boris", "course": "bus310"},
    {"student": "chen", "course": "cs240"},
    {"student": "daria", "course": "bus310"},
]

COURSE_SESSIONS = [
    {"course": "cs101", "session_type": "lecture", "start": dt(0, 2), "end": dt(0, 4), "location": "A-101", "week": "week1"},
    {"course": "cs101", "session_type": "lab", "start": dt(2, 3), "end": dt(2, 5), "location": "Lab 3", "week": "week1"},
    {"course": "cs240", "session_type": "lecture", "start": dt(1, 1), "end": dt(1, 3), "location": "A-205", "week": "week1"},
    {"course": "cs240", "session_type": "seminar", "start": dt(3, 2), "end": dt(3, 4), "location": "A-207", "week": "week1"},
    {"course": "bus310", "session_type": "workshop", "start": dt(4, 4), "end": dt(4, 6), "location": "Innovation Hub", "week": "week1"},
]

EXAM_SCHEDULES = [
    {"course": "cs101", "date": dt(21, 3), "room": "A-201", "format": "written"},
    {"course": "cs240", "date": dt(23, 2), "room": "A-205", "format": "oral"},
    {"course": "bus310", "date": dt(25, 4), "room": "B-101", "format": "project"},
]

GRADE_RECORDS = [
    {"student": "anna", "course": "cs101", "grade": "A", "gpa": Decimal("4.00")},
    {"student": "anna", "course": "cs240", "grade": "B", "gpa": Decimal("3.00")},
    {"student": "boris", "course": "cs101", "grade": "B+", "gpa": Decimal("3.30")},
    {"student": "boris", "course": "bus310", "grade": "A-", "gpa": Decimal("3.70")},
    {"student": "chen", "course": "cs240", "grade": "A", "gpa": Decimal("4.00")},
]

DEADLINES = [
    {"student": "anna", "title": "Scholarship essay", "due": dt(5), "category": "admin"},
    {"student": "boris", "title": "Lab report", "due": dt(3), "category": "academic"},
    {"student": "chen", "title": "Visa check-in", "due": dt(7), "category": "immigration"},
]

NOTIFICATIONS = [
    {"recipient": "anna", "channel": "email", "subject": "Workshop reminder", "body": "Join the AI workshop on Friday.", "status": "sent"},
    {"recipient": None, "channel": "in_app", "subject": "Campus Wi-Fi", "body": "Maintenance window scheduled this weekend.", "status": "pending"},
]

ROOMS = [
    ("hub", {"name": "Study Hub 1", "location": "Library 2F", "capacity": 12, "equipment": ["monitor", "whiteboard"]}),
    ("lab", {"name": "Innovation Lab", "location": "Tech Park", "capacity": 18, "equipment": ["3D printer", "VR set"]}),
    ("conf", {"name": "Conference Room B", "location": "Admin Building", "capacity": 14, "equipment": ["projector", "speakerphone"]}),
]

ROOM_BOOKINGS = [
    {"room": "hub", "user": "anna", "start": dt(1, 1), "end": dt(1, 3), "purpose": "Study group"},
    {"room": "lab", "user": "boris", "start": dt(2, 2), "end": dt(2, 5), "purpose": "Project sprint"},
]

EVENTS = [
    (
        "ai_day",
        {
            "title": "AI Demo Day",
            "description": "Showcase of student AI projects.",
            "category": "tech",
            "date_time": dt(3, 4),
            "location": "Innovation Hall",
            "max_attendees": 60,
            "current_attendees": 2,
            "registration_type": "attendee",
        },
    ),
    (
        "tour",
        {
            "title": "Campus Tour",
            "description": "Guided walk for applicants.",
            "category": "admissions",
            "date_time": dt(4, 1),
            "location": "Main Gate",
            "max_attendees": 30,
            "current_attendees": 1,
            "registration_type": "attendee",
        },
    ),
]

EVENT_REGISTRATIONS = [
    {"event": "ai_day", "user": "anna", "status": "registered"},
    {"event": "ai_day", "user": "boris", "status": "registered"},
    {"event": "tour", "user": "daria", "status": "registered"},
]

NEWS = [
    {"title": "New AI Lab Opens", "category": "research", "body": "State-of-the-art lab launched.", "published_at": dt(-1)},
    {"title": "Dorm Renovation", "category": "campus", "body": "North dormitory gets upgrades.", "published_at": dt(-2)},
]

CLUBS = [
    {"name": "Robotics Club", "description": "Build autonomous robots.", "meeting_schedule": "Wed 18:00", "contact": "roboclub@univ.ru"},
    {"name": "Debate Society", "description": "Weekly debates and competitions.", "meeting_schedule": "Fri 17:00", "contact": "debate@univ.ru"},
]

AI_SOURCES = [
    {"source_type": "youtube", "reference": "https://youtu.be/demo1", "title": "Linear Algebra Lecture", "metadata": {"duration": 3600}},
    {"source_type": "pdf", "reference": "s3://bucket/notes.pdf", "title": "Distributed Systems Notes", "metadata": {"pages": 48}},
]

AI_QUERIES = [
    {"query_text": "Explain AVL trees", "response_text": "AVL trees maintain balance using rotations."},
]

AI_QUIZZES = [
    {"course": "cs240", "prompt": "Generate quiz for recursion"},
]

AI_SUMMARIES = [
    {"source": "Lecture on graphs", "summary": "Graphs connect vertices via edges; DFS and BFS explore structures."},
]

AI_TRANSCRIPTIONS = [
    {"audio_ref": "gs://bucket/q&a.mp3", "transcript": "Professor answers on grading policy."},
]

DEAN_REQUESTS = [
    {"user": "anna", "request_type": "certificate", "payload": {"certificate": "Enrollment letter"}, "status": "submitted"},
]

ADMISSION_PROGRAMS = [
    ("cs_bsc", {"title": "BSc Computer Science", "description": "Four-year CS track", "duration_years": 4, "tuition": Decimal("420000.00"), "faculty": "Computer Science", "requirements": "High school diploma, math exam"}),
    ("design_ba", {"title": "BA Design", "description": "Studio-focused track", "duration_years": 4, "tuition": Decimal("380000.00"), "faculty": "Design", "requirements": "Portfolio, interview"}),
]

ADMISSION_EVENTS = [
    (
        "open_day_jan",
        {
            "title": "Open Day January",
            "event_type": "open_day",
            "description": "Campus presentations and guided tours",
            "date_time": dt(6, 2),
            "location": "Main Auditorium",
            "max_attendees": 100,
            "current_attendees": 2,
        },
    ),
    (
        "design_tour",
        {
            "title": "Design Tour",
            "event_type": "tour",
            "description": "Studios visit",
            "date_time": dt(8, 3),
            "location": "Design Center",
            "max_attendees": 25,
            "current_attendees": 1,
        },
    ),
]

ADMISSION_EVENT_BOOKINGS = [
    {
        "event": "open_day_jan",
        "applicant_name": "Maria Volkova",
        "email": "maria.volkova@example.com",
        "phone": "+7 900 123 4567",
    },
    {
        "event": "open_day_jan",
        "applicant_name": "Alexey Petrov",
        "email": "alexey.petrov@example.com",
        "phone": "+7 900 765 4321",
    },
    {
        "event": "design_tour",
        "applicant_name": "Elena Sokolova",
        "email": "elena.sokolova@example.com",
        "note": "Interested in product design program",
    },
]

ADMISSION_APPLICATIONS = [
    ("app_ivan", {"applicant_name": "Ivan Applicant", "email": "ivan.applicant@mail.com", "program": "cs_bsc", "status": "review"}),
    ("app_sara", {"applicant_name": "Sara Global", "email": "sara.global@mail.com", "program": "design_ba", "status": "documents"}),
]

ADMISSION_DOCUMENTS = [
    {"application": "app_ivan", "file_name": "passport.pdf", "file_type": "pdf", "storage_url": "s3://docs/passport.pdf"},
    {"application": "app_sara", "file_name": "portfolio.zip", "file_type": "zip", "storage_url": "s3://docs/portfolio.zip"},
]

ADMISSION_FAQ = [
    {"question": "What are tuition deadlines?", "response": "Invoices are due August 10th."},
]

TEACHING_ATTENDANCE = [
    {
        "course": "cs101",
        "professor": "ilya",
        "session_date": dt(0),
        "attendance": [
            {"student": "anna", "present": True},
            {"student": "boris", "present": True},
        ],
    },
    {
        "course": "cs240",
        "professor": "ilya",
        "session_date": dt(1),
        "attendance": [
            {"student": "anna", "present": False},
            {"student": "chen", "present": True},
        ],
    },
]

TEACHING_GRADE_UPLOADS = [
    {
        "course": "cs101",
        "professor": "ilya",
        "payload": [
            {"student": "anna", "grade": "A"},
            {"student": "boris", "grade": "B+"},
        ],
    }
]

TEACHING_SUBMISSIONS = [
    {"course": "cs101", "student": "anna", "title": "HW1", "status": "graded", "grade": "A"},
    {"course": "cs240", "student": "chen", "title": "Lab2", "status": "submitted"},
]

TEACHING_ANNOUNCEMENTS = [
    {"course": "cs240", "professor": "ilya", "message": "Extra office hours on Thursday."},
]

TEACHING_FEEDBACK = [
    {"course": "cs101", "student": "anna", "rating": 5, "comment": "Great explanations."},
    {"course": "cs101", "student": "boris", "rating": 4, "comment": "Would like more examples."},
]

VACATION_REQUESTS = [
    {"employee": "sofia", "start": dt(10), "end": dt(15), "vacation_type": "paid", "status": "approved"},
    {"employee": "ilya", "start": dt(30), "end": dt(35), "vacation_type": "unpaid", "status": "pending"},
]

BUSINESS_TRIPS = [
    {"employee": "sofia", "destination": "Moscow", "start": dt(12), "end": dt(14), "purpose": "HR summit", "status": "approved"},
]

HR_CERTIFICATES = [
    {"employee": "ilya", "certificate_type": "employment", "status": "ready", "download_url": "https://files/employment.pdf"},
    {"employee": "sofia", "certificate_type": "income", "status": "processing"},
]

HR_NOTIFICATIONS = [
    {"title": "HR Webinar", "body": "Join Thursday for policy updates.", "published_at": dt(-3)},
]

DORM_ROOMS = [
    {"student": "anna", "room_number": "A-201", "building": "North", "balance": Decimal("1200.00")},
    {"student": "chen", "room_number": "C-310", "building": "International", "balance": Decimal("800.00")},
]

DORM_REQUESTS = [
    {"student": "anna", "request_type": "maintenance", "description": "Heating issue", "status": "open"},
]

DORM_PAYMENTS = [
    {"student": "anna", "amount": Decimal("20000.00"), "reference": "TXN123"},
    {"student": "chen", "amount": Decimal("22000.00"), "reference": "TXN124"},
]

LIBRARY_BOOKS = [
    ("deep_learning", {"title": "Deep Learning", "author": "Ian Goodfellow", "keywords": ["ai", "ml"], "available_copies": 2}),
    ("design_things", {"title": "Design of Everyday Things", "author": "Don Norman", "keywords": ["design"], "available_copies": 1}),
    ("pm_essentials", {"title": "Project Management Essentials", "author": "Rita Mulcahy", "keywords": ["pm"], "available_copies": 3}),
]

LIBRARY_DIGITAL_ASSETS = [
    {"book": "deep_learning", "format": "pdf", "access_url": "https://library/dl.pdf", "metadata": {"size": "5MB"}},
    {"book": "pm_essentials", "format": "epub", "access_url": "https://library/pm.epub", "metadata": {"size": "2MB"}},
]

LIBRARY_RESERVATIONS = [
    {"book": "design_things", "student": "daria", "status": "pending"},
]

LIBRARY_LOANS = [
    {"book": "deep_learning", "student": "anna", "borrowed_at": dt(-2), "due_at": dt(12), "status": "borrowed"},
    {"book": "pm_essentials", "student": "boris", "borrowed_at": dt(-1), "due_at": dt(10), "status": "borrowed"},
]

SUPPORT_QUERIES = [
    {"user": "anna", "question": "How to reset Wi-Fi password?", "answer": "Use the IT portal to reset."},
]

SUPPORT_TICKETS = [
    {"user": "boris", "category": "it", "subject": "Laptop issue", "description": "Screen flicker in lab", "status": "open"},
]

AI_ADVISOR = [
    {"user": "anna", "topic": "career", "prompt": "Which electives help AI research?", "response": "Choose CS240 and math electives."},
]


async def seed_initial_data() -> None:
    async with AsyncSessionLocal() as session:
        if await _table_has_rows(session, users_table):
            return

        user_map = await _insert_with_keys(session, users_table, USERS)
        course_map = await _insert_with_keys(session, courses_table, _prepare_courses(user_map))
        room_map = await _insert_with_keys(session, rooms_table, ROOMS)
        event_map = await _insert_with_keys(session, events_table, EVENTS)
        program_map = await _insert_with_keys(session, admission_programs, ADMISSION_PROGRAMS)
        book_map = await _insert_with_keys(session, library_books, LIBRARY_BOOKS)

        await _bulk_insert(session, course_enrollments, _prepare_course_enrollments(user_map, course_map))
        await _bulk_insert(session, course_sessions, _prepare_course_sessions(course_map))
        await _bulk_insert(session, exam_schedules, _prepare_exam_schedules(course_map))
        await _bulk_insert(session, grade_records, _prepare_grade_records(user_map, course_map))
        await _bulk_insert(session, deadlines_table, _prepare_deadlines(user_map))
        await _bulk_insert(session, notifications_table, _prepare_notifications(user_map))
        await _bulk_insert(session, room_bookings, _prepare_room_bookings(room_map, user_map))
        await _bulk_insert(session, event_registrations, _prepare_event_registrations(event_map, user_map))
        await _bulk_insert(session, news_table, NEWS)
        await _bulk_insert(session, clubs_table, CLUBS)
        await _bulk_insert(session, ai_sources, AI_SOURCES)
        await _bulk_insert(session, ai_queries, AI_QUERIES)
        await _bulk_insert(session, ai_quizzes, _prepare_ai_quizzes(course_map))
        await _bulk_insert(session, ai_summaries, AI_SUMMARIES)
        await _bulk_insert(session, ai_transcriptions, AI_TRANSCRIPTIONS)
        await _bulk_insert(session, dean_requests, _prepare_dean_requests(user_map))
        app_map = await _insert_with_keys(session, admission_applications, _prepare_admission_applications(program_map))
        await _bulk_insert(session, admission_documents, _prepare_admission_documents(app_map))
        await _bulk_insert(session, admission_faq_queries, ADMISSION_FAQ)
        await _bulk_insert(session, teaching_attendance, _prepare_teaching_attendance(user_map, course_map))
        await _bulk_insert(session, teaching_grade_uploads, _prepare_teaching_grade_uploads(user_map, course_map))
        await _bulk_insert(session, teaching_submissions, _prepare_teaching_submissions(user_map, course_map))
        await _bulk_insert(session, teaching_announcements, _prepare_teaching_announcements(user_map, course_map))
        await _bulk_insert(session, teaching_feedback, _prepare_teaching_feedback(user_map, course_map))
        await _bulk_insert(session, vacation_requests, _prepare_vacation_requests(user_map))
        await _bulk_insert(session, business_trip_requests, _prepare_business_trips(user_map))
        await _bulk_insert(session, hr_certificates, _prepare_hr_certificates(user_map))
        await _bulk_insert(session, hr_notifications, HR_NOTIFICATIONS)
        await _bulk_insert(session, dorm_rooms, _prepare_dorm_rooms(user_map))
        await _bulk_insert(session, dorm_requests, _prepare_dorm_requests(user_map))
        await _bulk_insert(session, dorm_payments, _prepare_dorm_payments(user_map))
        await _bulk_insert(session, library_digital_assets, _prepare_library_digital_assets(book_map))
        await _bulk_insert(session, library_reservations, _prepare_library_reservations(book_map, user_map))
        await _bulk_insert(session, library_loans, _prepare_library_loans(book_map, user_map))
        await _bulk_insert(session, support_queries, _prepare_support_queries(user_map))
        await _bulk_insert(session, support_tickets, _prepare_support_tickets(user_map))
        await _bulk_insert(session, ai_advisor_sessions, _prepare_ai_advisor(user_map))
        admission_event_map = await _insert_with_keys(session, admission_events, ADMISSION_EVENTS)
        await _bulk_insert(session, admission_event_bookings, _prepare_admission_event_bookings(admission_event_map))

        await session.commit()


async def _table_has_rows(session, table) -> bool:
    result = await session.execute(select(list(table.primary_key.columns)[0]).limit(1))
    return result.first() is not None


async def _insert_with_keys(session, table, rows_with_keys):
    if not rows_with_keys:
        return {}
    keys = []
    rows = []
    for key, record in rows_with_keys:
        keys.append(key)
        rows.append(record)
    result = await session.execute(insert(table).returning(table.c.id), rows)
    ids = result.scalars().all()
    return dict(zip(keys, ids))


async def _bulk_insert(session, table, rows):
    if rows:
        await session.execute(insert(table), rows)


def _prepare_courses(user_map):
    prepared = []
    for key, record in COURSES:
        data = record.copy()
        teacher_key = data.pop("teacher", None)
        data["teacher_id"] = user_map.get(teacher_key) if teacher_key else None
        prepared.append((key, data))
    return prepared


def _prepare_course_enrollments(user_map, course_map):
    return [
        {
            "student_id": user_map[item["student"]],
            "course_id": course_map[item["course"]],
        }
        for item in COURSE_ENROLLMENTS
    ]


def _prepare_course_sessions(course_map):
    return [
        {
            "course_id": course_map[item["course"]],
            "session_type": item["session_type"],
            "start_time": item["start"],
            "end_time": item["end"],
            "location": item["location"],
            "week_label": item["week"],
        }
        for item in COURSE_SESSIONS
    ]


def _prepare_exam_schedules(course_map):
    return [
        {
            "course_id": course_map[item["course"]],
            "exam_date": item["date"],
            "room": item["room"],
            "exam_format": item["format"],
        }
        for item in EXAM_SCHEDULES
    ]


def _prepare_grade_records(user_map, course_map):
    return [
        {
            "student_id": user_map[item["student"]],
            "course_id": course_map[item["course"]],
            "grade": item["grade"],
            "gpa_points": item["gpa"],
        }
        for item in GRADE_RECORDS
    ]


def _prepare_deadlines(user_map):
    return [
        {
            "student_id": user_map[item["student"]],
            "title": item["title"],
            "due_date": item["due"],
            "category": item["category"],
            "status": "open",
        }
        for item in DEADLINES
    ]


def _prepare_notifications(user_map):
    rows = []
    for item in NOTIFICATIONS:
        rows.append(
            {
                "recipient_id": user_map.get(item["recipient"]) if item["recipient"] else None,
                "channel": item["channel"],
                "subject": item["subject"],
                "body": item["body"],
                "status": item["status"],
            }
        )
    return rows


def _prepare_room_bookings(room_map, user_map):
    return [
        {
            "room_id": room_map[item["room"]],
            "user_id": user_map[item["user"]],
            "start_time": item["start"],
            "end_time": item["end"],
            "purpose": item["purpose"],
            "status": "confirmed",
        }
        for item in ROOM_BOOKINGS
    ]


def _prepare_event_registrations(event_map, user_map):
    return [
        {
            "event_id": event_map[item["event"]],
            "user_id": user_map[item["user"]],
            "status": item["status"],
        }
        for item in EVENT_REGISTRATIONS
    ]


def _prepare_ai_quizzes(course_map):
    quizzes = []
    for item in AI_QUIZZES:
        questions = [
            {
                "question": f"{item['prompt']} - core concept",
                "options": ["A", "B", "C", "D"],
                "answer": "A",
            }
        ]
        quizzes.append(
            {
                "course_id": course_map[item["course"]],
                "prompt": item["prompt"],
                "questions": questions,
            }
        )
    return quizzes


def _prepare_dean_requests(user_map):
    return [
        {
            "user_id": user_map[item["user"]],
            "request_type": item["request_type"],
            "payload": item["payload"],
            "status": item["status"],
        }
        for item in DEAN_REQUESTS
    ]


def _prepare_admission_applications(program_map):
    prepared = []
    for key, record in ADMISSION_APPLICATIONS:
        data = record.copy()
        program_key = data.pop("program", None)
        data["program_id"] = program_map.get(program_key) if program_key else None
        prepared.append((key, data))
    return prepared


def _prepare_admission_documents(app_map):
    return [
        {
            "application_id": app_map[item["application"]],
            "file_name": item["file_name"],
            "file_type": item["file_type"],
            "storage_url": item["storage_url"],
        }
        for item in ADMISSION_DOCUMENTS
    ]


def _prepare_teaching_attendance(user_map, course_map):
    rows = []
    for item in TEACHING_ATTENDANCE:
        attendance = [
            {"student_id": user_map[e["student"]], "present": e["present"]} for e in item["attendance"]
        ]
        rows.append(
            {
                "course_id": course_map[item["course"]],
                "professor_id": user_map[item["professor"]],
                "session_date": item["session_date"],
                "attendance": attendance,
            }
        )
    return rows


def _prepare_teaching_grade_uploads(user_map, course_map):
    rows = []
    for item in TEACHING_GRADE_UPLOADS:
        payload = [
            {"student_id": user_map[e["student"]], "grade": e["grade"]} for e in item["payload"]
        ]
        rows.append(
            {
                "course_id": course_map[item["course"]],
                "professor_id": user_map[item["professor"]],
                "payload": payload,
            }
        )
    return rows


def _prepare_teaching_submissions(user_map, course_map):
    return [
        {
            "course_id": course_map[item["course"]],
            "student_id": user_map[item["student"]],
            "title": item["title"],
            "status": item["status"],
            "grade": item.get("grade"),
        }
        for item in TEACHING_SUBMISSIONS
    ]


def _prepare_teaching_announcements(user_map, course_map):
    return [
        {
            "course_id": course_map[item["course"]],
            "professor_id": user_map[item["professor"]],
            "message": item["message"],
        }
        for item in TEACHING_ANNOUNCEMENTS
    ]


def _prepare_teaching_feedback(user_map, course_map):
    return [
        {
            "course_id": course_map[item["course"]],
            "student_id": user_map[item["student"]],
            "rating": item["rating"],
            "comment": item["comment"],
        }
        for item in TEACHING_FEEDBACK
    ]


def _prepare_vacation_requests(user_map):
    return [
        {
            "employee_id": user_map[item["employee"]],
            "start_date": item["start"],
            "end_date": item["end"],
            "vacation_type": item["vacation_type"],
            "status": item["status"],
        }
        for item in VACATION_REQUESTS
    ]


def _prepare_business_trips(user_map):
    return [
        {
            "employee_id": user_map[item["employee"]],
            "destination": item["destination"],
            "start_date": item["start"],
            "end_date": item["end"],
            "purpose": item["purpose"],
            "status": item["status"],
        }
        for item in BUSINESS_TRIPS
    ]


def _prepare_hr_certificates(user_map):
    return [
        {
            "employee_id": user_map[item["employee"]],
            "certificate_type": item["certificate_type"],
            "status": item["status"],
            "download_url": item.get("download_url"),
        }
        for item in HR_CERTIFICATES
    ]


def _prepare_dorm_rooms(user_map):
    return [
        {
            "student_id": user_map[item["student"]],
            "room_number": item["room_number"],
            "building": item["building"],
            "balance": item["balance"],
        }
        for item in DORM_ROOMS
    ]


def _prepare_dorm_requests(user_map):
    return [
        {
            "student_id": user_map[item["student"]],
            "request_type": item["request_type"],
            "description": item["description"],
            "status": item["status"],
        }
        for item in DORM_REQUESTS
    ]


def _prepare_dorm_payments(user_map):
    return [
        {
            "student_id": user_map[item["student"]],
            "amount": item["amount"],
            "reference": item["reference"],
        }
        for item in DORM_PAYMENTS
    ]


def _prepare_library_digital_assets(book_map):
    return [
        {
            "book_id": book_map[item["book"]],
            "format": item["format"],
            "access_url": item["access_url"],
            "metadata": item["metadata"],
        }
        for item in LIBRARY_DIGITAL_ASSETS
    ]


def _prepare_library_reservations(book_map, user_map):
    return [
        {
            "book_id": book_map[item["book"]],
            "student_id": user_map[item["student"]],
            "status": item["status"],
        }
        for item in LIBRARY_RESERVATIONS
    ]


def _prepare_library_loans(book_map, user_map):
    return [
        {
            "book_id": book_map[item["book"]],
            "student_id": user_map[item["student"]],
            "borrowed_at": item["borrowed_at"],
            "due_at": item["due_at"],
            "status": item["status"],
        }
        for item in LIBRARY_LOANS
    ]


def _prepare_support_queries(user_map):
    return [
        {
            "user_id": user_map[item["user"]],
            "question": item["question"],
            "answer": item["answer"],
        }
        for item in SUPPORT_QUERIES
    ]


def _prepare_support_tickets(user_map):
    return [
        {
            "user_id": user_map[item["user"]],
            "category": item["category"],
            "subject": item["subject"],
            "description": item["description"],
            "status": item["status"],
        }
        for item in SUPPORT_TICKETS
    ]


def _prepare_ai_advisor(user_map):
    return [
        {
            "user_id": user_map[item["user"]],
            "topic": item["topic"],
            "prompt": item["prompt"],
            "response": item["response"],
        }
        for item in AI_ADVISOR
    ]


def _prepare_admission_event_bookings(admission_event_map):
    return [
        {
            "event_id": admission_event_map[item["event"]],
            "applicant_name": item["applicant_name"],
            "email": item["email"],
            "phone": item.get("phone"),
            "note": item.get("note"),
            "status": "confirmed",
        }
        for item in ADMISSION_EVENT_BOOKINGS
    ]
