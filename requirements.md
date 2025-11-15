# **MAX Bot - Complete Feature Structure**

## **🎯 Target Audiences (4 Roles)**

1. **Applicants (Абитуриенты)** - Prospective students
2. **Students (Студенты)** - Enrolled students
3. **Employees (Сотрудники вуза)** - University staff
4. **Leadership (Руководители вуза)** - University administration

---

## **🔐 Authentication System**

- **Guest Mode (Applicants)**
  - No authentication required
  - View-only access to public information

- **Authenticated Users (Students, Employees, Leadership)**
  - Email-based authentication
  - OTP code sent to email
  - User data stored in DB:
    - Name (Russian)
    - Name (English)
    - Role (Student/Employee/Leadership)
    - Foreign status (boolean flag)
    - Dorm info (if applicable)
    - Department/Faculty

---

## **🌍 Language Support**

- **Initial Setup**
  - Language selection on first interaction
  - Options: Russian 🇷🇺 | English 🇬🇧
  - Stored in user preferences

- **Language Switcher**
  - Available in main menu at any time (button)

---

## **📱 Main Menu Structure (Role-Based)**

### **1️⃣ APPLICANTS (Абитуриенты)** - Guest Mode

```txt
🏠 Main Menu
│
├── 📚 Admission (Поступление)
│   ├── ℹ️ About University
│   │   ├── Programs & Faculties
│   │   ├── Campus Information
│   │   ├── Student Life
│   │   └── Admission Requirements
│   │
│   ├── 📅 Open Day Registration
│   │   ├── View Available Dates
│   │   ├── Book Open Day Visit
│   │   │   ├── Select Date
│   │   │   ├── Select Time Slot
│   │   │   ├── Enter Contact Info (Name, Email, Phone)
│   │   │   └── ✅ Confirm Booking
│   │   └── My Bookings
│   │
│   ├── 🏛️ Campus Tour
│   │   ├── View Available Tours
│   │   ├── Book Tour
│   │   │   ├── Select Date
│   │   │   ├── Select Time
│   │   │   ├── Group/Individual
│   │   │   └── ✅ Confirm Booking
│   │   └── My Tour Bookings
│   │
│   └── 📞 Contact Admissions Office
│       ├── Phone Numbers
│       ├── Email Addresses
│       └── Office Hours & Location
│
├── 📄 Admission Documents
│   ├── Required Documents List
│   ├── Document Submission Deadline
│   └── 📅 Book Appointment (Submit Documents)
│       ├── Select Date
│       ├── Select Time Slot
│       ├── Enter Applicant Info
│       └── ✅ Confirm Appointment
│
└── 🌐 Language Settings
    └── Switch Language (RU/EN)
```

---

### **2️⃣ STUDENTS (Студенты)** - Authenticated

