package usecase

import (
	"context"
	"fmt"
	"github.com/dinizgab/booking-mvp/internal/services/storage"
	"github.com/google/uuid"
	"mime/multipart"
	"strings"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type (
	courtUseCaseImpl struct {
		courtRepository repository.CourtRepository
		uploadStorage   storage.StorageUploader
	}

	CourtUseCase interface {
		Create(ctx context.Context, court entity.Court, photos []*multipart.FileHeader) error
		FindByID(ctx context.Context, id string) (entity.Court, error)
		ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error)
		ListCompanyCourtsShowcase(ctx context.Context, companyID string) ([]entity.Court, error)
		ListBookingsByID(ctx context.Context, id string) ([]entity.Booking, error)
		ListAvailableBookingSlots(ctx context.Context, id string, date string) ([]entity.Booking, error)
		Update(ctx context.Context, id string, court entity.Court) error
		Delete(ctx context.Context, id string) error
	}
)

func NewCourtUseCase(courtRepository repository.CourtRepository, uploadStorage storage.StorageUploader) CourtUseCase {
	return &courtUseCaseImpl{
		courtRepository: courtRepository,
		uploadStorage:   uploadStorage,
	}
}

func (u *courtUseCaseImpl) Create(ctx context.Context, court entity.Court, photos []*multipart.FileHeader) error {
	court.ID = uuid.New().String()
	err := u.courtRepository.Create(ctx, &court)
	if err != nil {
		return err
	}

	var photosEntities []entity.CourtPhoto
	for position, fh := range photos {
		photoId := uuid.New().String()
		fileType := strings.Split(fh.Filename, ".")[1]
		filename := fmt.Sprintf("%s.%s", photoId, fileType)

		file, err := fh.Open()
		if err != nil {
			return fmt.Errorf("CourtUsecase.Create - could not open photo file: %w", err)
		}

		publicURL, err := u.uploadStorage.UploadFile(ctx, court.ID, filename, file)
		if err != nil {
			return err
		}

		photo := entity.CourtPhoto{
			ID:       photoId,
			CourtId:  court.ID,
			Path:     publicURL,
			Position: position,
			IsCover:  position == 0,
		}

		photosEntities = append(photosEntities, photo)

		err = file.Close()
		if err != nil {
			return fmt.Errorf("CourtUsecase.Create - could not close photo file: %w", err)
		}
	}

	err = u.courtRepository.InsertPhotos(ctx, photosEntities)
	if err != nil {
		// TODO - Delete photos from storage if error occurs in the database
		return err
	}

	return nil
}

func (u *courtUseCaseImpl) FindByID(ctx context.Context, id string) (entity.Court, error) {
	court, err := u.courtRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Court{}, err
	}

	return court, nil
}

func (u *courtUseCaseImpl) ListBookingsByID(ctx context.Context, id string) ([]entity.Booking, error) {
	bookings, err := u.courtRepository.ListBookingsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (u *courtUseCaseImpl) ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error) {
	courts, err := u.courtRepository.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return courts, nil
}

func (u *courtUseCaseImpl) ListCompanyCourtsShowcase(ctx context.Context, companyID string) ([]entity.Court, error) {
	courts, err := u.courtRepository.ListCompanyCourtsShowcase(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return courts, nil
}

func (u *courtUseCaseImpl) ListAvailableBookingSlots(ctx context.Context, id string, date string) ([]entity.Booking, error) {
	bookings, err := u.courtRepository.ListAvailableBookingSlots(ctx, id, date)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (u *courtUseCaseImpl) Update(ctx context.Context, id string, court entity.Court) error {
	err := u.courtRepository.Update(ctx, id, court)
	if err != nil {
		return err
	}

	return nil
}

func (u *courtUseCaseImpl) Delete(ctx context.Context, id string) error {
	err := u.courtRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
