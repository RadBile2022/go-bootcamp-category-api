package repository

import (
	"database/sql"
	"errors"
	"go-boot-category-api/model"
)

type Category interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category *model.Category) error
	Update(category *model.Category) error
	Delete(id int) error
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategory(db *sql.DB) Category {
	return &categoryRepo{db: db}
}

func (repo *categoryRepo) GetAll() ([]model.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var p model.Category
		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, p)
	}

	return categories, nil
}

func (repo *categoryRepo) Create(category *model.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}

// GetByID - ambil kategori by ID
func (repo *categoryRepo) GetByID(id int) (*model.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var p model.Category
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *categoryRepo) Update(category *model.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return nil
}

func (repo *categoryRepo) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return err
}