```txt
🏠 Main Menu
│
├── 📚 Education (Обучение)
│   ├── 📅 Schedule
│   │   ├── This Week's Schedule
│   │   ├── Next Week's Schedule
│   │   ├── Full Semester Schedule
│   │   └── Exam Schedule
│   │
│   ├── 💬 Teacher Feedback
│   │   ├── View My Courses
│   │   ├── Select Course
│   │   ├── Rate Teacher (1-5 stars)
│   │   ├── Leave Comment (Optional)
│   │   └── ✅ Submit Feedback
│   │
│   └── ➕ Electives Registration
│       ├── Browse Available Electives
│       ├── Filter by Category/Department
│       ├── View Elective Details
│       │   ├── Course Description
│       │   ├── Teacher Info
│       │   ├── Schedule & Location
│       │   └── Available Spots
│       └── ✅ Register for Elective
│
├── 🚀 Project Activities (Проектная деятельность)
│   ├── 💡 Submit My Project
│   │   ├── Project Title
│   │   ├── Description
│   │   ├── Required Team Size
│   │   ├── Skills Needed
│   │   └── ✅ Submit Project
│   │
│   ├── 👥 Build Team
│   │   ├── My Project
│   │   ├── Team Requests
│   │   ├── Accept/Reject Members
│   │   └── Current Team Members
│   │
│   ├── 🔍 Browse Projects
│   │   ├── Filter by Category
│   │   ├── View Project Details
│   │   └── ✅ Apply to Join
│   │
│   └── 📋 My Projects
│       ├── Projects I Created
│       ├── Projects I Joined
│       └── Task Notifications
│           ├── New Tasks Assigned
│           ├── Deadlines Approaching
│           └── Task Updates
│
├── 💼 Career (Карьера)
│   ├── 📞 Career Center Consultation
│   │   ├── Book Consultation
│   │   │   ├── Select Date
│   │   │   ├── Select Topic (CV, Interview, Career Path)
│   │   │   └── ✅ Confirm Booking
│   │   └── My Consultations
│   │
│   └── 💼 Job Board
│       ├── Browse Vacancies
│       ├── Filter by Category/Industry
│       ├── View Job Details
│       │   ├── Company Info
│       │   ├── Requirements
│       │   ├── Salary & Benefits
│       │   └── Application Deadline
│       ├── ✅ Apply for Job
│       └── My Applications
│
├── 🏛️ Dean's Office (Деканат)
│   ├── 📄 Request Certificates
│   │   ├── Certificate Types:
│   │   │   ├── Certificate of Enrollment (Справка об обучении)
│   │   │   ├── Transcript (Академическая справка)
│   │   │   ├── Scholarship Certificate (Справка о стипендии)
│   │   │   └── Other Documents
│   │   ├── Select Certificate Type
│   │   ├── Specify Purpose (if needed)
│   │   ├── Delivery Method (Email PDF / Pick up)
│   │   └── ✅ Submit Request
│   │   └── Track Request Status
│   │
│   ├── 💳 Tuition Payment
│   │   ├── View Balance
│   │   ├── Payment History
│   │   └── 💰 Pay Now (External Link)
│   │
│   ├── 💵 Apply for Compensation
│   │   ├── Compensation Types
│   │   ├── Fill Application Form
│   │   ├── Upload Documents
│   │   └── ✅ Submit Application
│   │   └── Track Status
│   │
│   ├── 📅 Book Appointment (Dean's Office)
│   │   ├── Select Service (Documents, Transfer, Leave)
│   │   ├── Select Date
│   │   ├── Select Time Slot
│   │   └── ✅ Confirm Appointment
│   │
│   └── 📝 Submit Application
│       ├── Transfer Application (Перевод)
│       │   ├── Transfer Type (Faculty/Program/University)
│       │   ├── Reason
│       │   ├── Upload Supporting Docs
│       │   └── ✅ Submit
│       └── Academic Leave (Академический отпуск)
│           ├── Leave Type (Medical/Personal/Military)
│           ├── Duration (From - To)
│           ├── Reason
│           ├── Upload Supporting Docs (Medical cert, etc.)
│           └── ✅ Submit
│           └── Track Status
│
├── 🏠 Dormitory (Общежитие)
│   ├── 💰 Payment
│   │   ├── Check Balance
│   │   │   ├── Current Balance
│   │   │   ├── Debt Amount (if any)
│   │   │   ├── Next Payment Due Date
│   │   │   └── Payment History
│   │   └── 💳 Pay Now (External Link - from config)
│   │
│   ├── 🛎️ Additional Services
│   │   ├── Browse Services (Laundry, Cleaning, etc.)
│   │   ├── View Prices
│   │   ├── Select Service
│   │   └── ✅ Order Service
│   │
│   ├── 🎫 Guest Pass
│   │   ├── Request Guest Pass
│   │   │   ├── Guest Name
│   │   │   ├── Guest ID/Passport
│   │   │   ├── Visit Date
│   │   │   ├── Duration (Hours)
│   │   │   └── ✅ Submit Request
│   │   └── My Guest Passes
│   │
│   └── 🔧 Maintenance Request
│       ├── Submit Repair Request
│       │   ├── Issue Type (Plumbing, Electrical, Furniture, etc.)
│       │   ├── Description
│       │   ├── Urgency (Low/Medium/High)
│       │   ├── Upload Photo (Optional)
│       │   └── ✅ Submit
│       └── Track Requests
│           ├── Open Requests
│           ├── In Progress
│           └── Completed
│
├── 🎭 Extracurricular (Внеучебная деятельность)
│   ├── 📅 Events Calendar
│   │   ├── View All Events
│   │   ├── Filter by Category (Sports, Arts, Academic, Social)
│   │   ├── Filter by Date
│   │   └── View Event Details
│   │       ├── Event Name
│   │       ├── Date & Time
│   │       ├── Location
│   │       ├── Description
│   │       ├── Available Spots
│   │       └── Registration Status
│   │
│   ├── ✅ Register for Event
│   │   ├── As Attendee
│   │   └── As Participant
│   │
│   └── 📋 My Events
│       ├── Upcoming Events
│       ├── Past Events
│       └── Cancel Registration
│
├── 📚 Library (Библиотека)
│   ├── 🔍 Search Books
│   │   ├── Search by Title
│   │   ├── Search by Author
│   │   ├── Search by ISBN
│   │   └── Browse by Category
│   │
│   ├── 📖 Order Physical Book
│   │   ├── Select Book
│   │   ├── Choose Pickup Location
│   │   ├── ✅ Reserve Book
│   │   └── Track Reservation
│   │
│   ├── 💻 E-Library Access
│   │   ├── Browse E-Books
│   │   ├── Access E-Journals
│   │   └── 🔗 Open E-Library Portal (External Link)
│   │
│   └── 📋 My Library
│       ├── Books on Loan
│       ├── Due Dates
│       ├── Reservations
│       └── Fines (if any)
│
├── 🛂 Visa Services (For Foreign Students)
│   ├── 📋 My Visa Status
│   │   ├── Current Visa Type
│   │   ├── Issue Date
│   │   ├── Expiration Date
│   │   ├── Days Until Expiration
│   │   └── Visa Status (Valid/Expiring Soon/Expired)
│   │
│   ├── 🔄 Visa Renewal
│   │   ├── View Renewal Process
│   │   ├── Required Documents Checklist
│   │   ├── Submit Renewal Application
│   │   │   ├── Upload Required Documents
│   │   │   ├── Fill Application Form
│   │   │   └── ✅ Submit
│   │   └── Track Renewal Status
│   │       ├── Application Received
│   │       ├── Under Review
│   │       ├── Approved/Rejected
│   │       └── Ready for Pickup
│   │
│   └── 📅 Book Visa Office Appointment
│       ├── Select Service Type
│       ├── Select Date
│       ├── Select Time
│       └── ✅ Confirm
│
├── ⚙️ Settings
│   ├── 👤 My Profile
│   │   ├── Name (RU)
│   │   ├── Name (EN)
│   │   ├── Email
│   │   ├── Faculty/Department
│   │   ├── Year of Study
│   │   ├── Student ID
│   │   ├── Dorm Info (if applicable)
│   │   └── Foreign Student Status
│   │
│   ├── 🌐 Language
│   │   └── Switch Language (RU/EN)
│   │
│   └── 🔔 Notifications
│       ├── Enable/Disable Notifications
│       └── Notification Preferences
│
└── ℹ️ Help & Support
    ├── FAQ
    ├── Contact Support
    └── Report Issue
```

