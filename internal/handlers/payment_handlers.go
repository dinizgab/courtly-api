package handlers

import (
	"log"

	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

func GetBookingPaymentStatus(uc usecase.PaymentUsecase) func(*gin.Context) {
    return func(c *gin.Context) {
        bookingID := c.Query("id")
        status, err := uc.GetBookingPaymentStatusByID(c.Request.Context(), bookingID)
        if err != nil {
            log.Println(err)
            c.JSON(500, gin.H{"error": "Failed to get payment status"})
            return
        }

        c.JSON(200, gin.H{"status": status})
    }
}

func GetBookingChargeInformation(uc usecase.PaymentUsecase) func(*gin.Context) {
    return func(c *gin.Context) {
        bookingID := c.Param("id")
        charge, err := uc.GetBookingChargeInformation(c.Request.Context(), bookingID)
        if err != nil {
            log.Println(err)
            c.JSON(500, gin.H{"error": "Failed to get booking charge information"})
            return
        }
        c.JSON(200, gin.H{
            "charge": charge,
        })
    }
}
