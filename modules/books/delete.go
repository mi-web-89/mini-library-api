package books

import (
	"context"

	"gorm.io/gorm"
)

type deletebook struct {
	repository *gorm.DB
	get
}

func (d deletebook) Delete(ctx context.Context, id uint) error {
	book, err := d.get.GetByID(ctx, id)

	if err != nil {
		return err
	}

	// soft delete
	if err := d.repository.Delete(&book).Error; err != nil {
		return err
	}

	return nil

}