---

### **3️⃣ EMPLOYEES (Сотрудники вуза)** - Authenticated

```txt
🏠 Main Menu
│
├── ✈️ Business Trips (Командировки)
│   ├── ➕ Submit Trip Request
│   │   ├── Destination
│   │   ├── Travel Dates (From - To)
│   │   ├── Purpose/Conference Name
│   │   ├── Estimated Budget
│   │   ├── Upload Invitation (if any)
│   │   └── ✅ Submit Request
│   │
│   ├── 📋 My Trip Requests
│   │   ├── Pending Approval
│   │   │   ├── View Status
│   │   │   ├── Approval Progress (Dept Head → Finance → Rector)
│   │   │   └── Estimated Processing Time
│   │   ├── Approved Trips
│   │   └── Rejected Trips (with reason)
│   │
│   ├── 📊 Submit Trip Report
│   │   ├── Select Trip
│   │   ├── Upload Receipts (Multiple)
│   │   │   ├── OCR Automatic Processing
│   │   │   ├── Taxi
│   │   │   ├── Accommodation
│   │   │   ├── Meals
│   │   │   └── Conference Fee
│   │   ├── Total Expenses
│   │   ├── Budget vs Actual
│   │   └── ✅ Submit Report
│   │
│   └── 📅 Book Travel Office Appointment
│       ├── Select Date
│       ├── Select Time
│       └── ✅ Confirm
│
├── 🏖️ Vacation (Отпуска)
│   ├── ➕ Submit Vacation Request
│   │   ├── Vacation Type
│   │   │   ├── Paid Vacation (Оплачиваемый отпуск)
│   │   │   └── Unpaid Leave (Неоплачиваемый отпуск)
│   │   ├── Start Date
│   │   ├── End Date
│   │   ├── Duration (Auto-calculated working days)
│   │   ├── Reason (if unpaid)
│   │   ├── Substitute Assignment
│   │   │   ├── Select Colleague
│   │   │   └── Specify Responsibilities
│   │   └── ✅ Submit Request
│   │
│   ├── 📊 My Vacation Balance
│   │   ├── Total Days Allocated
│   │   ├── Days Used
│   │   ├── Days Pending Approval
│   │   ├── Days Available
│   │   └── Days Expiring Soon
│   │
│   ├── 📋 My Vacation Requests
│   │   ├── Pending Approval
│   │   │   ├── View Status
│   │   │   ├── Approval Progress
│   │   │   └── Expected Decision Date
│   │   ├── Approved Vacations
│   │   │   └── Add to Calendar
│   │   └── Rejected Requests (with reason)
│   │
│   └── 👥 Team Calendar
│       ├── View Team Vacations
│       ├── Filter by Date Range
│       └── Check Availability
│
├── 🏢 Office Services (Офис)
│   ├── 📄 Request Certificates
│   │   ├── Certificate Types:
│   │   │   ├── Employment Certificate (Справка с места работы)
│   │   │   ├── Tax Form (2-НДФЛ)
│   │   │   ├── Salary Certificate
│   │   │   └── Other Documents
│   │   ├── Select Certificate Type
│   │   ├── Purpose (Visa, Bank, etc.)
│   │   ├── Delivery Method
│   │   │   ├── Email (PDF with digital signature)
│   │   │   └── Pick up from HR (Paper with stamp)
│   │   ├── Expected Processing Time
│   │   └── ✅ Submit Request
│   │   └── Track Request Status
│   │
│   ├── 🎫 Guest Pass (Office)
│   │   ├── Request Guest Pass
│   │   │   ├── Guest Name
│   │   │   ├── Guest ID/Passport
│   │   │   ├── Visit Date
│   │   │   ├── Time (From - To)
│   │   │   ├── Building Access Required
│   │   │   ├── Parking Needed (Yes/No)
│   │   │   └── ✅ Submit Request
│   │   ├── Guest Receives SMS/Email with QR Code
│   │   └── My Guest Passes
│   │       ├── Upcoming Visits
│   │       └── Past Visits
│   │
│   └── 📅 Book HR Office Appointment
│       ├── Select Service (Documents, Questions, etc.)
│       ├── Select Date
│       ├── Select Time
│       └── ✅ Confirm
│
├── 🎭 Extracurricular (Внеучебная деятельность)
│   ├── 📅 Events Calendar
│   │   ├── View All Events
│   │   ├── Filter by Category (Sports, Arts, Social)
│   │   ├── Filter by Date
│   │   └── View Event Details
│   │
│   ├── ✅ Register for Event
│   │   ├── As Attendee
│   │   └── As Participant
│   │
│   └── 📋 My Events
│       ├── Upcoming Events
│       ├── Past Events
│       └── Cancel Registration
│
├── 🛂 Visa Services (For Foreign Employees)
│   ├── 📋 My Visa Status
│   │   ├── Current Visa Type
│   │   ├── Issue Date
│   │   ├── Expiration Date
│   │   ├── Days Until Expiration
│   │   └── Visa Status (Valid/Expiring Soon/Expired)
│   │
│   ├── 🔄 Visa Renewal
│   │   ├── View Renewal Process
│   │   ├── Required Documents Checklist
│   │   ├── Submit Renewal Application
│   │   │   ├── Upload Required Documents
│   │   │   ├── Fill Application Form
│   │   │   └── ✅ Submit
│   │   └── Track Renewal Status
│   │
│   └── 📅 Book Visa Office Appointment
│       ├── Select Service Type
│       ├── Select Date
│       ├── Select Time
│       └── ✅ Confirm
│
├── ⚙️ Settings
│   ├── 👤 My Profile
│   │   ├── Name (RU)
│   │   ├── Name (EN)
│   │   ├── Email
│   │   ├── Department/Faculty
│   │   ├── Position/Title
│   │   ├── Employee ID
│   │   └── Foreign Employee Status
│   │
│   ├── 🌐 Language
│   │   └── Switch Language (RU/EN)
│   │
│   └── 🔔 Notifications
│       └── Notification Preferences
│
└── ℹ️ Help & Support
    ├── FAQ
    ├── Contact Support
    └── Report Issue
```

