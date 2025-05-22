package handlers

import (
	"encoding/json"
	"log"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/usecase"
	"github.com/gin-gonic/gin"
)

func CreateCourt(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid form data"})
			return
		}

		courtInfos := form.Value["court_info"][0]
		var court entity.Court
		if err := json.Unmarshal([]byte(courtInfos), &court); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// TODO - Save court photos
		//files := form.File
		//photos := make([]*multipart.FileHeader, 0)
		//for i, fhArr := range files {
		//	if strings.HasPrefix(i, "photo_") {
		//		photos = append(photos, fhArr[0])
		//	}
		//}

		err = uc.Create(c.Request.Context(), court)
		if err != nil {
			log.Println(err)
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
			log.Println(err)
			c.JSON(404, gin.H{"error": "Court not found"})
			return
		}

		c.JSON(200, court)
	}
}

func ListCourtsByCompany(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		companyID := c.Param("id")
		courts, err := uc.ListByCompany(c.Request.Context(), companyID)
		if err != nil {
			log.Println(err)
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
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to list bookings"})
			return
		}

		c.JSON(200, bookings)
	}
}

func ListCompanyCourtShowcase(uc usecase.CourtUseCase) func(*gin.Context) {
    return func(c *gin.Context) {
        companyID := c.Param("id")
        courts, err := uc.ListCompanyCourtsShowcase(c.Request.Context(), companyID)
        if err != nil {
            log.Println(err)
            c.JSON(500, gin.H{"error": "Failed to list courts"})
            return
        }

        c.JSON(200, courts)
    }
}

func UpdateCourt(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		form, err := c.MultipartForm()
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid form data"})
			return
		}

		courtInfos := form.Value["court_info"][0]
		var court entity.Court
		if err := json.Unmarshal([]byte(courtInfos), &court); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// TODO - Save court photos
		//files := form.File
		//photos := make([]*multipart.FileHeader, 0)
		//for i, fhArr := range files {
		//	if strings.HasPrefix(i, "photo_") {
		//		photos = append(photos, fhArr[0])
		//	}
		//}

		err = uc.Update(c.Request.Context(), id, court)
		if err != nil {
			log.Println(err)
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
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to delete court"})
			return
		}

		c.JSON(200, gin.H{"message": "Court deleted successfully"})
	}
}
