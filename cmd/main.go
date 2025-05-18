package main

import (
	"log"

	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/handlers"
	"github.com/dinizgab/booking-mvp/internal/repository"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	config := config.New()

	db, err := database.New(config.DB)
	if err != nil {
		log.Println(err)
	}

	companyRepository := repository.NewCompanyRepository(db)
	courtRepository := repository.NewCourtRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

    courtUsecase := usecase.NewCourtUseCase(courtRepository)
	_ = usecase.NewCompanyUsecase(companyRepository)
	_ = usecase.NewBookingUsecase(bookingRepository)

    router := gin.Default()

    router.POST("/courts", handlers.CreateCourt(courtUsecase))
    router.GET("/courts/:id", handlers.FindCourtByID(courtUsecase))
    router.GET("/companies/:company_id/courts", handlers.ListCourtsByCompany(courtUsecase))
    router.GET("/courts/:id/bookings", handlers.ListCourtBookingsByID(courtUsecase))
    router.PUT("/courts/:id", handlers.UpdateCourt(courtUsecase))
    router.DELETE("/courts/:id", handlers.DeleteCourt(courtUsecase))
}