---

### **4️⃣ LEADERSHIP (Руководители вуза)** - Authenticated

```txt
🏠 Main Menu
│
├── 📰 News Aggregator (Агрегатор новостной ленты)
│   ├── 📊 News Feed
│   │   ├── All Mentions (University name in media)
│   │   ├── Filter by Source
│   │   │   ├── News Websites
│   │   │   ├── Social Media
│   │   │   ├── Academic Publications
│   │   │   └── Local Media
│   │   ├── Filter by Date Range
│   │   └── View Article Details
│   │       ├── Article Title
│   │       ├── Source
│   │       ├── Date Published
│   │       ├── Summary
│   │       └── 🔗 Full Article Link
│   │
│   ├── 📈 Sentiment Analysis (Optional for v2)
│   │   ├── Positive Mentions
│   │   ├── Neutral Mentions
│   │   └── Negative Mentions
│   │
│   └── 🔔 News Alerts
│       ├── Enable/Disable Alerts
│       └── Alert Frequency (Real-time/Daily Digest)
│
├── 🎭 Extracurricular (Внеучебная деятельность)
│   ├── 📅 Events Calendar
│   │   ├── View All Events
│   │   ├── Filter by Category
│   │   ├── Filter by Date
│   │   └── View Event Details
│   │
│   ├── ✅ Register for Event
│   │   ├── As Attendee
│   │   └── As Participant
│   │
│   └── 📋 My Events
│       ├── Upcoming Events
│       └── Past Events
│
├── ⚙️ Settings
│   ├── 👤 My Profile
│   │   ├── Name (RU)
│   │   ├── Name (EN)
│   │   ├── Email
│   │   ├── Position/Title
│   │   └── Department
│   │
│   ├── 🌐 Language
│   │   └── Switch Language (RU/EN)
│   │
│   └── 🔔 Notifications
│       └── Notification Preferences
│
└── ℹ️ Help & Support
    ├── FAQ
    ├── Contact Support
    └── Report Issue
```

