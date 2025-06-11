# 🚀 Getting Started – Payroll System

---

## 🛠️ Requirements

- Go 1.20+
- PostgreSQL (recommended version 13+)
- Docker (optional, for running integrations tests)

---

## 📦 Installation

1. **Clone the repository**

```bash
git clone git@github.com:kristra/dealls-case-study.git
cd dealls-case-study
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

### 1. Create PostgreSQL Database and Seed Initial Data

```bash
go run cmd/script/main.go init
```

This command will:

- Create the database using the environment variables configured in `.env`
- Seed the following data:

  - 2 default roles: `Admin`, `Employee`
  - 1 admin user
  - 100 employee users, each with a salary, username, and password

### 2. Run the application

```bash
go run cmd/app/main.go
```

Your server will start on:
👉 `http://localhost:8080`

---

## 🧪 Running Tests

### ✅ Run All Tests (Unit + Integration)

```bash
go test ./...
```

### ✅ Run Unit Tests Only (Exclude Integration)

```bash
go test $(go list ./... | grep -v /internal/handlers)
```

### ✅ Run Integration Tests Only

```bash
go test ./internal/handlers
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

> All users created by the seed script use the default password `"password"` and have incremental usernames:
> `user1`, `user2`, ..., `user100`.
>
> These users are assigned the **Employee** role, and their salaries increase with their ID (e.g. `user1` has salary `1000`, `user2` has `2000`, etc.).
>
> An **admin** user is also created with username `"admin"` and the same password.

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

### `POST /api/v1/attendances/check-in`

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

### `POST /api/v1/attendances/check-out`

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

### `POST /api/v1/attendances/overtime`

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

### `POST /api/v1/reimbursements`

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

All `/api/v1/payrolls` routes require authentication with a **Bearer token** belonging to a user with the **Admin** role.

### `POST /api/v1/payrolls/{year}/{month}`

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

### `POST /api/v1/payrolls/{year}/{month}/run`

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

### `GET /api/v1/payrolls/{year}/{month}/summary`

Returns a summary of all employee payslips for the specified month and year.

#### Response (200 OK)

```json
{
  "message": "success",
  "data": {
    "payroll_id": 1,
    "year": 2025,
    "month": 6,
    "total_take_home": 109000000,
    "payslips": [
      {
        "user_id": 2,
        "username": "johndoe",
        "base_salary": 4000000,
        "overtime_pay": 100000,
        "reimbursement": 50000,
        "total_pay": 4150000
      },
      {
        "user_id": 3,
        "username": "janedoe",
        "base_salary": 4500000,
        "overtime_pay": 200000,
        "reimbursement": 100000,
        "total_pay": 4800000
      }
    ]
  }
}
```

---

## 🧾 Payslip

### `GET /api/v1/payslips/{year}/{month}`

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
    "base_salary": 40000,
    "overtime_pay": 4000,
    "reimbursement": 1000,
    "total_salary": 45000,
    "monthly_salary": 44000,
    "expected_working_days": 22,
    "days_attended": 20,
    "hourly_rate": 250,
    "overtime_rate_per_hour": 500,
    "total_hours_worked": 160,
    "total_overtime_hours": 8,
    "attendance_breakdown": [
      { "date": "2025-06-01" },
      { "date": "2025-06-02" },
      ...
    ],
    "overtime_breakdown": [
      { "date": "2025-06-08", "hours_worked": 2 },
      { "date": "2025-06-09", "hours_worked": 2 },
      ...
    ],
    "reimbursement_breakdown": [
      { "date": "2025-06-10", "amount": 500, "description": "Taxi to office" },
      { "date": "2025-06-12", "amount": 500, "description": "Client lunch" }
    ]
  }
}
```
