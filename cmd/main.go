package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dinizgab/booking-mvp/internal/auth"
	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/handlers"
	"github.com/dinizgab/booking-mvp/internal/repository"
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

    fmt.Println(os.Getenv("DATABASE_URL"))

	config := config.New()
	db, err := database.New(config.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	authService := auth.NewAuthService([]byte(config.API.JwtSecret))

	companyRepository := repository.NewCompanyRepository(db)
	courtRepository := repository.NewCourtRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

	courtUsecase := usecase.NewCourtUseCase(courtRepository)
	companyUsecase := usecase.NewCompanyUsecase(companyRepository, authService)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepository)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000", "https://courtly-red.vercel.app"},
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
	protected.Use(auth.AuthMiddleware(authService))
	{
		protected.GET("/companies/:id/dashboard", handlers.GetCompanyDashboard(companyUsecase))

		protected.POST("/courts", handlers.CreateCourt(courtUsecase))
		protected.GET("/courts/:id", handlers.FindCourtByID(courtUsecase))
		protected.GET("/courts/:id/bookings", handlers.ListCourtBookingsByID(courtUsecase))
		protected.PUT("/courts/:id", handlers.UpdateCourt(courtUsecase))
		protected.DELETE("/courts/:id", handlers.DeleteCourt(courtUsecase))

		protected.GET("/companies/:id/courts", handlers.ListCourtsByCompany(courtUsecase))
		protected.GET("/companies/:id", handlers.FindCompanyByID(companyUsecase))
		protected.PUT("/companies/:id", handlers.UpdateCompanyInformations(companyUsecase))
		protected.GET("/companies/:id/bookings", handlers.ListBookingsByCompany(bookingUsecase))
		protected.GET("/bookings/:id", handlers.FindBookingByID(bookingUsecase))
		protected.PATCH("/companies/:company_id/bookings/:booking_id/confirm", handlers.ConfirmBooking(bookingUsecase))
	}

	public := router.Group("/showcase")
	{
		public.GET("/companies/:id/courts", handlers.ListCompanyCourtShowcase(courtUsecase))
		public.GET("/courts/:id", handlers.FindCourtByIDShowcase(courtUsecase))
		public.GET("/courts/:id/available-slots", handlers.ListAvailableBookingSlots(courtUsecase))
		public.GET("/bookings", handlers.FindBookingByIDShowcase(bookingUsecase))
		public.POST("/courts/:id/bookings", handlers.CreateNewBooking(bookingUsecase))
	}

	router.Run(fmt.Sprintf(":%s", config.API.Port))
}
