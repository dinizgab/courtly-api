package handlers

import (
	"log"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

func CreateNewCompany(uc usecase.CompanyUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		var company entity.Company
		if err := c.ShouldBindJSON(&company); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		token, err := uc.Create(c.Request.Context(), company)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to create company"})
			return
		}

		c.JSON(201, gin.H{
            "message": "Company created successfully",
            "token":   token,
            "company": company,
        })
	}
}

func LoginCompany(uc usecase.CompanyUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		token, err := uc.Login(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			log.Println(err)
			if err == entity.ErrInvalidCredentials {
				c.JSON(401, gin.H{"error": "Invalid credentials"})
			}
			return
		}

		c.JSON(200, gin.H{
			"message": "Login successful",
			"token":   token,
		})
	}
}
