package service

import (
	"simple-crud/models"
	"simple-crud/repository"
	"simple-crud/util"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetSalesSummary() (*util.SalesSummary, error) {
	totalRevenue, totalTransaksi, err := s.repo.GetSalesSummary()
	if err != nil {
		return nil, err
	}
	topSellingProduct, err := s.repo.GetTopSellingProduct()
	if err != nil {
		return nil, err
	}

	var summary util.SalesSummary
	summary.TotalRevenue = totalRevenue
	summary.TotalTransaksi = totalTransaksi
	summary.ProdukTerlaris = util.ProdukTerlaris{
		Nama:       topSellingProduct.Name,
		QtyTerjual: topSellingProduct.QtySold,
	}

	return &summary, nil
}

// Expose top selling product for handler usage
func (s *TransactionService) GetTopSellingProduct() (*models.TopSellingProduct, error) {
	return s.repo.GetTopSellingProduct()
}

// GetSalesSummaryByRange menerima start_date dan end_date (format: YYYY-MM-DD) untuk menampilkan ringkasan pada rentang tanggal.
// Jika salah satu parameter kosong, gunakan perhitungan "hari ini" dengan GetSalesSummary().
func (s *TransactionService) GetSalesSummaryByRange(startDate, endDate string) (*util.SalesSummary, error) {
	if startDate == "" || endDate == "" {
		return s.GetSalesSummary()
	}

	totalRevenue, totalTransaksi, err := s.repo.GetSalesSummaryByRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	topSellingProduct, err := s.repo.GetTopSellingProductByRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &util.SalesSummary{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: util.ProdukTerlaris{
			Nama:       topSellingProduct.Name,
			QtyTerjual: topSellingProduct.QtySold,
		},
	}, nil
}

// GetTopSellingProductByRange mengembalikan produk terlaris pada rentang tanggal.
// Jika salah satu parameter kosong, fallback ke produk terlaris hari ini.
func (s *TransactionService) GetTopSellingProductByRange(startDate, endDate string) (*models.TopSellingProduct, error) {
	if startDate == "" || endDate == "" {
		return s.GetTopSellingProduct()
	}
	return s.repo.GetTopSellingProductByRange(startDate, endDate)
}
