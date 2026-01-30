package service

import (
	model "simple-crud/models"
	"simple-crud/repository"
)

type CategoriesService interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category model.Category) model.Category
	Update(id int, category model.Category) error
	Delete(id int) error
}

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Create(category model.Category) (model.Category, error) {
	c, err := s.repo.Create(category)
	if err != nil {
		return model.Category{}, err
	}
	return *c, nil
}

func (s *CategoryService) Update(id int, category model.Category) error {
	return s.repo.Update(id, category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
