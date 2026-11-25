package handlers

import (
	"log"
	"mini-library-api/modules/books"
	"mini-library-api/pkg/constant"
	"mini-library-api/pkg/dtos"
	"mini-library-api/pkg/fileupload"
	"mini-library-api/pkg/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v3"
)

type Books struct {
	books    books.Service
	uploader fileupload.Uploader
}

func (b Books) GetAll(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	sort := ctx.Query("sort")
	search := ctx.Query("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "title",
	}

	booksData, err := b.books.GetAll(ctx.Context(), pagination)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "fetch books failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "all books",
		Data:    booksData,
	})
}

func (b Books) GetById(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: INVALID_ID_MESSAGE,
		})
	}

	book, err := b.books.GetByID(ctx.Context(), uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[books.Book]{
		Status:  "success",
		Message: "book found",
		Data:    book,
	})
}

func (b Books) GetByCategory(ctx fiber.Ctx) error {
	param := ctx.Params("categoryId")

	categoryID, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: INVALID_ID_MESSAGE,
		})
	}

	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	sort := ctx.Query("sort")
	search := ctx.Query("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "title",
	}

	booksData, err := b.books.GetByCategory(ctx.Context(), pagination, uint(categoryID))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "fetch books failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "books by category",
		Data:    booksData,
	})
}

func (b Books) Create(ctx fiber.Ctx) error {
	bookReq, ok := ctx.Locals("validatedBody").(*books.BookRequest)

	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "validasi failed",
		})
	}

	file, err := ctx.FormFile("file")

	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "file not found",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	isFileValid := utils.ValidateFile(ext)

	if !isFileValid {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid file format",
		})
	}

	resp, err := b.uploader.UploadFile(ctx.Context(), file, uploader.UploadParams{})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "upload failed",
		})
	}

	bookReq.ImageLink = resp.SecureURL
	bookReq.ImageID = resp.PublicID

	book, err := b.books.Create(ctx.Context(), bookReq)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(dtos.Response[books.Book]{
		Status:  "success",
		Message: "book created",
		Data:    book,
	})
}

func (b Books) Update(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: INVALID_ID_MESSAGE,
		})
	}

	bookRecord, err := b.books.GetByID(ctx.Context(), uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	bookReq := ctx.Locals("validatedBody").(*books.BookRequest)

	file, err := ctx.FormFile("file")

	if file != nil {
		if err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
				Status:  "failed",
				Message: "file not found",
			})
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		if !utils.ValidateFile(ext) {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
				Status:  "failed",
				Message: "invalid file format",
			})
		}

		resp, err := b.uploader.UploadFile(ctx.Context(), file, uploader.UploadParams{})

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
				Status:  "failed",
				Message: "upload failed",
			})
		}

		bookReq.ImageLink = resp.SecureURL
		bookReq.ImageID = resp.PublicID

		if bookRecord.ImageID != "" {
			if _, err := b.uploader.DeleteFile(ctx.Context(), bookRecord.ImageID, constant.RESOURCE_IMAGE, uploader.UploadParams{}); err != nil {
				return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
					Status:  "failed",
					Message: "failed to delete old file",
				})
			}
		}
	} else {
		bookReq.ImageLink = bookRecord.ImageLink
		bookReq.ImageID = bookRecord.ImageID
	}

	book, err := b.books.Update(ctx.Context(), uint(id), bookReq)

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[books.Book]{
		Status:  "success",
		Message: "book updated",
		Data:    book,
	})
}

func (b Books) Delete(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: INVALID_ID_MESSAGE,
		})
	}

	book, err := b.books.GetByID(ctx, uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	if err := b.books.Delete(ctx.Context(), uint(id)); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	if book.ImageID != "" {
		if _, err = b.uploader.DeleteFile(ctx.Context(), book.ImageID, constant.RESOURCE_IMAGE, uploader.UploadParams{}); err != nil {
			log.Printf("failed to delete book image %s: %v", book.ImageID, err)
		}
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[any]{
		Status:  "success",
		Message: "book deleted",
	})
}

func NewBooks(books books.Service, uploader fileupload.Uploader) Books {
	booksHandler := Books{
		books:    books,
		uploader: uploader,
	}

	return booksHandler
}
