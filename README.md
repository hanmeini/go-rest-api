# 🎬 Go Flix API

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-blue.svg)](https://www.postgresql.org/)
[![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-orange.svg)](https://swagger.io/)

A modern REST API for movie management built with Go, featuring JWT authentication, PostgreSQL database, and comprehensive Swagger documentation.

## 🚀 Features

- **🔐 JWT Authentication** - Secure login/logout with token-based authentication
- **🎭 Movie Management** - Full CRUD operations for movie data
- **🗄️ PostgreSQL Database** - Robust data persistence with PostgreSQL
- **📚 Swagger Documentation** - Interactive API documentation
- **🛡️ Middleware Security** - Authentication middleware for protected routes
- **🔄 Soft Delete** - Safe deletion with audit trails
- **🌐 CORS Support** - Cross-origin resource sharing enabled
- **📊 Structured Logging** - JSON-based logging with slog
- **⚡ High Performance** - Built with Go for optimal performance

## 📋 Table of Contents

- [Installation](#-installation)
- [Configuration](#-configuration)
- [Database Setup](#-database-setup)
- [API Documentation](#-api-documentation)
- [Usage](#-usage)
- [Project Structure](#-project-structure)
- [API Endpoints](#-api-endpoints)
- [Authentication](#-authentication)
- [Contributing](#-contributing)
- [License](#-license)

## 🛠️ Installation

### Prerequisites

- Go 1.25.1 or higher
- PostgreSQL 13 or higher
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/go-flix-api.git
cd go-flix-api
```

### Install Dependencies

```bash
go mod download
```

### Build the Application

```bash
go build cmd/server/main.go
```

## ⚙️ Configuration

Create a `config.yml` file in the root directory:

```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "your_password"
  dbname: "go_flix_db"

jwt:
  secret: "your_jwt_secret_key"

users:
  - username: "admin"
    password: "admin123"
  - username: "user1"
    password: "password123"
```

### Environment Variables

You can also use environment variables:

```bash
export PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=go_flix_db
export JWT_SECRET=your_jwt_secret_key
```

## 🗄️ Database Setup

### 1. Create Database

```sql
CREATE DATABASE go_flix_db;
```

### 2. Run Schema

Execute the SQL schema from `database/schema.sql`:

```sql
CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY,
    judul VARCHAR(255) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    tahun_rilis INT NOT NULL,
    sutradara VARCHAR(100) NOT NULL,
    pemeran TEXT[] NOT NULL,
    
    -- Audit Columns
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    version INT DEFAULT 1
);
```

### 3. Verify Connection

```bash
psql -h localhost -p 5432 -U postgres -d go_flix_db -c "\dt"
```

## 📚 API Documentation

### Swagger UI

Once the server is running, access the interactive API documentation at:

**🌐 [http://localhost:8080/swagger/](http://localhost:8080/swagger/)**

### API Specification Files

- **JSON**: `docs/swagger.json`
- **YAML**: `docs/swagger.yaml`

## 🚀 Usage

### Start the Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Health Check

```bash
curl http://localhost:8080/health
```

## 📁 Project Structure

```
go-flix-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── config/
│   └── config.go               # Configuration management
├── database/
│   └── schema.sql              # Database schema
├── docs/                       # Generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── auth/                   # Authentication module
│   │   ├── handler.go          # Auth HTTP handlers
│   │   └── middleware.go       # Auth middleware
│   ├── middleware/
│   │   └── auth_middleware.go  # JWT middleware
│   └── movie/                  # Movie management module
│       ├── handler.go          # Movie HTTP handlers
│       ├── repository.go       # Database operations
│       └── service.go          # Business logic
├── models/
│   └── movie.go                # Data models
├── config.yml                  # Configuration file
├── go.mod                      # Go module file
├── go.sum                      # Go module checksums
└── README.md                   # This file
```

## 🔗 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/login` | User login | ❌ |
| POST | `/api/logout` | User logout | ✅ |

### Movies

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/movies` | Get all movies | ✅ |
| GET | `/api/movies/{id}` | Get movie by ID | ✅ |
| POST | `/api/movies` | Create new movie | ✅ |
| PUT | `/api/movies/{id}` | Update movie | ✅ |
| DELETE | `/api/movies/{id}` | Delete movie | ✅ |

### System

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | ❌ |
| GET | `/swagger/` | Swagger UI | ❌ |

## 🔐 Authentication

### Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Using JWT Token

Include the token in the Authorization header:

```bash
curl -X GET http://localhost:8080/api/movies \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Logout

```bash
curl -X POST http://localhost:8080/api/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📝 Example Usage

### Create a Movie

```bash
curl -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "judul": "The Avengers",
    "genre": "Action",
    "tahun_rilis": 2012,
    "sutradara": "Joss Whedon",
    "pemeran": ["Robert Downey Jr.", "Chris Evans", "Scarlett Johansson"]
  }'
```

### Get All Movies

```bash
curl -X GET http://localhost:8080/api/movies \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update a Movie

```bash
curl -X PUT http://localhost:8080/api/movies/{movie-id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "judul": "The Avengers: Endgame",
    "tahun_rilis": 2019
  }'
```

## 🧪 Testing

### Using Swagger UI

1. Open [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
2. Click "Authorize" and enter your JWT token
3. Test endpoints interactively

### Using cURL

```bash
# Login and get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}' | \
  jq -r '.token')

# Use token for authenticated requests
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/movies
```

## 🔧 Development

### Regenerate Swagger Documentation

```bash
swag init -g cmd/server/main.go
```

### Run Tests

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

### Lint Code

```bash
golangci-lint run
```

## 🧰 Troubleshooting & Tips

### pq: relation "movies" does not exist

- Pastikan aplikasi terhubung ke database yang benar sesuai `config.yml` (dbname `go_flix_db`).
- Buat tabel jika belum ada, lalu restart server:

```sql
CREATE TABLE IF NOT EXISTS public.movies (
  id UUID PRIMARY KEY,
  judul VARCHAR(255) NOT NULL,
  genre VARCHAR(100) NOT NULL,
  tahun_rilis INT NOT NULL,
  sutradara VARCHAR(100) NOT NULL,
  pemeran TEXT[] NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  created_by VARCHAR(100),
  updated_by VARCHAR(100),
  version INT DEFAULT 1
);
```

Verifikasi:

```sql
SELECT table_schema, table_name FROM information_schema.tables
WHERE table_schema='public' AND table_name='movies';
```

### Swagger 500 / doc.json error

- Regenerate docs: `swag init -g cmd/server/main.go`
- Pastikan import docs di `cmd/server/main.go`:
  - `import _ "go-flix-api/docs"`
- Akses UI: `http://localhost:8080/swagger/`

### DBeaver tidak menampilkan data

- Pastikan terkoneksi ke DB `go_flix_db`, schema `public` (centang schema, lalu Refresh/F5).
- Jalankan query langsung: `SELECT COUNT(*) FROM public.movies;`

### Windows PowerShell notes

- Pisahkan perintah (hindari `&&`), jalankan satu per satu.

## 🍿 Seeding Data (5 Film Cepat)

### Via API

1) Login untuk mendapatkan token

```bash
curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}'
```

2) Tambahkan beberapa film (ulang sesuai kebutuhan)

```bash
TOKEN=... # isi dari langkah login

curl -X POST http://localhost:8080/api/movies \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"judul":"Inception","genre":"Sci-Fi","tahun_rilis":2010,"sutradara":"Christopher Nolan","pemeran":["Leonardo DiCaprio","Joseph Gordon-Levitt"]}'

curl -X POST http://localhost:8080/api/movies \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"judul":"The Dark Knight","genre":"Action","tahun_rilis":2008,"sutradara":"Christopher Nolan","pemeran":["Christian Bale","Heath Ledger"]}'

curl -X POST http://localhost:8080/api/movies \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"judul":"Interstellar","genre":"Sci-Fi","tahun_rilis":2014,"sutradara":"Christopher Nolan","pemeran":["Matthew McConaughey","Anne Hathaway"]}'

curl -X POST http://localhost:8080/api/movies \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"judul":"Avengers: Endgame","genre":"Action","tahun_rilis":2019,"sutradara":"Anthony Russo","pemeran":["Robert Downey Jr.","Chris Evans"]}'

curl -X POST http://localhost:8080/api/movies \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"judul":"Parasite","genre":"Thriller","tahun_rilis":2019,"sutradara":"Bong Joon-ho","pemeran":["Song Kang-ho","Choi Woo-shik"]}'
```

### Via SQL (langsung di DB)

```sql
INSERT INTO public.movies (id,judul,genre,tahun_rilis,sutradara,pemeran,created_at,updated_at,version)
VALUES
(gen_random_uuid(),'Inception','Sci-Fi',2010,'Christopher Nolan',ARRAY['Leonardo DiCaprio','Joseph Gordon-Levitt'],NOW(),NOW(),1),
(gen_random_uuid(),'The Dark Knight','Action',2008,'Christopher Nolan',ARRAY['Christian Bale','Heath Ledger'],NOW(),NOW(),1),
(gen_random_uuid(),'Interstellar','Sci-Fi',2014,'Christopher Nolan',ARRAY['Matthew McConaughey','Anne Hathaway'],NOW(),NOW(),1),
(gen_random_uuid(),'Avengers: Endgame','Action',2019,'Anthony Russo',ARRAY['Robert Downey Jr.','Chris Evans'],NOW(),NOW(),1),
(gen_random_uuid(),'Parasite','Thriller',2019,'Bong Joon-ho',ARRAY['Song Kang-ho','Choi Woo-shik'],NOW(),NOW(),1);
```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Error**
   - Verify PostgreSQL is running
   - Check database credentials in `config.yml`
   - Ensure database exists

2. **JWT Token Invalid**
   - Check token expiration (default: 1 hour)
   - Verify JWT secret in configuration
   - Ensure proper Authorization header format

3. **Port Already in Use**
   - Change port in `config.yml`
   - Kill existing process: `taskkill /F /IM main.exe`

### Logs

The application uses structured JSON logging. Check console output for detailed error messages.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push to branch: `git push origin feature/new-feature`
5. Submit a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new features
- Update documentation for API changes
- Use meaningful commit messages

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Your Name** - *Initial work* - [YourGitHub](https://github.com/yourusername)

## 🙏 Acknowledgments

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [SQLX](https://github.com/jmoiron/sqlx) - SQL extensions
- [Swagger](https://swagger.io/) - API documentation

---

⭐ **Star this repository if you found it helpful!**