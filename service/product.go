package service

import (
	"go-boot-category-api/framework/repository"
	"go-boot-category-api/model"
)

type Product interface {
	GetAll() ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id int) error
}

type productService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) Product {
	return &productService{repo: repo}
}

func (s *productService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) Create(data *model.Product) error {
	return s.repo.Create(data)
}

func (s *productService) GetByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) Update(product *model.Product) error {
	return s.repo.Update(product)
}

func (s *productService) Delete(id int) error {
	return s.repo.Delete(id)
}
