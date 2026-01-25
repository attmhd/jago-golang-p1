package repository

import (
	"errors"
	"simple-crud/model"
)

type CategoryRepository interface {
	GetAllCategories() []model.Category
	GetCategoryByID(id int) (model.Category, error)
	CreateCategory(c model.Category) model.Category
	UpdateCategory(id int, c model.Category) error
	DeleteCategory(id int) error
}

type categoryRepository struct {
	categories []model.Category
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{
		categories: []model.Category{
			{ID: 1, Name: "Electronics"},
			{ID: 2, Name: "Clothing"},
			{ID: 3, Name: "Books"},
		},
	}
}

func (r *categoryRepository) GetAllCategories() []model.Category {
	return r.categories
}

func (r *categoryRepository) GetCategoryByID(id int) (model.Category, error) {
	for _, c := range r.categories {
		if c.ID == id {
			return c, nil
		}
	}
	return model.Category{}, errors.New("category not found")
}

func (r *categoryRepository) CreateCategory(c model.Category) model.Category {
	r.categories = append(r.categories, c)
	return c
}

func (r *categoryRepository) UpdateCategory(id int, c model.Category) error {
	for i, category := range r.categories {
		if category.ID == id {
			r.categories[i] = c
			return nil
		}
	}
	return errors.New("category not found")
}

func (r *categoryRepository) DeleteCategory(id int) error {
	for i, category := range r.categories {
		if category.ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return nil
		}
	}
	return errors.New("category not found")
}
