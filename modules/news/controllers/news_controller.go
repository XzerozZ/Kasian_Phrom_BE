package controllers

import (
	"mime/multipart"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/news/usecases"

	"github.com/gofiber/fiber/v2"
)

type NewsController struct {
	newsusecase usecases.NewsUseCase
}

func NewNewsController(newsusecase usecases.NewsUseCase) *NewsController {
	return &NewsController{newsusecase: newsusecase}
}

func (c *NewsController) CreateNewsHandler(ctx *fiber.Ctx) error {
	req := &entities.News{
		Title: ctx.FormValue("title"),
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "invalid form data",
			"result":      nil,
		})
	}

	types := form.Value["type"]
	descs := form.Value["desc"]
	if len(types) != len(descs) {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "type and desc counts do not match",
			"result":      nil,
		})
	}

	files := form.File["images"]
	if len(files) == 0 {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":      "Error",
            "status_code": fiber.StatusBadRequest,
            "message":     "At least one image is required",
            "result":      nil,
        })
    }

	var fileHeaders []multipart.FileHeader
    for _, file := range files {
        fileHeaders = append(fileHeaders, *file)
    }

	for i := 0; i < len(types); i++ {
		req.Dialog = append(req.Dialog, entities.Dialog{
			Type: types[i],
			Desc: descs[i],
		})
	}

	if req.Title == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "title cannot be empty",
			"result":      nil,
		})
	}

	if len(req.Dialog) == 0 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "dialogs cannot be empty",
			"result":      nil,
		})
	}

	data, err := c.newsusecase.CreateNews(req, fileHeaders, ctx)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      	fiber.ErrInternalServerError.Message,
			"status_code": 	fiber.ErrInternalServerError.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":		"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"Nursing house created successfully",
		"result":      	data,
	})
}

func (c *NewsController) GetAllNewsHandler(ctx *fiber.Ctx) error {
	data, err := c.newsusecase.GetAllNews()
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      	fiber.ErrInternalServerError.Message,
			"status_code": 	fiber.ErrInternalServerError.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      	"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"Nursing houses retrieved successfully",
		"result":      	data,
	})
}

func (c *NewsController) GetNewsByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.newsusecase.GetNewsByID(id)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      	"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"Nursing house retrieved successfully",
		"result":      	data,
	})
}

func (c *NewsController) GetNewsNextIDHandler(ctx *fiber.Ctx) error {
	data, err := c.newsusecase.GetNewsNextID()
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      	"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"Nursing house retrieved successfully",
		"result":      	data,
	})
}

func (c *NewsController) UpdateNewsByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var news entities.News
	if err := ctx.BodyParser(&news); err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}
	

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "invalid form data",
			"result":      nil,
		})
	}
	
	types := form.Value["type"]
	descs := form.Value["desc"]
	if len(types) != len(descs) {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "type and desc counts do not match",
			"result":      nil,
		})
	}

	files := form.File["images"]
	if len(files) == 0 {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":      "Error",
            "status_code": fiber.StatusBadRequest,
            "message":     "At least one image is required",
            "result":      nil,
        })
    }

	var fileHeaders []multipart.FileHeader
    for _, file := range files {
        fileHeaders = append(fileHeaders, *file)
    }

	for i := 0; i < len(types); i++ {
		news.Dialog = append(news.Dialog, entities.Dialog{
			Type: types[i],
			Desc: descs[i],
		})
	}

	if news.Title == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "title cannot be empty",
			"result":      nil,
		})
	}

	if len(news.Dialog) == 0 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "dialogs cannot be empty",
			"result":      nil,
		})
	}

	updatedNews, err := c.newsusecase.UpdateNewsByID(id, news, fileHeaders, ctx)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      	"Success",
		"status_code": 	fiber.StatusOK,
		"message":     	"Nursing house retrieved successfully",
		"result":      	updatedNews,
	})
}