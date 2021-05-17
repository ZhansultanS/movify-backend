package main

import (
	"flag"
	"fmt"
	"github.com/DARKestMODE/movify/internal/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "host=localhost user=postgres password=postgres dbname=movify port=5432 sslmode=disable TimeZone=Asia/Shanghai", "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	db.AutoMigrate(&data.Movie{}, &data.Genre{})

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	baseUrl, _ := url.Parse(fmt.Sprintf("http://localhost%s/v1/healthcheck", srv.Addr))

	logger.Printf("starting %s server on %s %s", cfg.env, srv.Addr, baseUrl)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func openDB(cfg config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
