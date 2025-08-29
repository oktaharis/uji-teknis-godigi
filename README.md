# Uji Teknis Backend Engineer – GoDigi

<p align="center">
  <img src="https://img.icons8.com/ios-filled/100/database.png" alt="GoDigi Logo" width="80"/>
</p>

<h1 align="center">🚀 GoDigi Backend API</h1>

<p align="center">
  RESTful API untuk sistem manajemen leads dan projects dengan autentikasi JWT.<br/>
  Dibangun dengan Go, Gin Framework, dan MySQL.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go&logoColor=white"/>
  <img src="https://img.shields.io/badge/Gin_Framework-1.9-0096D6?logo=go&logoColor=white"/>
  <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white"/>
  <img src="https://img.shields.io/badge/JWT-Auth-000000?logo=json-web-tokens&logoColor=white"/>
  <img src="https://img.shields.io/badge/GORM-1.25-000000?logo=go&logoColor=white"/>
</p>

---

## 🚀 Cara Menjalankan

### 1. Clone Repository
```bash
git clone https://github.com/oktaharis/uji-teknis-godigi.git
cd uji-teknis-godigi
```

### 2. Import Database
- Buka **phpMyAdmin**
- Buat database baru `godigi`
- Import file **dataset.sql** yang ada di root project

### 3. Konfigurasi Environment
```bash
cp .env.example .env
```
Sesuaikan koneksi DB dan JWT secret di file `.env`

### 4. Jalankan Aplikasi
```bash
go mod tidy
go run cmd/api/main.go
```
Server berjalan di `http://localhost:8080`

---

## 🛠️ Tools
- **Go 1.20+** - Bahasa pemrograman
- **Gin Framework** - Web framework
- **GORM** - ORM untuk MySQL
- **JWT** - Autentikasi token
- **VSCode REST Client** - Testing API (pakai file `godigi.rest`)
- **cURL + jq** - Opsional, untuk test via terminal

---

## 📡 API Endpoints

### 🔐 Authentication
- `POST /auth/register` - Register user baru
- `POST /auth/login` - Login user
- `POST /auth/logout` - Logout (memerlukan token)
- `POST /auth/forgot-password` - Lupa password
- `POST /auth/reset-password` - Reset password
- `GET  /me` - Get profile user (memerlukan token)

### 📋 Leads Management
- `POST /leads` - Create new lead
- `GET  /leads` - Get all leads
- `GET  /leads/:id` - Get lead by ID
- `PUT  /leads/:id` - Update lead
- `DELETE /leads/:id` - Delete lead
- `GET  /leads/summary` - Get leads summary

### 📂 Projects Management
- `POST /projects` - Create new project
- `GET  /projects` - Get all projects
- `GET  /projects/:id` - Get project by ID
- `PUT  /projects/:id` - Update project
- `DELETE /projects/:id` - Delete project

### 👨‍💼 Admin Management (role=admin)
- `POST /admin/users` - Create new user
- `GET  /admin/users` - Get all users
- `GET  /admin/users/:id` - Get user by ID
- `PUT  /admin/users/:id` - Update user
- `DELETE /admin/users/:id` - Delete user

---

## ⚡ Quick Test with cURL

### Set Variabel Dasar
```bash
BASE_URL=http://localhost:8080
NAME="Okta Haris"
EMAIL="okta@example.com"
PASS="secret123"
NEWPASS="newSecret123"
```

### Healthcheck
```bash
curl -s $BASE_URL/healthz | jq
```

### Register
```bash
curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" | jq
```

### Login
```bash
TOKEN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" | jq -r '.data.token')
echo $TOKEN
```

### Profile (Me)
```bash
curl -s "$BASE_URL/me" -H "Authorization: Bearer $TOKEN" | jq
```

### Logout
```bash
curl -s -X POST "$BASE_URL/auth/logout" -H "Authorization: Bearer $TOKEN" | jq
```

### Forgot Password
```bash
RESET_TOKEN=$(curl -s -X POST "$BASE_URL/auth/forgot-password" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\"}" | jq -r '.data.reset_token')
echo $RESET_TOKEN
```

### Reset Password
```bash
curl -s -X POST "$BASE_URL/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$RESET_TOKEN\",\"new_password\":\"$NEWPASS\"}" | jq
```

### Login dengan Password Baru
```bash
TOKEN2=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$NEWPASS\"}" | jq -r '.data.token')
echo $TOKEN2
```

