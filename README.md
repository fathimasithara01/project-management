# 🚀 Project Management API

> Production-ready backend system demonstrating **Clean Architecture**, authentication, and scalable API design, suitable for real-world applications.

A full-featured **Project Management System** built with **Golang (Gin), GORM, PostgreSQL, and JWT**, following **Clean Architecture principles**.  
It enables organizations to efficiently manage **users, projects, and tasks**, with **role-based access control, secure authentication**, and **scalable backend design**.

---

## 🏗️ Architecture

This project follows **Clean Architecture + Layered Design** for maintainability and scalability.

### 📦 Layers

- **Controller (Handler)** – Handles HTTP requests and responses  
- **Usecase (Service)** – Contains business logic and application rules  
- **Repository (Data Layer)** – Database interactions using GORM  
- **Middleware** – JWT authentication & role-based authorization  
- **Domain (Models)** – Core business entities  
- **Infrastructure** – Concrete implementations for external services (DB, etc.)  
- **Utils** – Helpers: validation, JWT, pagination, formatting  

### 🔄 Request Flow
Client → Controller → Usecase → Repository → Database


---

## 🧠 Features

- 🔐 JWT Authentication  
- 🛡️ Role-Based Access Control (Admin & Developer)  
- 👤 User Management  
- 📁 Project Management  
- ✅ Task Management  
- 📄 Pagination & Filtering  
- ⚙️ Clean Architecture  
- ✅ Proper HTTP Status Codes & Error Handling  
- ✅ Input Validation  
- 🔄 Database Transactions for critical operations  

---

## 🗄️ Database Design

### Entities

#### 👤 User

| Field       | Type     | Notes                   |
|------------|---------|------------------------|
| id         | uint    | Primary Key             |
| name       | string  | Not null               |
| email      | string  | Unique                 |
| password   | string  | Hashed, not exposed    |
| role       | string  | admin / developer       |

**Relationships:**

- One User → Many Projects (creator)  
- One User → Many Tasks (assignee)  

#### 📁 Project

| Field       | Type     | Notes                   |
|------------|---------|------------------------|
| id         | uint    | Primary Key             |
| name       | string  | Not null               |
| description| text    | Optional               |
| created_by | uint    | Foreign Key → User      |

**Relationships:**

- One Project → Many Tasks  

#### ✅ Task

| Field       | Type     | Notes                              |
|------------|---------|-----------------------------------|
| id         | uint    | Primary Key                        |
| title      | string  | Not null                           |
| description| text    | Optional                            |
| status     | string  | todo / in-progress / done           |
| project_id | uint    | Foreign Key → Project              |
| assigned_to| uint    | Foreign Key → User                 |
| due_date   | datetime| Optional                            |

### 🔗 ER Diagram

User (1) ---- (M) Project
User (1) ---- (M) Task
Project (1) ---- (M) Task


---

## ⚙️ Setup Instructions

### 1️⃣ Clone Repository

```bash
git clone https://github.com/fathimasithara01/project-management.git
cd project-management

## 2️⃣ Install Dependencies
go mod tidy

3️⃣ Environment Variables

Create a .env file based on .env.example:

PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=project_management
DB_SSLMODE=disable

JWT_SECRET=your_secret_key

4️⃣ Database Setup & Migrations
Ensure PostgreSQL is running
Create database: project_management
GORM AutoMigrate will create tables automatically on server start
Optional: manual SQL migrations in migrations/:
psql -U postgres -d project_management -f migrations/0001_create_users.sql
psql -U postgres -d project_management -f migrations/0002_create_projects.sql
psql -U postgres -d project_management -f migrations/0003_create_tasks.sql
5️⃣ Run the Server
go run cmd/server/main.go

Server will run at: http://localhost:8080

📡 API Documentation
🔐 Authentication
Method	Endpoint	Description
POST	/api/v1/auth/register	Register a new user
POST	/api/v1/auth/login	Login and get JWT token
👤 Users (Admin Only)
Method	Endpoint	Description
GET	/api/v1/users	List all users
GET	/api/v1/users/:id	Get user by ID
POST	/api/v1/users	Create new user
📁 Projects
Method	Endpoint	Description
POST	/api/v1/projects	Create project (Admin)
GET	/api/v1/projects	List all projects
GET	/api/v1/projects/:id	Get project by ID
PUT	/api/v1/projects/:id	Update project (Admin)
DELETE	/api/v1/projects/:id	Delete project (Admin)
✅ Tasks
Method	Endpoint	Role Access	Description
POST	/api/v1/tasks	Admin	Create task
GET	/api/v1/tasks	Admin/Developer	List tasks (pagination & filters)
GET	/api/v1/tasks/:id	Admin/Developer	Get task by ID
PUT	/api/v1/tasks/:id	Admin	Update task
PATCH	/api/v1/tasks/:id/status	Admin/Developer	Update task status
DELETE	/api/v1/tasks/:id	Admin	Delete task
🔍 Filtering & Pagination

Filtering:

/tasks?project_id=1
/tasks?status=todo
/tasks?assigned_to=2

Pagination:

/tasks?page=1&limit=10
📌 Sample Response
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Build API",
    "status": "in-progress"
  }
}

🔐 Security
JWT-based authentication
Role-Based Access Control (RBAC)
Middleware protects routes
Ownership validation ensures secure operations
🔄 Transactions

Critical operations like task assignment and status updates are wrapped in database transactions for data consistency and reliability.

🧪 Testing
APIs tested with Postman
Postman collection: docs/postman_collection.json
🛠️ Tech Stack
Language: Go
Framework: Gin
ORM: GORM
Database: PostgreSQL
Authentication: JWT
Architecture: Clean Architecture + Layered Design


🎯 Future Improvements
Docker + docker-compose support
Swagger / OpenAPI documentation
Unit & integration tests
CI/CD pipeline
Redis caching for performance

👩‍💻 Author

Fathima Sithara
Backend Engineer (Golang)