package main

type ShapeService struct {
	config *Config
	repo *ShapeRepo
}

func NewShapeService(config *Config, repo *ShapeRepo) *ShapeService {
	return &ShapeService{
		config: config,
		repo: repo,
	}
}

func (service *ShapeService) CreateShape(uid string, shape string, color string, size string) *int8 {
	return service.repo.CreateShape(uid, shape, color, size)
}

func (service *ShapeService) DeleteShape(uid string) *int8 {
	return service.repo.DeleteShape(uid)
}

func (service *ShapeService) GetShape() *ShapeObject {
	return service.repo.GetShape()
}

func (service *ShapeService) SetShape(uid string, shape string) *int8 {
	return service.repo.SetShape(uid, shape)
}

func (service *ShapeService) SetColor(uid string, color string) *int8 {
	return service.repo.SetColor(uid, color)
}

func (service *ShapeService) SetSize(uid string, size string) *int8 {
	return service.repo.SetSize(uid, size)
}