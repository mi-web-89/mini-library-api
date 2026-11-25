package books

import (
	"mime/multipart"
	"mini-library-api/modules/categories"
	"time"

	"gorm.io/gorm"
)

// struct for DTO / response
type Book struct {
	ID          uint                `json:"id" gorm:"primaryKey" `
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Publisher   string              `json:"publisher"`
	Year        uint                `json:"year"`
	Stock       uint                `json:"stock"`
	ImageLink   string              `json:"image_link"`
	ImageID     string              `json:"image_id"`
	CategoryID  uint                `json:"category_id"`
	Category    categories.Category `json:"category"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `json:"deleted_at" gorm:"index"`
}

// karena nantinya akan mengupload sebuah file, maka request ini akan menjadi sebuah form-data
// imageID dan ImageLink terisi ketika file berhasil di upload
type BookRequest struct {
	Title       string `form:"title" validate:"required"`
	Description string `form:"description" validate:"required"`
	Publisher   string `form:"publisher" validate:"required"`
	Year        uint   `form:"year" validate:"required"`
	Stock       uint   `form:"stock" validate:"required"`
	CategoryID  uint   `form:"category_id" validate:"required"`
	ImageLink   string
	ImageID     string
	File        *multipart.FileHeader
}
