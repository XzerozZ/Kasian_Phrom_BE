package controllers

import (
	"mime/multipart"
  	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/usecases"

	"github.com/gofiber/fiber/v2"
)

type NhController struct {
	nhusecase usecases.NhUseCase
}

func NewNhController(nhusecase usecases.NhUseCase) *NhController {
	return &NhController{nhusecase: nhusecase}
}

func (c *NhController) CreateNhHandler(ctx *fiber.Ctx) error {
	var nursingHouse entities.NursingHouse
	if err := ctx.BodyParser(&nursingHouse); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      	fiber.ErrBadRequest.Message,
			"status_code": 	fiber.ErrBadRequest.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"message":     "Failed to parse form data",
			"status_code": fiber.StatusBadRequest,
		})
	}

	files := form.File["images"]
	var fileHeaders []multipart.FileHeader
    for _, file := range files {
        fileHeaders = append(fileHeaders, *file)
    }

	data, err := c.nhusecase.CreateNh(nursingHouse, fileHeaders, ctx)
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

func (c *NhController) GetAllNhHandler(ctx *fiber.Ctx) error {
	data, err := c.nhusecase.GetAllNh()
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

func (c *NhController) GetAllActiveNhHandler(ctx *fiber.Ctx) error {
	data, err := c.nhusecase.GetActiveNh()
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

func (c *NhController) GetAllInactiveNhHandler(ctx *fiber.Ctx) error {
	data, err := c.nhusecase.GetInactiveNh()
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

func (c *NhController) GetNhByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := c.nhusecase.GetNhByID(id)
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

func (c *NhController) GetNhNextIDHandler(ctx *fiber.Ctx) error {
	data, err := c.nhusecase.GetNhNextID()
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

func (c *NhController) UpdateNhByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var nursingHouse entities.NursingHouse
	if err := ctx.BodyParser(&nursingHouse); err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      	fiber.ErrNotFound.Message,
			"status_code": 	fiber.ErrNotFound.Code,
			"message":     	err.Error(),
			"result":      	nil,
		})
	}

	form, err := ctx.MultipartForm()
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":      "Error",
            "status_code": fiber.StatusBadRequest,
            "message":     "Failed to parse form data",
            "result":      nil,
        })
    }

    files := form.File["images"]
    var fileHeaders []multipart.FileHeader
    for _, file := range files {
        fileHeaders = append(fileHeaders, *file)
    }

	updatedNh, err := c.nhusecase.UpdateNhByID(id, nursingHouse, fileHeaders, ctx)
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
		"result":      	updatedNh,
	})
}

func (c *NhController) GetNhImages(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	images, err := c.nhusecase.GetImagesByNhID(id)
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
		"result":      	images,
	})
}