package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env is missing")
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database init error: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	log.Println("⚡ Terkoneksi ke PostgreSQL dengan sukses!")

	// Jalankan migrasi schema
	_, err = DB.Exec(initSchemaSQL)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi schema: %v", err)
	}
	log.Println("🗄️  Migrasi database berhasil dijalankan!")
}
