package handlers

import (
	"log"
	"time"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

func ListBookingsByCompany(uc usecase.BookingUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		companyID := c.Query("company_id")
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		var startDate, endDate *time.Time
        layout := time.RFC3339

		if startDateStr != "" {
			t, err := time.Parse(layout, startDateStr)
			if err != nil {
                log.Println(err)
				c.JSON(400, gin.H{"error": "Invalid start_date format"})
				return
			}
			startDate = &t
		}

		if endDateStr != "" {
			t, err := time.Parse(layout, endDateStr)
			if err != nil {
                log.Println(err)
				c.JSON(400, gin.H{"error": "Invalid end_date format"})
				return
			}
			endDate = &t
		}

		filter := entity.BookingFilter{
			StartDate: startDate,
			EndDate:   endDate,
		}

        bookings, err := uc.ListByCompanyID(c.Request.Context(), companyID, filter)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to list courts"})
			return
		}

		c.JSON(200, bookings)
	}
}

func ConfirmBooking(uc usecase.BookingUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		companyID := c.Param("company_id")
		bookingID := c.Param("booking_id")

		var input struct {
			VerificationCode string `json:"verification_code"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		err := uc.ConfirmBooking(c.Request.Context(), companyID, bookingID, input.VerificationCode)
		if err != nil {
			log.Println(err)
			if err == entity.ErrInvalidVerificationCode {
				c.JSON(400, gin.H{"error": "Invalid verification code"})
				return
			}

			c.JSON(500, gin.H{"error": "Failed to confirm booking"})
			return
		}

		c.JSON(200, gin.H{"message": "Booking confirmed"})
	}
}

func FindBookingByID(uc usecase.BookingUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		bookingID := c.Param("id")
		booking, err := uc.FindByID(c.Request.Context(), bookingID)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to find booking"})
			return
		}

		c.JSON(200, booking)
	}
}

func FindBookingByIDShowcase(uc usecase.BookingUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		bookingID := c.Query("id")

		booking, err := uc.FindByIDShowcase(c.Request.Context(), bookingID)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to find booking"})
			return
		}

		c.JSON(200, booking)
	}
}