---

## **🗄️ Database Schema Requirements**

### **Users Table**

```txt
- user_id (primary key)
- max_user_id (MAX messenger ID)
- email
- role (applicant/student/employee/leadership)
- name_ru
- name_en
- is_foreigner (boolean)
- language_preference (ru/en)
- created_at
- last_active
```

### **Students Table** (extends Users)

```txt
- student_id (foreign key)
- faculty
- department
- year_of_study
- student_number
- dorm_building (nullable)
- dorm_room (nullable)
- dorm_payment_due_date
- dorm_balance
```

### **Employees Table** (extends Users)

```txt
- employee_id (foreign key)
- department
- position
- employee_number
- hire_date
```

### **Visa Information Table**

```txt
- visa_id (primary key)
- user_id (foreign key)
- visa_type
- issue_date
- expiration_date
- renewal_status (none/pending/approved/rejected)
- documents_uploaded (JSON)
```

### **Appointments Table**

```txt
- appointment_id
- user_id
- appointment_type (dean/hr/admissions/visa)
- date
- time_slot
- status (pending/confirmed/cancelled/completed)
- notes
```

### **Applications Table**

```txt
- application_id
- user_id
- application_type (academic_leave/transfer/vacation/business_trip/certificate)
- status (pending/approved/rejected)
- submitted_date
- approval_chain (JSON)
- documents (JSON)
```

