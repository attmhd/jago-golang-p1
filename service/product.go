package service

import (
	model "simple-crud/models"
	"simple-crud/repository"
)

type ProductServices interface {
	GetAll() ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	Create(product *model.Product) (*model.Product, error)
	Update(product *model.Product) (*model.Product, error)
	Delete(id int) error
}

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product *model.Product) (*model.Product, error) {
	created, err := s.repo.Create(product)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(created.ID)
}

func (s *ProductService) Update(product *model.Product) (*model.Product, error) {
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return s.repo.GetByID(product.ID)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
