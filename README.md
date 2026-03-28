#  Project Management API

> Production-ready backend system demonstrating **Clean Architecture**, authentication, and scalable API design.

A full-featured **Project Management System** built with **Golang (Gin), GORM, PostgreSQL, and JWT**, following **Clean Architecture principles**.
It enables organizations to efficiently manage **users, projects, and tasks** with **secure authentication, role-based access control**, and a **scalable backend architecture**.

---

##  Architecture

This project follows **Clean Architecture + Layered Design** to ensure maintainability, scalability, and separation of concerns.

###  Layers

* **Handler (Controller)** – Handles HTTP requests and responses
* **Usecase (Service)** – Contains business logic and application rules
* **Repository (Data Layer)** – Handles database operations using GORM
* **Middleware** – JWT authentication & role-based authorization
* **Domain (Entities)** – Core business models
* **Infrastructure** – External services (DB, config, etc.)
* **Utils** – Helpers (JWT, validation, pagination, responses)

###  Request Flow

Client → Handler → Usecase → Repository → Database

---

##  Features

*  JWT Authentication
*  Role-Based Access Control (Admin & Developer)
*  User Management
*  Project Management
*  Task Management
*  Pagination & Filtering
*  Clean Architecture
*  Input Validation
*  Proper HTTP Status Codes
*  Database Transactions for critical operations

---

##  Database Design

### Entities

####  User

| Field    | Type   | Notes                |
| -------- | ------ | -------------------- |
| id       | uint   | Primary Key          |
| name     | string | Not null             |
| email    | string | Unique               |
| password | string | Hashed (not exposed) |
| role     | string | admin / developer    |

**Relationships:**

* One User → Many Projects (creator)
* One User → Many Tasks (assignee)

---

####  Project

| Field       | Type   | Notes              |
| ----------- | ------ | ------------------ |
| id          | uint   | Primary Key        |
| name        | string | Not null           |
| description | text   | Optional           |
| created_by  | uint   | Foreign Key → User |

**Relationships:**

* One Project → Many Tasks

---

####  Task

| Field       | Type     | Notes                     |
| ----------- | -------- | ------------------------- |
| id          | uint     | Primary Key               |
| title       | string   | Not null                  |
| description | text     | Optional                  |
| status      | string   | todo / in-progress / done |
| project_id  | uint     | Foreign Key → Project     |
| assigned_to | uint     | Foreign Key → User        |
| due_date    | datetime | Optional                  |

---

###  ER Diagram

```
User (1) ---- (M) Project
User (1) ---- (M) Task
Project (1) ---- (M) Task
```

---

##  Setup Instructions

### 1️ Clone Repository

```bash
git clone https://github.com/fathimasithara01/project-management.git
cd project-management
```

### 2 Install Dependencies

```bash
go mod tidy
```

---

### 3️ Environment Variables

Create a `.env` file:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=project_management
DB_SSLMODE=disable

JWT_SECRET=your_secret_key
```

---

### 4️ Database Setup

* Ensure PostgreSQL is running
* Create database:

```sql
CREATE DATABASE project_management;
```

* GORM AutoMigrate will create tables automatically on startup

(Optional manual migrations)

```bash
psql -U postgres -d project_management -f migrations/0001_create_users.sql
psql -U postgres -d project_management -f migrations/0002_create_projects.sql
psql -U postgres -d project_management -f migrations/0003_create_tasks.sql
```

---

### 5️ Run the Server

```bash
go run cmd/server/main.go
```

Server runs at:
 http://localhost:8080

---

## 📡 API Documentation

###  Authentication

| Method | Endpoint              | Description           |
| ------ | --------------------- | --------------------- |
| POST   | /api/v1/auth/register | Register new user     |
| POST   | /api/v1/auth/login    | Login & get JWT token |

---

###  Users (Admin Only)

| Method | Endpoint          | Description    |
| ------ | ----------------- | -------------- |
| GET    | /api/v1/users     | List users     |
| GET    | /api/v1/users/:id | Get user by ID |
| POST   | /api/v1/users     | Create user    |

---

###  Projects

| Method | Endpoint             | Description            |
| ------ | -------------------- | ---------------------- |
| POST   | /api/v1/projects     | Create project (Admin) |
| GET    | /api/v1/projects     | List projects          |
| GET    | /api/v1/projects/:id | Get project by ID      |
| PUT    | /api/v1/projects/:id | Update project (Admin) |
| DELETE | /api/v1/projects/:id | Delete project (Admin) |

---

###  Tasks

| Method | Endpoint                 | Access          | Description        |
| ------ | ------------------------ | --------------- | ------------------ |
| POST   | /api/v1/tasks            | Admin           | Create task        |
| GET    | /api/v1/tasks            | Admin/Developer | List tasks         |
| GET    | /api/v1/tasks/:id        | Admin/Developer | Get task by ID     |
| PUT    | /api/v1/tasks/:id        | Admin           | Update task        |
| PATCH  | /api/v1/tasks/:id/status | Admin/Developer | Update task status |
| DELETE | /api/v1/tasks/:id        | Admin           | Delete task        |

---

##  Filtering & Pagination

### Filtering

```
/tasks?project_id=1
/tasks?status=todo
/tasks?assigned_to=2
```

### Pagination

```
/tasks?page=1&limit=10
```

---

##  Sample Response

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Build API",
    "status": "in-progress"
  }
}
```

---

##  Security

* JWT-based authentication
* Role-Based Access Control (RBAC)
* Protected routes via middleware
* Ownership validation for secure operations

---

##  Transactions

Critical operations (e.g., task assignment, updates) are handled using **database transactions** to ensure consistency.

---

##  Testing

* APIs tested using Postman
* Collection available at:
  `docs/postman_collection.json`

---

##  Tech Stack

* **Language:** Go
* **Framework:** Gin
* **ORM:** GORM
* **Database:** PostgreSQL
* **Authentication:** JWT
* **Architecture:** Clean Architecture

---

##  Future Improvements

* Docker & docker-compose support
* Swagger / OpenAPI documentation
* Unit & integration tests
* CI/CD pipeline
* Redis caching

---

##  Author

**Fathima Sithara**
Backend Engineer (Golang)
