package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	//logger
	//db driver
}


//Mount
func (app *application) mount() http.Handler{
	r := chi.NewRouter()
	r.Use(middleware.RequestID)	//for Rate limiting
	r.Use(middleware.RealIP) //also for rate limiting + analytics and tracing
	// r.Use(middleware.Logger) //
	r.Use(middleware.Recoverer) //For crashouts
	r.Use(cors)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})


	// Current
	r.Get("/api/logout", handleLogOut)
	r.Get("/api/validate/", handleAPIValidation)
	r.Get("/oauth/lastfm/login", handleLoginFlow)
	r.Get("/oauth/lastfm/callback", handleCallbackFlow)

	return r
}


//Run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: h,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Minute * 1,
	}

	return srv.ListenAndServe()
}

type config struct {
	addr string //Address
	db dbConfig
}

type dbConfig struct {
	dsn string //Domain String: user=YY password=XX 
}
