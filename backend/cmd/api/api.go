package main

import (
	"backend-lastfm/internal/auth"
	blend "backend-lastfm/internal/blending"
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"backend-lastfm/internal/musicbrainz"
	network "backend-lastfm/internal/network"
	shared "backend-lastfm/internal/shared"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
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
		200,
	)

	authCfg := auth.Config{
		ExpiryDuration: time.Duration(app.config.sessionExpiry) * time.Second,
		// ExpiryDuration:     time.Duration(app.config.sessionExpiry) * time.Second,
		FrontendCookieName: "sid",
		FrontendURL:        os.Getenv("FRONTEND_URL"),
		BackendURL:         os.Getenv("BACKEND_URL"),
	}

	authRepo := auth.NewRedisStateStore(rdb, authCfg.ExpiryDuration) // Placeholder nil, replace with actual Redis client
	authService := auth.NewAuthService(authRepo, authCfg)
	authHandler := auth.NewAuthHandler(
		*authService,
		authCfg,
	)

	MBsqlxDB := sqlx.MustConnect("pgx", os.Getenv("MUSICBRAINZ_DB_DSN"))
	MBsqlxDB.SetMaxOpenConns(25)
	MBsqlxDB.SetMaxIdleConns(25)
	MBsqlxDB.SetConnMaxLifetime(5 * time.Minute)

	BlendifysqlxDB := sqlx.MustConnect("pgx", os.Getenv("BLENDIFY_DB_DSN"))
	BlendifysqlxDB.SetMaxOpenConns(25)
	BlendifysqlxDB.SetMaxIdleConns(25)
	BlendifysqlxDB.SetConnMaxLifetime(5 * time.Minute)

	mbRepo := musicbrainz.NewPostgresMusicBrainzRepo(MBsqlxDB)
	mbService := musicbrainz.NewMBService(mbRepo)

	blendRepo := blend.NewBlendStore(rdb, BlendifysqlxDB)
	blendService := blend.NewBlendService(*blendRepo, *LastFMExternal, *mbService)
	blendHandler := blend.NewBlendHandler(
		os.Getenv("FRONTEND_URL"),
		"sid",
		*blendService,
		string(auth.UserKey),
	)

	sharedService := shared.NewSharedService(authService, blendService)
	sharedHandler := shared.NewSharedHandler(*sharedService)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/blend", func(r chi.Router) {
			r.Use(auth.ValidateCookie(*authHandler, *authService))
			// r.Get("/new", blendHandler.GetBlendPercentage)
			r.Get("/health", blendHandler.GetBlendHealth)
			// r.Post("/add/{permaLink}", blendHandler.AddBlendFromInviteLink)
			r.Post("/add", blendHandler.AddBlendFromInviteLink)
			r.Post("/delete", blendHandler.DeleteBlend)
			r.Get("/carddata", blendHandler.GetBlendPageData)
			r.Get("/cataloguedata", blendHandler.GetBlendedEntryData)
			r.Get("/userblends", blendHandler.GetUserBlends)
			r.Get("/usertopitems", blendHandler.GetUserTopItems)
			r.Get("/generatelink", blendHandler.GenerateNewLink)
			r.Get("/getpermalink", blendHandler.GetPermanentLink)
			r.Get("/userinfo", blendHandler.GetUserInfo)
			r.Get("/usertopgenres", blendHandler.GetUserTopGenres)
			r.Get("/blendtopgenres", blendHandler.GetBlendTopGenres)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Get("/login/{platform}", authHandler.HandleLastFMLoginFlow)
			r.Get("/callback/{platform}", authHandler.HandleLastFMCallbackFlow)
			r.Get("/validate", authHandler.HandleAPIValidation)

			r.Group(func(r chi.Router) {
				r.Use(auth.ValidateCookie(*authHandler, *authService))
				r.Post("/delete", sharedHandler.DeleteAllData)
				r.Post("/logout", authHandler.HandleLastFMLogOut)

			})
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
		WriteTimeout: time.Second * 120,
		IdleTimeout:  time.Minute * 1,
	}

	glog.Info("Server Started")
	glog.Infof("Address: %s", srv.Addr)
	glog.Infof("ReadTimeout: %f", srv.ReadTimeout.Seconds())
	glog.Infof("WriteTimeout: %f", srv.WriteTimeout.Seconds())
	glog.Infof("IdleTimeout: %f", srv.IdleTimeout.Seconds())
	glog.Infof("Session Valid Time: %s", (time.Duration(app.config.sessionExpiry) * time.Second).String())

	prod, _ := strconv.ParseBool(os.Getenv("PROD"))
	if prod {
		for _, o := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				glog.Info("Allowed Origin: " + o)
			}
		}
	}
	return srv.ListenAndServe()
}

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
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
