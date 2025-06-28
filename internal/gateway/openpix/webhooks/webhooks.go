package webhooks

import (
	"log"

	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TODO - Check its not updating booking status to confirmed if the payment is confirmed
func ConfirmedPaymentWebhook(puc usecase.PaymentUsecase, buc usecase.BookingUsecase) func(*gin.Context) {
    return func(c *gin.Context) {
        var in openpix.ChargeWebhookEvent
        err := c.ShouldBindJSON(&in)
        if err != nil {
            log.Println("Error binding JSON:", err)
            c.JSON(400, gin.H{"status": "error", "message": "Invalid request data"})
            return
        }

        err = puc.ConfirmPayment(c.Request.Context(), in.Charge)
        if err != nil {
            log.Println("Error confirming payment:", err)
            c.JSON(500, gin.H{"status": "error", "message": "Failed to confirm payment"})
            return
        }
    
        c.JSON(200, gin.H{"status": "success", "message": "Payment confirmed"})
    }
}

func ExpiredPaymentWebhook(uc usecase.PaymentUsecase) func(*gin.Context) {
    return func(c *gin.Context) {
        var in openpix.ChargeWebhookEvent
        err := c.ShouldBindJSON(&in)
        if err != nil {
            log.Println("Error binding JSON:", err)
            c.JSON(400, gin.H{"status": "error", "message": "Invalid request data"})
            return
        }

        err = uc.ExpirePayment(c.Request.Context(), in.Charge)
        if err != nil {
            log.Println("Error expiring payment:", err)
            c.JSON(500, gin.H{"status": "error", "message": "Failed to expire payment"})
            return
        }
    
        c.JSON(200, gin.H{"status": "success", "message": "Payment expired"})
    }
}

func CancelBookingWebhook(uc usecase.BookingUsecase) func(*gin.Context) {
    return func(c *gin.Context) {
        var in openpix.ChargeWebhookEvent
        err := c.ShouldBindJSON(&in)
        if err != nil {
            log.Println("Error binding JSON:", err)
            c.JSON(400, gin.H{"status": "error", "message": "Invalid request data"})
            return
        }

        //err = uc.CancelBooking(c.Request.Context(), in.Charge)
        if err != nil {
            log.Println("Error canceling booking:", err)
            c.JSON(500, gin.H{"status": "error", "message": "Failed to cancel booking"})
            return
        }
    
        c.JSON(200, gin.H{"status": "success", "message": "Booking canceled"})
    }
}
