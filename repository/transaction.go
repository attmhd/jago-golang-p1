package repository

import (
	"database/sql"
	"fmt"
	"simple-crud/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productName string
		var productID, price, stock int
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		subtotal := item.Quantity * price
		totalAmount += subtotal

		// Update hanya stock (tanpa kolom sold)
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	// Asumsi kolom created_at memiliki default NOW()
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}

// Mengembalikan produk terlaris hari ini (nama + qty terjual) dengan menghitung dari transaction_details
func (r *TransactionRepository) GetTopSellingProduct() (*models.TopSellingProduct, error) {
	var (
		name    string
		qtySold int
	)

	// Hitung qty terjual per produk dari transaction_details yang terjadi hari ini
	// Bergantung pada kolom created_at di tabel transactions
	query := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) AS qty_terjual
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query).Scan(&name, &qtySold)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no products found")
	}
	if err != nil {
		return nil, err
	}

	return &models.TopSellingProduct{
		Name:    name,
		QtySold: qtySold,
	}, nil
}

// Ringkasan penjualan "hari ini": total revenue dan jumlah transaksi
func (r *TransactionRepository) GetSalesSummary() (int, int, error) {
	var totalRevenue int
	var totalTransaksi int

	// Hitung total revenue hari ini
	err := r.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE DATE(created_at) = CURRENT_DATE").Scan(&totalRevenue)
	if err != nil {
		return 0, 0, err
	}

	// Hitung total transaksi hari ini
	err = r.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE DATE(created_at) = CURRENT_DATE").Scan(&totalTransaksi)
	if err != nil {
		return 0, 0, err
	}

	return totalRevenue, totalTransaksi, nil
}

// Ringkasan penjualan berdasarkan rentang tanggal [startDate, endDate] (format: YYYY-MM-DD)
func (r *TransactionRepository) GetSalesSummaryByRange(startDate, endDate string) (int, int, error) {
	var totalRevenue int
	var totalTransaksi int

	// Hitung total revenue pada rentang tanggal
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&totalRevenue)
	if err != nil {
		return 0, 0, err
	}

	// Hitung total transaksi pada rentang tanggal
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&totalTransaksi)
	if err != nil {
		return 0, 0, err
	}

	return totalRevenue, totalTransaksi, nil
}

// Produk terlaris berdasarkan rentang tanggal [startDate, endDate]
func (r *TransactionRepository) GetTopSellingProductByRange(startDate, endDate string) (*models.TopSellingProduct, error) {
	var name string
	var qtySold int

	query := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) AS qty_terjual
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, startDate, endDate).Scan(&name, &qtySold)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no products found in range")
	}
	if err != nil {
		return nil, err
	}

	return &models.TopSellingProduct{
		Name:    name,
		QtySold: qtySold,
	}, nil
}
