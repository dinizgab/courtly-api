package handlers

import (
	"log"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

func ListBookingsByCompany(uc usecase.BookingUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		companyID := c.Param("company_id")
		courts, err := uc.ListByCompanyID(c.Request.Context(), companyID)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to list courts"})
			return
		}

		c.JSON(200, courts)

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

