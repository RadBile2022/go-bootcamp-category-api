package service

import (
	"go-boot-category-api/framework/repository"
	"go-boot-category-api/model"
)

type Category interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category *model.Category) error
	Update(category *model.Category) error
	Delete(id int) error
}

type categoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) Category {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) Create(data *model.Category) error {
	return s.repo.Create(data)
}

func (s *categoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) Update(Category *model.Category) error {
	return s.repo.Update(Category)
}

func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
