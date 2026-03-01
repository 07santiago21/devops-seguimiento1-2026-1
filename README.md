# DevOps - Seguimiento #1 (2026-1)

## Project Overview

This project consists of building a RESTful API using Go (net/http) with real database persistence

The application will be deployed in two independent environments:

- 🧪 Test Environment
- 🚀 Production Environment

Each environment will have:
- Independent deployment
- Independent database
- Independent environment variables
- Independent CI/CD pipeline

The project will include:
- Version control with Git (GitMoji convention)
- CI/CD pipelines
- Minimum test coverage requirements


---

## Tech Stack

- **Language:** Go
- **HTTP:** net/http
- **Database:** PostgreSQL
- **ORM:** GORM
- **CI/CD:** GitHub Actions (to be implemented)

---

## Entities

The API will manage the following entities:

### 1️⃣ Student
- id (UUID)
- name
- last_name
- age
- created_at

### 2️⃣ Course
- id (UUID)
- name
- description
- credits
- capacity
- created_at

### 3️⃣ Enrollment
- id (UUID)
- student_id (FK)
- course_id (FK)
- status (active / completed / cancelled)
- enrollment_date
- total_amount

## Environment Variables

Create a `.env` file in the root directory:

```env
DATABASE_HOST=
DATABASE_PORT=
DATABASE_USER=
DATABASE_PASSWORD=
DATABASE_NAME=