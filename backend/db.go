package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type MyDB struct {
	DB *sql.DB
}

var (
	instance *MyDB
	once     sync.Once
)

func InitDB() {

	once.Do(func() {

		db, err := sql.Open("sqlite3", "compras.db")

		if err != nil {
			log.Fatalf("❌ Error al abrir SQLite: %v", err)
		}

		if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
			log.Fatalf("❌ Error activando foreign keys: %v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("❌ No se pudo conectar a SQLite: %v", err)
		}

		log.Println("✅ Conectado a SQLite")

		instance = &MyDB{DB: db}
	})
}

func GetDB() *sql.DB {
	if instance == nil {
		InitDB()
	}
	return instance.DB
}
