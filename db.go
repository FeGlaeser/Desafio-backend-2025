package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	// Carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis do ambiente.")
	}

	// Tenta usar a DATABASE_URL diretamente
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Se DATABASE_URL não existir, monta manualmente com variáveis individuais
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")

		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	}

	// Abre conexão
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao abrir conexão com o banco:", err)
	}

	// Testa conexão
	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar com o banco:", err)
	}

	log.Println("Banco de dados conectado com sucesso!")
}
