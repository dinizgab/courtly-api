package openpix

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func ConfirmPaymentWebhook() func(*gin.Context) {
    return func(c *gin.Context) {
        in, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.JSON(400, gin.H{"status": "error", "message": "Failed to read request body"})
            return
        }
    
        fmt.Println("Received webhook data:", string(in))

        c.JSON(200, gin.H{"status": "success", "message": "Payment confirmed"})
    }
}