### Leads CRUD + Summary

#### Create Lead
```bash
LEAD_ID=$(curl -s -X POST "$BASE_URL/leads" \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d '{"company_name":"PT Godigi","contact_name":"Budi","email":"budi@godigi.co.id","phone":"081234567890","source":"website","industry":"IT","region":"Jabodetabek","sales_rep":"Okta","status":"New","notes":"Inquired via landing page"}' | jq -r '.data.id')
echo $LEAD_ID
```

#### Get Lead by ID
```bash
curl -s "$BASE_URL/leads/$LEAD_ID" -H "Authorization: Bearer $TOKEN2" | jq
```

#### List Leads
```bash
curl -s "$BASE_URL/leads?q=Godigi&status=New&page=1&per_page=10" -H "Authorization: Bearer $TOKEN2" | jq
```

#### Update Lead
```bash
curl -s -X PUT "$BASE_URL/leads/$LEAD_ID" \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d '{"status":"Qualified","notes":"Booked demo"}' | jq
```

#### Summary Leads
```bash
curl -s "$BASE_URL/leads/summary" -H "Authorization: Bearer $TOKEN2" | jq
```

#### Delete Lead
```bash
curl -s -X DELETE "$BASE_URL/leads/$LEAD_ID" -H "Authorization: Bearer $TOKEN2" | jq
```

### Projects CRUD

#### Create Project
```bash
PROJECT_ID=$(curl -s -X POST "$BASE_URL/projects" \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d '{"name":"Godigi CRM","description":"Internal CRM untuk tim sales","status":"planned","start_date":"2025-09-01"}' | jq -r '.data.id')
echo $PROJECT_ID
```

#### List Projects
```bash
curl -s "$BASE_URL/projects?q=CRM&status=planned&page=1&per_page=10" -H "Authorization: Bearer $TOKEN2" | jq
```

#### Get Project by ID
```bash
curl -s "$BASE_URL/projects/$PROJECT_ID" -H "Authorization: Bearer $TOKEN2" | jq
```

#### Update Project
```bash
curl -s -X PUT "$BASE_URL/projects/$PROJECT_ID" \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress","end_date":"2025-12-31"}' | jq
```

#### Delete Project
```bash
curl -s -X DELETE "$BASE_URL/projects/$PROJECT_ID" -H "Authorization: Bearer $TOKEN2" | jq
```

### Admin – Users Management

**Note:** Pastikan role user sudah diubah menjadi admin di database:
```sql
UPDATE users SET role='admin' WHERE email='okta@example.com';
```

#### Create User
```bash
curl -s -X POST "$BASE_URL/admin/users" \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d '{"name":"Admin User","email":"admin@example.com","password":"admin123","role":"admin"}' | jq
```

#### List Users
```bash
curl -s "$BASE_URL/admin/users?q=okta&page=1&per_page=10" -H "Authorization: Bearer $TOKEN2" | jq
```

#### Get User by ID
```bash
curl -s "$BASE_URL/admin/users/1" -H "Authorization: Bearer $TOKEN2" | jq
```

#### Update User
```bash
curl -s -X PUT "$BASE_URL/admin/users/1" \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}' | jq
```

#### Delete User
```bash
curl -s -X DELETE "$BASE_URL/admin/users/1" -H "Authorization: Bearer $TOKEN2" | jq
```

---

## 📁 Project Structure

```
cmd/
└── api/
    └── main.go          # Application entry point
internal/
├── config/              # Configuration setup
├── controllers/         # HTTP controllers
├── middleware/          # Custom middleware
├── models/             # Database models
├── repositories/       # Data access layer
├── services/           # Business logic
└── utils/              # Utility functions
pkg/
└── database/           # Database connection
.env                    # Environment variables
godigi.rest             # API test requests
dataset.sql             # Database schema and sample data
```

---

## 📋 Requirements

- MySQL 8.0+
- Go 1.20+
- Git

---

## 🐛 Troubleshooting

Jika遇到 masalah koneksi database, pastikan:
1. MySQL service berjalan
2. Konfigurasi database di `.env` benar
3. Database `godigi` sudah dibuat
4. Tabel leads dan projects sudah diimport dari `dataset.sql`

---

## 📄 License

Proyek ini dibuat untuk tujuan uji teknis Backend Engineer di GoDigi.

---

**Happy Coding! 🚀**