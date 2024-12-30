package controllers

import (
	"strconv"
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
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "invalid form data",
			"result":      nil,
		})
	}

	title := form.Value["title"]
	imageTitleFile, err := ctx.FormFile("image_title")
	if err != nil && err != fiber.ErrUnprocessableEntity {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid image_title file",
			"result":      nil,
		})
	}

	imageDescFile, _ := ctx.FormFile("image_desc")
	if len(title) == 0  {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "title cannot be empty",
			"result":      nil,
		})
	}

	req := &entities.News{
		Title: title[0],
	}

	types := form.Value["type"]
	descs := form.Value["desc"]
	bolds := form.Value["bold"]
	if len(types) != len(descs) || len(types) != len(bolds) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Mismatch in count of 'type', 'desc', and 'bold'",
			"result":      nil,
		})
	}

	for i := 0; i < len(types); i++ {
		bold, err := strconv.ParseBool(bolds[i])
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "Invalid value for 'bold', must be true or false",
				"result":      nil,
			})
		}

		req.Dialog = append(req.Dialog, entities.Dialog{
			Type: types[i],
			Desc: descs[i],
			Bold: bold,
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

	data, err := c.newsusecase.CreateNews(req, imageTitleFile, imageDescFile, ctx)
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
		"message":     	"News created successfully",
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
		"message":     	"News retrieved successfully",
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
		"message":     	"News retrieved successfully",
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
		"message":     	"News retrieved successfully",
		"result":      	data,
	})
}

func (c *NewsController) UpdateNewsByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "invalid form data",
			"result":      nil,
		})
	}
	
	title := form.Value["title"]
	imageTitleFile, _ := ctx.FormFile("image_title")
	imageDescFile, _ := ctx.FormFile("image_desc")

	if len(title) == 0 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "title cannot be empty",
			"result":      nil,
		})
	}

	news := entities.News{
		Title: title[0],
	}

	types := form.Value["type"]
	descs := form.Value["desc"]
	bolds := form.Value["bold"]
	if len(types) != len(descs) || len(types) != len(bolds) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Mismatch in count of 'type', 'desc', and 'bold'",
			"result":      nil,
		})
	}

	for i := 0; i < len(types); i++ {
		bold, err := strconv.ParseBool(bolds[i])
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "Invalid value for 'bold', must be true or false",
				"result":      nil,
			})
		}

		news.Dialog = append(news.Dialog, entities.Dialog{
			Type: types[i],
			Desc: descs[i],
			Bold: bold,
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

	updatedNews, err := c.newsusecase.UpdateNewsByID(id, news, imageTitleFile, imageDescFile, ctx)
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
		"message":     	"News retrieved successfully",
		"result":      	updatedNews,
	})
}

func (c *NewsController) DeleteNewsByIDHandler(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    err := c.newsusecase.DeleteNewsByID(id)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":      "Error",
            "status_code": fiber.StatusInternalServerError,
            "message":     err.Error(),
            "result":      nil,
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":      "Success",
        "status_code": fiber.StatusOK,
        "message":     "News deleted successfully",
        "result":      nil,
    })
}
