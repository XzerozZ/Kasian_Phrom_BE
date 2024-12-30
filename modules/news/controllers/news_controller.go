package controllers

import (
	"strings"
	"strconv"
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
	bolds := form.Value["bold"]
	if len(types) != len(descs) || len(types) != len(bolds) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Mismatch in count of 'type', 'desc', and 'bold'",
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
	bolds := form.Value["bold"]
	if len(types) != len(descs) || len(types) != len(bolds) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Mismatch in count of 'type', 'desc', and 'bold'",
			"result":      nil,
		})
	}

	var deleteImages []string
    if imagesStr := ctx.FormValue("delete_images"); imagesStr != "" {
		deleteImages = strings.Split(imagesStr, ",")
	} else {
		deleteImages = []string{} 
	}

	var fileHeaders []multipart.FileHeader
    if files := form.File["images"]; len(files) > 0 {
        for _, file := range files {
            fileHeaders = append(fileHeaders, *file)
        }
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

	if len(deleteImages) > 0 || len(fileHeaders) > 0 {
        existingNews, err := c.newsusecase.GetNewsByID(id)
        if err != nil {
            return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":      "Error",
                "status_code": fiber.StatusNotFound,
                "message":     "News not found",
                "result":      nil,
            })
        }
        
        remainingImagesCount := len(existingNews.Images) - len(deleteImages) + len(fileHeaders)
        if len(fileHeaders) > 0 {
			remainingImagesCount += len(fileHeaders)
		}
		
		if remainingImagesCount < 1 {
            return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":      "Error",
                "status_code": fiber.StatusBadRequest,
                "message":     "News must have at least one image",
                "result":      nil,
            })
        }
    }

	if len(news.Dialog) == 0 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "dialogs cannot be empty",
			"result":      nil,
		})
	}

	updatedNews, err := c.newsusecase.UpdateNewsByID(id, news, fileHeaders, deleteImages, ctx)
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
