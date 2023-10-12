package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	db     *sql.DB
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API Server Port")
	flag.StringVar(
		&cfg.db.dsn,
		"db-dsn",
		"postgres://postgres:password@localhost:5432/result?sslmode=disable",
		"PostgresSql DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(
		&cfg.db.maxIdleTime,
		"db-max-idle-time",
		"15m",
		"PostgreSQL max connection idle time",
	)

	flag.Parse()

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &application{
		config: cfg,
		db:     db,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err, nil)
	}
}

func connectDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) serve() error {
	port := strconv.Itoa(app.config.port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.resultHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Print("starting server")
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) resultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
