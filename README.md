# Nama pembuat

Nama: Naufal Nur Fahriza
Role: Programmer Golang

# Go RESTful API - User CRUD

API RESTful yang dibangun dengan Go untuk melakukan operasi CRUD pada resource pengguna (user) dengan fitur autentikasi menggunakan Bearer Token dan MySQL sebagai database.

## Fitur

- Operasi CRUD lengkap pada resource user
- Autentikasi dengan Bearer Token menggunakan JWT
- Komunikasi data dalam format JSON
- Database MySQL untuk penyimpanan data
- Penanganan error yang baik
- Menggunakan goroutines untuk operasi database
- Database migration

## Teknologi yang Digunakan

- [Go](https://golang.org/) - Bahasa pemrograman
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router dan dispatcher
- [GORM](https://gorm.io/) - ORM library untuk Go
- [JWT-Go](https://github.com/golang-jwt/jwt) - Implementasi JWT untuk Go
- [GoDotEnv](https://github.com/joho/godotenv) - Untuk mengelola environment variables
- [MySQL](https://www.mysql.com/) - Database

## Prasyarat

- Go 1.16 atau lebih baru
- MySQL 5.7 atau lebih baru
- Git

## Cara Menjalankan Aplikasi

### Backend (Go API)

1. Clone repository ini:
   ```bash
   git clone https://github.com/username/go-restapi-crud.git
   cd go-restapi-crud
   ```

2. Install dependensi:
   ```bash
   go mod download
   ```

3. Siapkan database MySQL:
   ```sql
   CREATE DATABASE go_crud_api;
   ```

4. Setup file .env (copy dari .env):
   ```bash
   cp .env
   ```

5. Edit file .env dengan kredensial database dan konfigurasi lainnya:
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=go_crud_api
   DB_USERNAME=your_username
   DB_PASSWORD=your_password
   JWT_SECRET=your_jwt_secret
   ```

6. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

7. API sekarang berjalan di http://localhost:8080

### Frontend (Coming soon)

## Struktur API

### Auth Endpoints

| Method | Endpoint       | Deskripsi             | Autentikasi |
|--------|----------------|----------------------|-------------|
| POST   | /api/register  | Mendaftarkan user baru | Tidak       |
| POST   | /api/login     | Login user            | Tidak       |

### User Endpoints

| Method | Endpoint       | Deskripsi             | Autentikasi |
|--------|----------------|----------------------|-------------|
| POST   | /api/users     | Membuat user baru     | Ya          |
| GET    | /api/users     | Mendapatkan semua user| Ya          |
| GET    | /api/users/{id}| Mendapatkan user by ID| Ya          |
| PUT    | /api/users/{id}| Mengupdate user by ID | Ya          |
| DELETE | /api/users/{id}| Menghapus user by ID  | Ya          |

## Dokumentasi API

### Register User

**Request:**
```http
POST /api/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Respons Sukses:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2023-03-18T12:34:56Z",
    "updated_at": "2023-03-18T12:34:56Z"
  }
}
```

**Respons Error:**
```json
{
  "error": true,
  "message": "User with this email already exists"
}
```

## Postman Screenshot

Berikut adalah screenshot dari endpoint login di Postman:

![Postman Login Screenshot](./postman_screenshot/postman_login/screenshot.png)

### Login User

**Request:**
```http
POST /api/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Respons Sukses:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2023-03-18T12:34:56Z",
    "updated_at": "2023-03-18T12:34:56Z"
  }
}
```

**Respons Error:**
```json
{
  "error": true,
  "message": "Invalid email or password"
}
```

## Postman Screenshot

Berikut adalah screenshot dari endpoint login di Postman:

![Postman Login Screenshot](./postman_screenshot/postman_login/screenshot.png)

### Create User

**Request:**
```http
POST /api/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "password": "securepassword"
}
```

**Respons Sukses:**
```json
{
  "id": 2,
  "name": "Jane Smith",
  "email": "jane@example.com",
  "created_at": "2023-03-18T13:45:12Z",
  "updated_at": "2023-03-18T13:45:12Z"
}
```

**Respons Error:**
```json
{
  "error": true,
  "message": "User with this email already exists"
}
```

### Get All Users

**Request:**
```http
GET /api/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Respons Sukses:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2023-03-18T12:34:56Z",
    "updated_at": "2023-03-18T12:34:56Z"
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "email": "jane@example.com",
    "created_at": "2023-03-18T13:45:12Z",
    "updated_at": "2023-03-18T13:45:12Z"
  }
]
```

### Get User by ID

**Request:**
```http
GET /api/users/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Respons Sukses:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2023-03-18T12:34:56Z",
  "updated_at": "2023-03-18T12:34:56Z"
}
```

**Respons Error:**
```json
{
  "error": true,
  "message": "User not found"
}
```

### Update User

**Request:**
```http
PUT /api/users/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "John Doe Updated"
}
```

**Respons Sukses:**
```json
{
  "id": 1,
  "name": "John Doe Updated",
  "email": "john@example.com",
  "created_at": "2023-03-18T12:34:56Z",
  "updated_at": "2023-03-18T14:12:34Z"
}
```

**Respons Error:**
```json
{
  "error": true,
  "message": "User not found"
}
```

### Delete User

**Request:**
```http
DELETE /api/users/1
Authorization: Bearer eyJhbGciO