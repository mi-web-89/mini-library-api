package categories

import (
	"context"
	"mini-library-api/pkg/utils"

	"gorm.io/gorm"
)

type Service interface {
	GetAll(ctx context.Context, pagination utils.Pagination) (utils.Pagination, error)
	GetByID(ctx context.Context, id uint) (Category, error)
	Create(ctx context.Context, categoryReq *CategoryRequest) (Category, error)
	Update(ctx context.Context, id uint, categoryReq *CategoryRequest) (Category, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	getall
	get
	create
	update
	deletecategory
}

var _ Service = (*service)(nil)

func New(repository *gorm.DB) Service {
	return service{
		getall:         getall{repository: repository},
		get:            get{repository: repository},
		create:         create{repository: repository},
		update:         update{repository: repository, get: get{repository: repository}},
		deletecategory: deletecategory{repository: repository, get: get{repository: repository}},
	}
}
