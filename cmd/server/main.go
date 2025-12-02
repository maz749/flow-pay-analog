package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/maz749/flow-pay-analog/internal/handlers"
	"github.com/maz749/flow-pay-analog/internal/middleware"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/internal/scheduler"
	"github.com/maz749/flow-pay-analog/internal/telegram"
	"github.com/maz749/flow-pay-analog/pkg/config"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.New(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWT.Secret)
	subHandler := handlers.NewSubscriptionHandler(subRepo, notifRepo)
	telegramHandler := handlers.NewTelegramHandler(userRepo)

	// Initialize router
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware)

	// Public routes
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Protected routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

	// Auth routes
	api.HandleFunc("/auth/profile", authHandler.GetProfile).Methods("GET")

	// Subscription routes
	api.HandleFunc("/subscriptions", subHandler.Create).Methods("POST")
	api.HandleFunc("/subscriptions", subHandler.GetAll).Methods("GET")
	api.HandleFunc("/subscriptions/stats", subHandler.GetStats).Methods("GET")
	api.HandleFunc("/subscriptions/{id}", subHandler.GetByID).Methods("GET")
	api.HandleFunc("/subscriptions/{id}", subHandler.Update).Methods("PUT")
	api.HandleFunc("/subscriptions/{id}", subHandler.Delete).Methods("DELETE")

	// Telegram routes
	api.HandleFunc("/telegram/settings", telegramHandler.GetSettings).Methods("GET")
	api.HandleFunc("/telegram/settings", telegramHandler.UpdateSettings).Methods("PUT")

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Serve landing pages
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	}).Methods("GET")

	router.HandleFunc("/personal", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/personal.html")
	}).Methods("GET")

	router.HandleFunc("/business", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/business.html")
	}).Methods("GET")

	router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/about.html")
	}).Methods("GET")

	router.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/blog.html")
	}).Methods("GET")

	router.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/privacy.html")
	}).Methods("GET")

	// Signup redirects to home page with registration modal
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/#register", http.StatusSeeOther)
	}).Methods("GET")

	// Catch-all route for SPA (must be last)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	// Initialize Telegram bot (optional, only if token is provided)
	var bot *telegram.Bot
	var sched *scheduler.Scheduler

	if cfg.Telegram.BotToken != "" {
		bot, err = telegram.NewBot(cfg.Telegram.BotToken, db)
		if err != nil {
			log.Printf("Warning: Failed to initialize Telegram bot: %v", err)
			log.Println("Continuing without Telegram integration...")
		} else {
			// Start bot in goroutine
			go bot.Start()

			// Initialize and start scheduler
			sched = scheduler.NewScheduler(subRepo, notifRepo, userRepo, bot)
			sched.Start()
			defer sched.Stop()
		}
	} else {
		log.Println("Telegram bot token not provided, running without Telegram integration")
	}

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		if sched != nil {
			sched.Stop()
		}
		os.Exit(0)
	}()

	log.Printf("Server starting on %s", addr)
	log.Printf("Environment: %s", cfg.App.Environment)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
