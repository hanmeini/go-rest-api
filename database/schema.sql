-- Buat database baru (jalankan ini sekali secara manual di DBeaver)
-- CREATE DATABASE go_flix_db;

-- Skrip untuk membuat tabel movies
CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY,
    judul VARCHAR(255) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    tahun_rilis INT NOT NULL,
    sutradara VARCHAR(100) NOT NULL,
    pemeran TEXT[] NOT NULL,
    
    -- Kolom Audit Standar
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    version INT DEFAULT 1
);