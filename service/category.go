package service

import (
	"simple-crud/model"
	"simple-crud/repository"
)

type CategoryService interface {
	GetAll() []model.Category
	GetByID(id int) (model.Category, error)
	Create(category model.Category) model.Category
	Update(id int, category model.Category) error
	Delete(id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) GetAll() []model.Category {
	return s.repo.GetAllCategories()
}

func (s *categoryService) GetByID(id int) (model.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *categoryService) Create(category model.Category) model.Category {
	return s.repo.CreateCategory(category)
}

func (s *categoryService) Update(id int, category model.Category) error {
	return s.repo.UpdateCategory(id, category)
}

func (s *categoryService) Delete(id int) error {
	return s.repo.DeleteCategory(id)
}
