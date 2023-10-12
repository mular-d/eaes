package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

type Result struct {
	ID        int
	Name      string
	Image     string
	Chemistry int
	Biology   int
	Maths     int
	Civic     int
	English   int
	Aptitude  int
	Physics   int
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
	mux.HandleFunc("/result", app.resultHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Print("starting server")
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) resultHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	resultID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	row := app.db.QueryRow("SELECT * FROM results WHERE id = $1", resultID)

	var result Result
	err = row.Scan(&result.ID, &result.Name, &result.Image, &result.Chemistry, &result.Biology, &result.Maths, &result.Civic, &result.English, &result.Aptitude, &result.Physics)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Result not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResult)
}
