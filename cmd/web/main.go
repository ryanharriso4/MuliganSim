package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql" // New import
	"muligansim.ryanharris.net/internal/models"
)

type application struct {
	logger         *slog.Logger
	cards          *models.CardModel
	users          *models.UserModel
	sessionManager *scs.SessionManager
	templateCache  map[string]*template.Template
}

func main() {

	port := flag.String("port", ":4000", "Port that the website broadcast off of")
	dsn := flag.String("dsn", "put database info here", "MySQL data source name")
	key := flag.String("key", "put tls key here", "tls key")
	cert := flag.String("cert", "put tls cert here", "tls cert")

	flag.Parse()

	print(key, cert)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger:         logger,
		cards:          &models.CardModel{DB: db},
		users:          &models.UserModel{DB: db},
		sessionManager: sessionManager,
		templateCache:  templateCache,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Broadcasting on port ", "addr", srv.Addr)

	err = srv.ListenAndServeTLS(*cert, *key)
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
