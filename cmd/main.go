package main

import (
	"fmt"
	"log"

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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}
	config := config.New()

	db, err := database.New(config.DB)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	companyRepository := repository.NewCompanyRepository(db)
	courtRepository := repository.NewCourtRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

	courtUsecase := usecase.NewCourtUseCase(courtRepository)
	companyUsecase := usecase.NewCompanyUsecase(companyRepository)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepository)

	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/courts", handlers.CreateCourt(courtUsecase))
	router.GET("/courts/:id", handlers.FindCourtByID(courtUsecase))
	router.GET("/companies/:company_id/courts", handlers.ListCourtsByCompany(courtUsecase))
	router.GET("/courts/:id/bookings", handlers.ListCourtBookingsByID(courtUsecase))
	router.PUT("/courts/:id", handlers.UpdateCourt(courtUsecase))
	router.DELETE("/courts/:id", handlers.DeleteCourt(courtUsecase))

	router.GET("/companies/:company_id/bookings", handlers.ListBookingsByCompany(bookingUsecase))
	router.GET("/bookings/:id", handlers.FindBookingByID(bookingUsecase))
	router.PATCH("/companies/:company_id/bookings/:booking_id/confirm", handlers.ConfirmBooking(bookingUsecase))

	router.POST("/auth/signup", handlers.CreateNewCompany(companyUsecase))
	router.POST("auth/login", handlers.LoginCompany(companyUsecase))

	router.Run(fmt.Sprintf(":%s", config.API.Port))
}
