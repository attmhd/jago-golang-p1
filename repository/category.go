package repository

import (
	"database/sql"
	"errors"
	model "simple-crud/models"
)

type CategoriesRepository interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(c model.Category) (*model.Category, error)
	Update(id int, c model.Category) error
	Delete(id int) error
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*model.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"
	var c model.Category
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Kategori tidak ditemukan")
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(c model.Category) (*model.Category, error) {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	var id int
	err := r.db.QueryRow(query, c.Name, c.Description).Scan(&id)
	if err != nil {
		return nil, err
	}

	c.ID = id
	return &c, nil
}

func (r *CategoryRepository) Update(id int, c model.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := r.db.Exec(query, c.Name, c.Description, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Kategori tidak ditemukan")
	}

	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
