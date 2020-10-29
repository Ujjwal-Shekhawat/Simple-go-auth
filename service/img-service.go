package service

import "github.com/Ujjwal-Shekhawat/golang-gin-poc/model"

type ImgService interface {
	Save(model.Image) model.Image
	FindAll() []model.Image
}

type imgService struct {
	images []model.Image
}

// New - returns a pointer to the struct
func New() ImgService {
	return &imgService{}
}

func (service *imgService) Save(image model.Image) model.Image {
	service.images = append(service.images, image)
	return image
}

func (service *imgService) FindAll() []model.Image {
	return service.images
}
