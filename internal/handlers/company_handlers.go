package handlers

import (
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

func CreateNewCompany(uc usecase.CompanyUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		var company entity.Company
		if err := c.ShouldBindJSON(&company); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		err := uc.Create(c.Request.Context(), company)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create company"})
			return
		}

		c.JSON(201, company)
	}
}
