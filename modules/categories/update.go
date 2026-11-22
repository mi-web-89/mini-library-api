package categories

import (
	"context"
	"mini-library-api/db/models"

	"gorm.io/gorm"
)

type update struct {
	repository *gorm.DB
	get
}

func (u update) Update(ctx context.Context, id uint, categoryReq *CategoryRequest) (Category, error) {
	category := models.Category{
		Name:        categoryReq.Name,
		Description: categoryReq.Description,
	}

	result := u.repository.WithContext(ctx).Where("id = ?", id).Updates(&category)

	if err := result.Error; err != nil {
		return Category{}, err
	}

	record, err := u.get.GetByID(ctx, id)

	if err != nil {
		return Category{}, err
	}

	return record, nil
}
