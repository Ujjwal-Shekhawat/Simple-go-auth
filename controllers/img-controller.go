package controller

import (
	"github.com/Ujjwal-Shekhawat/golang-gin-poc/model"
	"github.com/Ujjwal-Shekhawat/golang-gin-poc/service"
	"github.com/gin-gonic/gin"
)

// ImgController - Image Controller Interface
type ImgController interface {
	FindAll() []model.Image
	Save(ctx *gin.Context) model.Image
}

type controller struct {
	service service.ImgService
}

func New(service service.ImgService) ImgController {
	return &controller{
		service: service,
	}
}

func (c *controller) FindAll() []model.Image {
	return c.service.FindAll()
}

func (c *controller) Save(ctx *gin.Context) model.Image {
	var image model.Image
	ctx.BindJSON(&image)
	c.service.Save(image)

	return image
}
