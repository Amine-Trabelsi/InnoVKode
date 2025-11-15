# Seed Data Overview

The `app.seed_data` module now seeds PostgreSQL with a coherent snapshot that mirrors the API surface described in `req.md` and `api.md`. Records are inserted by natural keys (students, courses, rooms, events, etc.), so every foreign key uses the IDs returned by the database instead of hard-coded numbers. By default, the app drops and recreates all tables (`RESET_DB_ON_STARTUP=true`) before reseeding so the dataset is always consistent; set that env var to `false` if you need to preserve live data.

Highlights:

- **People & Auth**: Eight core users covering students, employees, leadership, and applicants; students include dorm assignments, foreign-status flags, and course enrollments.
- **Academics**: Three courses with sessions, exams, grades, deadlines, notifications, attendance, submissions, announcements, and feedback so `/schedule`, `/exams`, `/grades`, and `/teaching/*` all return data.
- **Campus Life**: Rooms, bookings, events with RSVPs, clubs, news, dorm rooms/requests/payments, plus HR vacation/trip/certificate workflows.
- **Admissions**: Programs, open-day events, two applications with uploaded documents, and cached FAQ interactions for the applicant endpoints.
- **AI & Support**: Seeded RAG sources, queries, quizzes, summaries, transcriptions, advisor chats, and support tickets/queries.
- **Library & Facilities**: Physical/digital books, reservations, loans, and maintenance tickets, ensuring `/library/*` and `/dorms/*` respond with meaningful payloads.

The seed runs once during startup (after schema creation) and skips tables that already contain rows, making redeployments idempotent.
