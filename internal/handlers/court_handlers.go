package handlers

import (
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

func CreateCourt(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		var court entity.Court
		if err := c.ShouldBindJSON(&court); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		err := uc.Create(c.Request.Context(), court)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create court"})
			return
		}

		c.JSON(201, gin.H{"message": "Court created successfully"})
	}
}

func FindCourtByID(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		court, err := uc.FindByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(404, gin.H{"error": "Court not found"})
			return
		}

		c.JSON(200, court)
	}
}

func ListCourtsByCompany(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		companyID := c.Param("company_id")
		courts, err := uc.ListByCompany(c.Request.Context(), companyID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to list courts"})
			return
		}

		c.JSON(200, courts)
	}
}

func ListCourtBookingsByID(uc usecase.CourtUseCase) func(*gin.Context) {
    return func(c *gin.Context) {
        id := c.Param("id")
        bookings, err := uc.ListBookingsByID(c.Request.Context(), id)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to list bookings"})
            return
        }

        c.JSON(200, bookings)
    }
}

func UpdateCourt(uc usecase.CourtUseCase) func(*gin.Context) {
    return func(c *gin.Context) {
        id := c.Param("id")
        var court entity.Court
        if err := c.ShouldBindJSON(&court); err != nil {
            c.JSON(400, gin.H{"error": "Invalid input"})
            return
        }

        err := uc.Update(c.Request.Context(), id, court)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to update court"})
            return
        }

        c.JSON(200, gin.H{"message": "Court updated successfully"})
    }
}

func DeleteCourt(uc usecase.CourtUseCase) func(*gin.Context) {
    return func(c *gin.Context) {
        id := c.Param("id")
        err := uc.Delete(c.Request.Context(), id)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to delete court"})
            return
        }

        c.JSON(200, gin.H{"message": "Court deleted successfully"})
    }
}
