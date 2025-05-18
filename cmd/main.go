package main

import (
	"log"

	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/repository"
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

	_ = repository.NewCompanyRepository(db)
	_ = repository.NewCourtRepository(db)

}
