package books

import (
	"context"
	"mini-library-api/db/models"

	"gorm.io/gorm"
)

type update struct {
	repository *gorm.DB
	get
}

func (u update) Update(ctx context.Context, id uint, bookReq *BookRequest) (Book, error) {
	book := models.Book{
		Title:       bookReq.Title,
		Description: bookReq.Description,
		Publisher:   bookReq.Publisher,
		Year:        bookReq.Year,
		Stock:       bookReq.Stock,
		ImageLink:   bookReq.ImageLink,
		ImageID:     bookReq.ImageID,
		CategoryID:  bookReq.CategoryID,
	}

	result := u.repository.WithContext(ctx).Where("id = ?", id).Updates(&book)

	if err := result.Error; err != nil {
		return Book{}, err
	}

	record, err := u.get.GetByID(ctx, id)

	if err != nil {
		return Book{}, err
	}

	return record, nil
}
