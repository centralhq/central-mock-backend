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

func (service *ShapeService) GetShape() *ShapeObject {
	return service.repo.GetShape()
}

func (service *ShapeService) SetShape(shape string) *string {
	return service.repo.SetShape(shape)
}

func (service *ShapeService) SetColor(color string) *string {
	return service.repo.SetColor(color)
}

func (service *ShapeService) SetSize(size string) *string {
	return service.repo.SetSize(size)
}