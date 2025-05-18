package main

import (
	"log"

	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/repository"
	"github.com/dinizgab/booking-mvp/internal/usecase"
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

	_ = usecase.NewCourtUseCase(courtRepository)
	_ = usecase.NewCompanyUsecase(companyRepository)
	_ = usecase.NewBookingUsecase(bookingRepository)
}
