package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopherSocial/internal/db"
	"gopherSocial/internal/env"
	"gopherSocial/internal/store"
	"log"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8083"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://gopheruser:123@localhost:5436/gophersocial?sslmode=disable"),
			maxOpenConns: 30, // ideally you should get this value from environment variables
			maxIdleConns: 30,
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	// Open a new connection to the database
	database, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Panic(err)
		}
	}(database)
	logger.Info("DATABASE CONNECTION POOL ESTABLISHED")
	s := store.NewStorage(database)

	app := &application{
		config: cfg,
		store:  s,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
