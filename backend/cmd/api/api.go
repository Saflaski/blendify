package main

import (
	"backend-lastfm/internal/auth"
	blend "backend-lastfm/internal/blending"
	musicapi "backend-lastfm/internal/music_api/lastfm"
	network "backend-lastfm/internal/network"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/glog"
	"github.com/redis/go-redis/v9"
)

type application struct {
	config config
	//logger
	dbConfig dbConfig
	external externalConfig
}

// Mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID) //for Rate limiting
	r.Use(middleware.RealIP)    //also for rate limiting + analytics and tracing
	// r.Use(middleware.Logger) //
	r.Use(middleware.Recoverer) //For crashouts
	r.Use(network.Cors)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	//Connect to Redis Client

	rdb := redis.NewClient(&redis.Options{
		Addr:     app.config.db.addrString,
		Password: app.config.db.password,
		DB:       app.config.db.db,
		Protocol: app.config.db.protocol,
	})

	apiKey := app.config.external.apiKey
	lastFMURL := app.config.external.lastFMURL

	LastFMExternal := musicapi.NewLastFMExternalAdapter(
		apiKey,
		lastFMURL,
		true,
	)

	authCfg := auth.Config{
		ExpiryDuration:     time.Duration(app.config.sessionExpiry),
		FrontendCookieName: "sid",
		FrontendURL:        "http://localhost:5173",
	}

	authRepo := auth.NewRedisStateStore(rdb) // Placeholder nil, replace with actual Redis client
	authService := auth.NewAuthService(authRepo, authCfg)
	authHandler := auth.NewAuthHandler(
		authService,
		authCfg,
	)

	blendRepo := blend.NewRedisStateStore(rdb)
	blendService := blend.NewBlendService(*blendRepo, *LastFMExternal)
	blendHandler := blend.NewBlendHandler(
		"http://localhost:5173",
		"sid",
		*blendService)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/blend", func(r chi.Router) {
			r.Use(auth.ValidateCookie(*authHandler, authService))
			r.Get("/new", blendHandler.GetBlendPercentage)
			r.Post("/add", blendHandler.AddBlendFromInviteLink)
			r.Get("/generatelink", blendHandler.GenerateNewLink)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Get("/login/{platform}", authHandler.HandleLastFMLoginFlow)
			r.Post("/logout", authHandler.HandleLastFMLogOut)
			r.Get("/validate", authHandler.HandleAPIValidation)
			r.Get("/callback/{platform}", authHandler.HandleLastFMCallbackFlow)

		})
	})

	glog.Info("Mounted Handlers:")
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		glog.Infof("Route: %s %s\n", method, route)
		return nil
	})
	return r
}

// Run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		ReadTimeout:  time.Second * 10, //Blanket Read, Write, Idle Timeouts as safety net.
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute * 1,
	}

	glog.Info("Server Started")
	glog.Infof("Address: %s", srv.Addr)
	glog.Infof("ReadTimeout: %f", srv.ReadTimeout.Seconds())
	glog.Infof("WriteTimeout: %f", srv.WriteTimeout.Seconds())
	glog.Infof("IdleTimeout: %f", srv.IdleTimeout.Seconds())

	return srv.ListenAndServe()
}

type config struct {
	addr          string //Address
	db            dbConfig
	external      externalConfig
	sessionExpiry int
}

type dbConfig struct {
	addrString string
	password   string
	db         int
	protocol   int
}

type externalConfig struct {
	apiKey    string
	lastFMURL string
}
