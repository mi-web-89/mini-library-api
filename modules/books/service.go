package books

import (
	"context"
	"mini-library-api/pkg/utils"

	"gorm.io/gorm"
)

type Service interface {
	GetAll(ctx context.Context, pagination utils.Pagination) (utils.Pagination, error)
	GetByID(ctx context.Context, id uint) (Book, error)
	GetByCategory(ctx context.Context, pagination utils.Pagination, cotegoryId uint) (utils.Pagination, error)
	Create(ctx context.Context, bookReq *BookRequest) (Book, error)
	Update(ctx context.Context, id uint, bookReq *BookRequest) (Book, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	getall
	get
	getbycategory
	create
	update
	deletebook
}

var _ Service = (*service)(nil)

func New(repository *gorm.DB) Service {
	return service{
		getall:        getall{repository: repository},
		get:           get{repository: repository},
		getbycategory: getbycategory{repository: repository},
		create:        create{repository: repository},
		update:        update{repository: repository, get: get{repository: repository}},
		deletebook:    deletebook{repository: repository, get: get{repository: repository}},
	}
}
