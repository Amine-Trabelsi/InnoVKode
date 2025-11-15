# PROJECT

- define user issues

## Target audience

- guest
- student
- employee

## Features

- after which we send an otp code to the server 
- auth by email email for students and employees
- visa renewal system (status, and all info about it)
- the bot architecture must be only managed by buttons
- book apoitment to submit documents (avoid queues)
- language support
- - Issue applications of different types (справка с места работы, 2ндфл для сотрудников, акадимический отпуст, отпус(оплачиваемый и не оплачиваемый) )
- store in DB if student a mark that the is a foreigner
- reminder for dorm payment (and info about dept, also make requests for repair)

---

- liabrary system

## scalibility

- think about scalibility for
  - 4 million users
  - 920 government universities
  - 340 commercial universities

- the system must support scalling horizontally, services handlling message can be spawn as many as needed
- think about scalability of how this solution can be generalized

## university tools

- make your own system over the LMS of the university

- Univresity Tools: LMS (Moodle, Canvas, Blackboard, WebTutor, 1C-Образование, Smart University)
- Schedule

---

User Usage Story — Telegram University Assistant Bot (Category-Based Flow with Guest)

---

1. Start & Category Selection
Input: /start
Output:
“Welcome to the University Assistant Bot! 🎓
Please select your category:”

- 🧑‍🎓 Student
- 👨‍🏫 Professor
- 🙋 Guest

---

### 🧑‍🎓 If Student selected:

2. Main Menu (Student)
Output:
“Hi, Student! What would you like to do?”

- 📅 My Schedule
- 🧾 Certificates
- 🧠 Lecture Summaries
- 🧩 Quizzes
- ❓ Ask about Lectures
- 💼 Career

3. Lecture Summary
Input: “🧠 Lecture Summaries → Physics → 12 Nov”
Output:
“Summary of Physics (12 Nov):

- Topic: Newton’s Laws
- Key Points: Inertia, F = ma, Action–Reaction
- Example: Car acceleration on frictionless surface.”

4. Quiz Generation
Input: “🧩 Generate quiz for Physics lecture.”
Output:
1️⃣ What does F = ma represent?
2️⃣ Which law describes inertia?
3️⃣ Example of Newton’s 3rd law?
→ User answers → “✅ 3/3 correct! Great job!”

5. RAG Q&A
Input: “❓ What’s the difference between static and kinetic friction?”
Output: “Static friction resists motion from rest; kinetic friction acts when sliding occurs. Static > kinetic.”

6. Certificates or Career
Input: “🧾 Request certificate.”
Output: “Confirm details? (2nd year, Faculty of Chemistry)” → ✅ Confirm
→ “✅ Request submitted. You’ll be notified when it’s ready.”

---

### 👨‍🏫 If Professor selected:

2. Main Menu (Professor)
Output:
“Hello, Professor! Choose an option:”

- 📤 Upload Lecture Notes
- 🧩 Generate Quiz from Lecture
- 🧠 Create Lecture Summary
- 💬 Answer Student Questions

3. Lecture Upload
Input: “📤 Upload Lecture Notes → ‘Thermodynamics.pdf’”
Output: “✅ Lecture uploaded. Students can now view summaries and take quizzes.”

4. Generate Quiz
Input: “🧩 Generate quiz from uploaded lecture.”
Output: “✅ Quiz created and linked to your Thermodynamics lecture.”

---

### 🙋 If Guest selected

1. TODO

2. Main Menu (Guest)
    Output:
    “Welcome, Guest! Explore information about our university:”

    - 🏛 About the University
    - 📚 Programs Offered
    - 📅 Admission Deadlines
    - 🗓 Campus Events
    - 📞 Contact & Support

3. Explore Info
    Input: “📚 Programs Offered”
    Output:
    “Available Programs:

    - B.Sc. in Computer Science
    - B.A. in Linguistics
    - M.Sc. in Data Science
    Tap any to learn more.”

4. Admission Info
Input: “📅 Admission Deadlines”
Output:
“Application deadlines:

- Bachelor’s: June 30
- Master’s: July 15
- International: May 31”

---

→ End Result:
Each category (Student, Professor, Guest) gets a tailored, intuitive experience:

- Students manage academic tasks, summaries, quizzes, and Q&A.

- Professors upload materials and generate assessments.

- Guests explore university info and admissions — all through one unified Telegram bot.