### **Events Table**

```txt
- event_id
- title
- description
- category
- date_time
- location
- max_attendees
- current_attendees
- registration_type (attendee/participant)
```

### **Event Registrations Table**

```txt
- registration_id
- event_id
- user_id
- registration_type
- status (registered/cancelled)
```

---

## **🔧 Technical Implementation Notes**

### **Button Navigation Structure**

- Main menu with role-based buttons
- Breadcrumb navigation (Back buttons at each level)
- Inline keyboards for selections
- Reply keyboards for frequently used actions

### **External Links (from Configuration)**

- Dorm payment URL
- E-library portal URL
- Tuition payment URL
- University website links

### **OTP Authentication Flow**

```txt
1. User sends /start
2. Bot: "Select language: RU | EN"
3. Bot: "Enter your email"
4. User enters email
5. Bot sends OTP to email
6. Bot: "Enter the code sent to your email"
7. User enters OTP
8. Bot validates OTP
9. Bot loads user profile from DB (role, name, etc.)
10. Bot displays role-specific main menu
```

### **Appointment Booking System**

- Time slots generated from configuration
- Conflict detection (no double-booking)
- Confirmation messages
- Reminder notifications (optional for v2)

### **Document Upload**

- Support for PDF, JPG, PNG
- File size limits
- OCR for receipt processing (business trips)

---

## **📋 Summary of Key Features**

✅ **Must-Have Features (All Roles):**

- Language selection (RU/EN)
- Button-only navigation
- Role-based access control

✅ **Applicants:**

- University information
- Open day booking
- Campus tour booking
- Document submission appointment

✅ **Students:**

- Schedule viewing
- Teacher feedback
- Project management
- Career services
- Dean's office services (certificates, applications)
- Dorm management (payment, guest pass, maintenance)
- Library services
- Extracurricular events

✅ **Employees:**

- Business trip management
- Vacation requests (paid/unpaid)
- Certificate requests (справка с места работы, 2-НДФЛ)
- Guest pass (office)
- Extracurricular events

✅ **Leadership:**

- News aggregator (university mentions)
- Extracurricular events

✅ **Foreign Users (Students & Employees):**

- Visa status tracking
- Visa renewal system
- Document checklist
- Appointment booking

✅ **Technical Features:**

- Email + OTP authentication
- Persistent storage (user profiles, foreigner flag, dorm info)
- External payment links (configurable)
- Appointment booking system
- Document upload support
