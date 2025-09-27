# JWT Authentication Setup

## Setup yang Diperlukan

### 1. Buat file `.env` di root project
Buat file `.env` dengan isi:
```
JWT_SECRET_KEY=your-super-secret-jwt-key-change-this-in-production
```

**PENTING:** Ganti `your-super-secret-jwt-key-change-this-in-production` dengan secret key yang kuat untuk production!

### 2. File yang Sudah Diperbarui untuk JWT

✅ **auth/auth.go** - Implementasi lengkap JWT authentication
✅ **handlers/auth_handler.go** - Menggunakan GenerateJWT() dan RevokeJWT()
✅ **middleware/auth_middleware.go** - Menggunakan ValidateJWT() untuk validasi
✅ **main.go** - Memuat JWT secret dari .env file
✅ **auth/auth_test.go** - Test cases untuk JWT functionality

### 3. Fitur JWT yang Tersedia

- **GenerateJWT(username)** - Membuat token JWT dengan masa berlaku 1 jam
- **ValidateJWT(token)** - Memvalidasi token dan memeriksa denylist
- **RevokeJWT(token)** - Menambahkan token ke denylist untuk logout
- **Denylist** - Menyimpan token yang sudah di-revoke di memory

### 4. Cara Penggunaan

1. **Login:** POST `/api/login` dengan username/password
2. **Gunakan token:** Tambahkan header `Authorization: Bearer <token>`
3. **Logout:** POST `/api/logout` dengan token di header Authorization

### 5. Endpoint yang Dilindungi

Semua endpoint `/api/movies/*` memerlukan JWT token yang valid.
