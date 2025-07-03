package handlers

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"strings"

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

		files := form.File
		photos := make([]*multipart.FileHeader, 0)
		for i, fhArr := range files {
			if strings.HasPrefix(i, "photo_") {
				photos = append(photos, fhArr[0])
			}
		}

		var courtSchedule []entity.CourtSchedule
		schedule := form.Value["schedule"][0]
		if err := json.Unmarshal([]byte(schedule), &courtSchedule); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid schedule input"})
			return
		}

		court.CourtSchedule = courtSchedule

		err = uc.Create(c, court, photos)
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

		var courtSchedule []entity.CourtSchedule
		schedule := form.Value["schedule"][0]
		if err := json.Unmarshal([]byte(schedule), &courtSchedule); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid schedule input"})
			return
		}
		court.CourtSchedule = courtSchedule

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

// TODO - Check if this handler is really needed
func FindCourtByIDShowcase(uc usecase.CourtUseCase) func(*gin.Context) {
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

func ListAvailableBookingSlots(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		courtID := c.Param("id")
		date := c.Query("date")
		slots, err := uc.ListAvailableBookingSlots(c.Request.Context(), courtID, date)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to list available slots"})
			return
		}

		c.JSON(200, slots)
	}
}

func CreateNewBooking(uc usecase.BookingUsecase) func(*gin.Context) {
	return func(c *gin.Context) {
		courtId := c.Param("id")

		var booking entity.Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		booking.CourtId = courtId

		id, err := uc.Create(c.Request.Context(), booking)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to create booking"})
			return
		}

		c.JSON(201, gin.H{"message": "Booking created successfully", "id": id})
	}
}

func ChangeCourtStatus(uc usecase.CourtUseCase) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		var court entity.Court
		if id == "" {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if err := c.ShouldBindJSON(&court); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		err := uc.UpdateCourtStatus(c, id, court)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to update court"})
			return
		}

		c.JSON(200, gin.H{"message": "Court updated successfully"})
	}
}
