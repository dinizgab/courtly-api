package webhooks

import (
	"fmt"

	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TODO - Check its not updating booking status to confirmed if the payment is confirmed
func ConfirmedPaymentWebhook(uc usecase.PaymentUsecase) func(*gin.Context) {
    return func(c *gin.Context) {
        var in openpix.ChargeWebhookEvent
        err := c.ShouldBindJSON(&in)
        if err != nil {
            fmt.Println("Error binding JSON:", err)
            c.JSON(400, gin.H{"status": "error", "message": "Invalid request data"})
            return
        }

        err = uc.ConfirmPayment(c.Request.Context(), in.Charge)
        if err != nil {
            fmt.Println("Error confirming payment:", err)
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
            fmt.Println("Error binding JSON:", err)
            c.JSON(400, gin.H{"status": "error", "message": "Invalid request data"})
            return
        }

        err = uc.ExpirePayment(c.Request.Context(), in.Charge)
        if err != nil {
            fmt.Println("Error expiring payment:", err)
            c.JSON(500, gin.H{"status": "error", "message": "Failed to expire payment"})
            return
        }
    
        c.JSON(200, gin.H{"status": "success", "message": "Payment expired"})
    }
}
