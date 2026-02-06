package repository

import (
	"database/sql"
	model "simple-crud/models"
)

type ProductRepositories interface {
	GetAll(name string) ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	Create(product *model.Product) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id int) error
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetAll(name string) ([]model.Product, error) {
	query := `
	SELECT
		p.id,
		p.category_id,
		c.name AS category_name,
		p.name,
		p.price,
		p.stock
	FROM products p
	JOIN categories c
		ON p.category_id = c.id
`

	args := []interface{}{}
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	// terminate query after optional WHERE
	query += ";"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.CategoryName,
			&product.Name,
			&product.Price,
			&product.Stock,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	query := `
		SELECT
			p.id,
			p.category_id,
			c.name AS category_name,
			p.name,
			p.price,
			p.stock
		FROM products p
		JOIN categories c
			ON p.category_id = c.id
		WHERE p.id = $1;
	`
	row := r.db.QueryRow(query, id)

	var product model.Product
	if err := row.Scan(
		&product.ID,
		&product.CategoryID,
		&product.CategoryName,
		&product.Name,
		&product.Price,
		&product.Stock,
	); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Create(product *model.Product) (*model.Product, error) {
	query := `
		INSERT INTO products (category_id, name, price, stock)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	row := r.db.QueryRow(query, product.CategoryID, product.Name, product.Price, product.Stock)
	if err := row.Scan(&product.ID); err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) Update(product *model.Product) error {
	query := `
		UPDATE products
		SET category_id = $2, name = $3, price = $4, stock = $5
		WHERE id = $1;
	`
	_, err := r.db.Exec(query, product.ID, product.CategoryID, product.Name, product.Price, product.Stock)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Delete(id int) error {
	query := `
		DELETE FROM products
		WHERE id = $1;
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
