package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dinizgab/booking-mvp/internal/services/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinizgab/booking-mvp/internal/auth"
	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
	"github.com/dinizgab/booking-mvp/internal/gateway/openpix/webhooks"
	"github.com/dinizgab/booking-mvp/internal/handlers"
	"github.com/dinizgab/booking-mvp/internal/repository"
	"github.com/dinizgab/booking-mvp/internal/services/notification"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading env file: %v", err)
		}
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error loading cfg: %v", err)
	}

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	authService := auth.NewAuthService(cfg.API.JwtSecret)
	emailRenderer, err := notification.NewHTMLRender(nil)
	if err != nil {
		log.Fatalf("Failed to create email renderer: %v", err)
	}

	pixGatewayClient := openpix.NewOpenPixClient(cfg.OpenPix)
	emailService := notification.NewEmailSender(emailRenderer, cfg.SMTP)
	storageUploadService := storage.NewSupabaseStorageUploader(cfg.Storage, "court-photos")

	companyRepository := repository.NewCompanyRepository(db)
	courtRepository := repository.NewCourtRepository(db)
	bookingRepository := repository.NewBookingRepository(db)
	paymentRepository := repository.NewPaymentRepository(db)

	pixPaymentUsecase := usecase.NewPixGatewayService(pixGatewayClient, bookingRepository, paymentRepository, emailService)
	courtUsecase := usecase.NewCourtUseCase(courtRepository, storageUploadService)
	companyUsecase := usecase.NewCompanyUsecase(companyRepository, authService, pixPaymentUsecase)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepository, pixPaymentUsecase, companyUsecase, courtUsecase)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://courtly-red.vercel.app", "https://courtly.com.br/"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Authorization", "Content-Type", "Accept", "Access-Control-Request-Headers",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/auth/signup", handlers.CreateNewCompany(companyUsecase))
	router.POST("/auth/login", handlers.LoginCompany(companyUsecase))

	protected := router.Group("/admin")
	protected.Use(auth.Middleware(authService))
	{
		protected.GET("/companies/:id/dashboard", handlers.GetCompanyDashboard(companyUsecase))
		protected.GET("/companies/:id/balance", handlers.GetCompanyBalance(pixPaymentUsecase))

		protected.POST("/companies/:id/withdraw", handlers.CreateWithdrawRequest(pixPaymentUsecase))

		protected.POST("/courts", handlers.CreateCourt(courtUsecase))
		protected.GET("/courts/:id", handlers.FindCourtByID(courtUsecase))
		protected.GET("/courts/:id/bookings", handlers.ListCourtBookingsByID(courtUsecase))
		protected.PUT("/courts/:id", handlers.UpdateCourt(courtUsecase))
		protected.DELETE("/courts/:id", handlers.DeleteCourt(courtUsecase))

		protected.GET("/companies/:id/courts", handlers.ListCourtsByCompany(courtUsecase))
		protected.GET("/companies/:id", handlers.FindCompanyByID(companyUsecase))
		protected.PUT("/companies/:id", handlers.UpdateCompanyInformations(companyUsecase))
		protected.GET("/bookings", handlers.ListBookingsByCompany(bookingUsecase))
		protected.GET("/bookings/:id", handlers.FindBookingByID(bookingUsecase))
		// TODO - (refactor) change this route name
		protected.PATCH("/companies/:company_id/bookings/:booking_id/confirm", handlers.ConfirmBooking(bookingUsecase))
	}

	public := router.Group("/showcase")
	{
        public.GET("/companies/:id", handlers.FindCompanyByIDShowcase(companyUsecase))
		public.GET("/companies/:id/courts", handlers.ListCompanyCourtShowcase(courtUsecase))
		public.GET("/courts/:id", handlers.FindCourtByIDShowcase(courtUsecase))
		public.GET("/courts/:id/available-slots", handlers.ListAvailableBookingSlots(courtUsecase))
		public.GET("/bookings", handlers.FindBookingByIDShowcase(bookingUsecase))
		public.POST("/courts/:id/bookings", handlers.CreateNewBooking(bookingUsecase))
		public.GET("/bookings/status", handlers.GetBookingPaymentStatus(pixPaymentUsecase))
		public.GET("/bookings/:id/charge", handlers.GetBookingChargeInformation(pixPaymentUsecase))
	}

	webhookRouter := router.Group("/webhooks")
	{
		webhookRouter.POST("/pix/confirmed", webhooks.ConfirmedPaymentWebhook(pixPaymentUsecase))
		webhookRouter.POST("/pix/expired", webhooks.ExpiredPaymentWebhook(pixPaymentUsecase))
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.API.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("cmd.main - Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
