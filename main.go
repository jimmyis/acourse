package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/acoshift/acourse/app"
	"github.com/acoshift/acourse/postgres"
	"github.com/acoshift/configfile"
	"github.com/acoshift/gzip"
	"github.com/acoshift/hsts"
	"github.com/acoshift/middleware"
	"github.com/acoshift/redirecthttps"
	_ "github.com/lib/pq"
)

func main() {
	time.Local = time.UTC

	config := configfile.NewReader("config")

	sqlURL := config.String("sql_url")

	// init databases
	db, err := sql.Open("postgres", sqlURL)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(5)

	userRepo, err := postgres.NewUserRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Init(app.Config{
		ProjectID:      config.String("project_id"),
		ServiceAccount: config.Bytes("service_account"),
		BucketName:     config.String("bucket"),
		EmailServer:    config.String("email_server"),
		EmailPort:      config.Int("email_port"),
		EmailUser:      config.String("email_user"),
		EmailPassword:  config.String("email_password"),
		EmailFrom:      config.String("email_from"),
		BaseURL:        config.String("base_url"),
		XSRFSecret:     config.String("xsrf_key"),
		DB:             db,
		RedisAddr:      config.String("redis_addr"),
		RedisPass:      config.String("redis_pass"),
		RedisPrefix:    config.String("redis_prefix"),
		SessionSecret:  config.Bytes("session_secret"),
		SlackURL:       config.String("slack_url"),
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	h := middleware.Chain(
		redirecthttps.New(redirecthttps.Config{Mode: redirecthttps.OnlyProxy}),
		hsts.New(hsts.PreloadConfig),
		gzip.New(gzip.DefaultConfig),
	)(app.Handler(userRepo))
	mux.Handle("/", h)

	// lets reverse proxy handle other settings
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Start server at :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
