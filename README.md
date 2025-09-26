# Go Flix API

REST API sederhana untuk manajemen data film, dibangun dengan Go, `gorilla/mux`, dan validasi menggunakan `go-playground/validator`. API ini menggunakan in-memory storage (map) sehingga mudah dijalankan tanpa dependensi database.

## Fitur
- CRUD film (create, read all/by ID, update, delete)
- Validasi input request (required fields, minimal tahun rilis 1888)
- CORS middleware (mendukung GET, POST, PUT, DELETE, OPTIONS)
- Health check endpoint
- Unit test untuk handler pembuatan film
- Token-based auth sederhana (Login/Logout) dengan in-memory token store
- Middleware otentikasi untuk melindungi semua endpoint film

## Teknologi
- Go (module: `go-flix-api`)
- `github.com/gorilla/mux` untuk routing
- `github.com/go-playground/validator/v10` untuk validasi
- `github.com/google/uuid` untuk ID

## Struktur Proyek
```
.
├─ main.go
├─ handlers/
│  ├─ handler_movie.go
│  ├─ auth_handler.go
│  └─ handler_movie_test.go
├─ middleware/
│  └─ auth_middleware.go
├─ auth/
│  ├─ auth.go
│  └─ auth_test.go
├─ models/
│  └─ movie.go
├─ go.mod
└─ go.sum
```

## Menjalankan Aplikasi
Pastikan Go terinstal. Dari root proyek:

```bash
go mod tidy
go run main.go
```

Server akan berjalan di:
- http://localhost:8080

Log startup akan menampilkan daftar endpoint yang tersedia.

## Konfigurasi (Users)
File `config.yaml` berisi daftar user untuk login:

```yaml
users:
  - username: "user1"
    password: "password123"
  - username: "user2"
    password: "password456"
```

`main.go` akan memanggil `auth.LoadConfig("config.yaml")` saat startup dan gagal bila file tidak valid.

## CORS
CORS sudah diaktifkan via `corsMiddleware` di `main.go`:
- Mengizinkan origin `*`
- Mengizinkan metode `GET, POST, PUT, DELETE, OPTIONS`
- Menangani preflight `OPTIONS`

Jika Anda membutuhkan kredensial (cookies/Authorization), ganti `*` menjadi origin spesifik (mis. `http://127.0.0.1:5500`) dan tambahkan:
```go
w.Header().Set("Access-Control-Allow-Credentials", "true")
```

## Endpoint API

Base path: `/api`

- POST `/api/login` — login, mengembalikan token
- POST `/api/logout` — logout, mencabut token saat ini

Semua endpoint film di bawah ini DIPROTEKSI middleware `AuthMiddleware`.
Header `Authorization: Bearer <token>` wajib disertakan.

- POST `/api/movies` — buat film baru
- GET `/api/movies` — ambil semua film
- GET `/api/movies/{id}` — ambil film berdasarkan ID
- PUT `/api/movies/{id}` — perbarui film berdasarkan ID (parsial, hanya field yang dikirim)
- DELETE `/api/movies/{id}` — hapus film berdasarkan ID

Utility:
- GET `/health` — health check

### Model
- Movie
```json
{
  "id": "uuid",
  "judul": "string",
  "genre": "string",
  "tahun_rilis": 2008,
  "sutradara": "string",
  "pemeran": ["string"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Request Model
- CreateMovieRequest (POST /api/movies)
```json
{
  "judul": "string (required)",
  "genre": "string (required)",
  "tahun_rilis": 1888,
  "sutradara": "string (required)",
  "pemeran": ["string"]  // required
}
```

- UpdateMovieRequest (PUT /api/movies/{id}) — semua field opsional
```json
{
  "judul": "string",
  "genre": "string",
  "tahun_rilis": 2000,
  "sutradara": "string",
  "pemeran": ["string"]
}
```

## Contoh cURL

- Login (dapatkan token)
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}' | jq -r .token)
echo "TOKEN=$TOKEN"
```

- Akses endpoint film (harus menyertakan Bearer token)
```bash
curl http://localhost:8080/api/movies \
  -H "Authorization: Bearer $TOKEN"
```

- Logout (mencabut token)
```bash
curl -X POST http://localhost:8080/api/logout \
  -H "Authorization: Bearer $TOKEN"
```

- Create
```bash
curl -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "judul": "The Dark Knight",
    "genre": "Action",
    "tahun_rilis": 2008,
    "sutradara": "Christopher Nolan",
    "pemeran": ["Christian Bale", "Heath Ledger"]
  }'
```

- Get all
```bash
curl http://localhost:8080/api/movies \
  -H "Authorization: Bearer $TOKEN"
```

- Get by ID
```bash
curl http://localhost:8080/api/movies/{id} \
  -H "Authorization: Bearer $TOKEN"
```

- Update
```bash
curl -X PUT http://localhost:8080/api/movies/{id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "genre": "Action/Crime",
    "pemeran": ["Christian Bale", "Heath Ledger", "Aaron Eckhart"]
  }'
```

- Delete
```bash
curl -X DELETE http://localhost:8080/api/movies/{id} \
  -H "Authorization: Bearer $TOKEN"
```

## Testing
Jalankan unit test:

```bash
go test ./...
```

Lebih detail:
```bash
go test -v ./handlers
go test -v ./auth
go test -v ./handlers -run ^TestCreateMovie_Success$
go test -cover ./...
```

## Catatan
- Database in-memory akan kosong ulang setiap restart.
- Token disimpan in-memory; server restart akan menghapus token aktif.
- Jika Anda mengubah origin frontend (mis. `localhost` vs `127.0.0.1`), pastikan setelan CORS sesuai.
- Peringatan `favicon.ico` 404 di frontend aman diabaikan; tambahkan `<link rel="icon" href="data:,">` jika perlu.
