# 🚀 Getting Started – Payroll System

---

## 🛠️ Requirements

- Go 1.20+
- PostgreSQL (recommended version 13+)
- Git
- Make (optional, for scripts)

---

## 📦 Installation

1. **Clone the repository**

```bash
git
cd
```

2. **Set up environment variables**

Create a `.env` file in the root:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=payroll_system_db
JWT_SECRET=your_super_secret_key
```

3. **Install Go dependencies**

```bash
go mod tidy
```

---

## 🧑‍💻 Running Locally

### 1. Create PostgreSQL database

```bash

```

### 2. Run the application

```bash

```

Your server will start on:
👉 `http://localhost:8080`

---

## 🧪 Running Tests

### Unit & Integration Tests

```bash
go test ./...
```

---

## 🛠 Project Structure

```
.
├── cmd
│   ├── app             # Application entry point
│   └── script          # Scripts entry point
├── docs                # Swagger docs
├── internal
│   ├── db              # DB setup
│   ├── dto             # Request/response schemas
│   ├── handlers        # Route handlers
│   ├── middlewares     # Middlewares (Auth)
│   ├── models          # GORM models
│   ├── route           # Route definitions
│   ├── seed            # DB seed helpers
│   └── utils           # Helpers (e.g., JWT, response formatting)
```

---

## 🧠 Credentials & Helpers

```sql

```

---

# 📘 Payroll System API Documentation

---

## 📖 API Docs (Swagger UI)

For full API details, refer to the Swagger spec or view them interactively in Swagger UI.

After running the app:

- Visit: `http://localhost:8080/swagger/index.html`

Or regenerate docs with:

```bash
swag init --parseDependency --parseInternal --d cmd/app/,internal/handlers
```

---

## 🔐 Authentication

### `POST /auth/login`

Authenticates a user using their username and password.

#### Request Body

```json
{
  "username": "johndoe",
  "password": "your_password"
}
```

#### Response

```json
{
  "message": "success",
  "data": {
    "token": "jwt_token_here",
    "user": {
      "id": 1,
      "username": "johndoe",
      "role": "Employee"
    }
  }
}
```

---

## 🔐 Authorization

All routes (except `/auth/login`) require authentication using a **Bearer token** passed in the request header.

### 🔑 Required Header

```http
Authorization: Bearer <JWT_TOKEN>
```

- Replace `<JWT_TOKEN>` with the token received from the `/auth/login` endpoint.
- This must be included in the **`Authorization`** header of every protected request.

---

## 👤 Attendance

### `POST /attendances/check-in`

Check-in for the current day.

- Only allowed once per day.
- Not allowed on weekends.

#### Response (200 OK)

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "date": "2025-06-01",
    "check_in_at": "2025-06-01T09:00:00Z",
    "check_out_at": null
  }
}
```

---

### `POST /attendances/check-out`

Check-out for the current day.

- Must have checked in first.
- Only allowed once per day.
- Not allowed on weekends.

#### Response (200 OK)

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "date": "2025-06-01",
    "check_in_at": "2025-06-01T09:00:00Z",
    "check_out_at": "2025-06-01T17:00:00Z"
  }
}
```

---

### `POST /attendances/overtime`

Submit overtime for the current day.

- Must be submitted **after check-out**
- Max **3 hours** allowed per day

#### Request Body

```json
{
  "hours_worked": 2.5
}
```

#### Response (200 OK)

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "hours_worked": 2.5
  }
}
```

---

## 💵 Reimbursements

### `POST /reimbursements`

Submit a reimbursement request.

#### Request Body

```json
{
  "amount": 100000,
  "description": "Taxi to client site"
}
```

#### Response (201 Created)

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "amount": 100000,
    "description": "Taxi to client site"
  }
}
```

---

## 🧮 Payroll

### `POST /payrolls/{year}/{month}`

Create or update a payroll for the specified year and month.

#### Request Body (optional)

```json
{
  "name": "June 2025 Payroll",
  "period_start": "2025-06-01",
  "period_end": "2025-06-30"
}
```

#### Response (200 OK)

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "name": "June 2025 Payroll",
    "period_start": "2025-06-01",
    "period_end": "2025-06-30",
    "status": "draft"
  }
}
```

---

### `POST /payrolls/{year}/{month}/run`

Run payroll generation for the specified year and month.

- Generates payslips for all employees.
- Changes status to `draft` → `pending`.
- Status will automatically change from `pending` → `processed` once the background task completes.
- Can only be run once per payroll.

#### Response

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "name": "June 2025 Payroll",
    "period_start": "2025-06-01",
    "period_end": "2025-06-30",
    "status": "pending"
  }
}
```

---

## 🧾 Payslip

### `GET /payslips/{year}/{month}`

Get the payslip for the authenticated user for the specified period.

#### Response (200 OK)

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 2,
    "month": 6,
    "year": 2025,
    "base_salary": 4200000,
    "overtime_pay": 200000,
    "reimbursement": 100000,
    "total_salary": 4500000,
    "total_hours_worked": 160,
    "total_overtime_hours": 10,
    "attendance_breakdown": "...",
    "overtime_breakdown": "...",
    "reimbursement_breakdown": "..."
  }
}
```
