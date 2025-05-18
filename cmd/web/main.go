package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox.samedarslan28.net/internal/models"
	"time"
)

// Define an application struct to hold the application-wide dependencies for the
// web application.
// For now, we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	snippets          models.SnippetModelInterface
	users             models.UserModelInterface
	templateCache     map[string]*template.Template
	formDecoder       *form.Decoder
	sessionManager    *scs.SessionManager
	isDebugModeActive bool
}

func main() {
	addr := flag.String("addr", ":4000", "http service address")
	debug := flag.Bool("debug", false, "enable debug mode")
	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost:5432/snippetbox?sslmode=disable", "database connection string")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			errorLog.Fatal(err)
		}
	}(db)

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		users:             &models.UserModel{DB: db},
		snippets:          &models.SnippetModel{DB: db},
		templateCache:     templateCache,
		formDecoder:       formDecoder,
		sessionManager:    sessionManager,
		isDebugModeActive: *debug,
	}

	infoLog.Printf("Starting server on %s", *addr)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		errorLog.Fatal(err)
		return
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
