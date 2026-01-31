# Simple CRUD - Pertemuan 2

Proyek ini adalah aplikasi sederhana CRUD untuk Products dan Categories menggunakan Gin (web framework), dengan arsitektur terpisah: repository, service, dan handler. Aplikasi berjalan sebagai HTTP server di port 8080.

Kini data `products` terhubung dengan `categories` melalui `category_id` dan seluruh response untuk `GET /products`, `GET /products/:id`, `POST /products`, dan `PUT /products/:id` mengembalikan struktur JSON dengan objek `category` yang nested (berisi `id` dan `name`) berdasarkan relasi.

## Fitur Utama

- CRUD Categories:
  - List semua kategori
  - Get kategori by ID
  - Create kategori
  - Update kategori
  - Delete kategori
- CRUD Products:
  - List semua produk dengan kategori nested
  - Get produk by ID dengan kategori nested
  - Create produk dan kembalikan kategori nested
  - Update produk dan kembalikan kategori nested
  - Delete produk
- Health check endpoint untuk monitoring
- Struktur kode modular:
  - `repository` untuk akses data (PostgreSQL via `database/sql`)
  - `service` untuk business logic
  - `handler` untuk HTTP layer (Gin)
  
## Struktur Proyek

```
pertemuan-1/
  ├─ handler/
  │  ├─ category.go           // HTTP handlers (Gin) untuk kategori
  │  ├─ product.go            // HTTP handlers (Gin) untuk produk, response nested category
  │  └─ response.go           // DTO untuk response Product/Category
  ├─ models/
  │  ├─ category.go           // Definisi struct entity Category
  │  └─ product.go            // Definisi struct entity Product (termasuk CategoryID & CategoryName)
  ├─ repository/
  │  ├─ category.go           // Implementasi repository Category (DB)
  │  └─ product.go            // Implementasi repository Product (DB, JOIN category)
  ├─ service/
  │  ├─ category.go           // Business logic dan interface CategoryService
  │  └─ product.go            // Business logic dan interface ProductService
  └─ main.go                  // Bootstrap server, wiring repo -> service -> handler, routes
```

## Arsitektur dan Alur

1. Repository
   - `repository.ProductRepository`
     - `GetAll() ([]model.Product, error)` — SELECT dengan JOIN ke `categories` untuk mendapatkan `CategoryName`
     - `GetByID(id int) (*model.Product, error)` — SELECT dengan JOIN untuk satu produk
     - `Create(product *model.Product) (*model.Product, error)` — INSERT dengan `RETURNING id`, lalu service akan memanggil `GetByID` untuk melengkapi `CategoryName`
     - `Update(product *model.Product) error` — UPDATE berdasarkan `id`
     - `Delete(id int) error` — DELETE berdasarkan `id`
   - `repository.CategoryRepository` — operasi dasar kategori (GetAll, GetByID, Create, Update, Delete)

2. Service
   - `service.ProductService`
     - `GetAll() ([]model.Product, error)`
     - `GetByID(id int) (*model.Product, error)`
     - `Create(product *model.Product) (*model.Product, error)` — setelah insert, fetch lagi dengan `GetByID` agar `CategoryName` terisi
     - `Update(product *model.Product) (*model.Product, error)` — setelah update, fetch lagi dengan `GetByID` agar `CategoryName` terisi
     - `Delete(id int) error`
   - `service.CategoryService` — abstraksi business logic kategori

3. Handler
   - `handler.ProductHandler`
     - Mengubah hasil dari service (model flat) menjadi DTO dengan `category` nested:
       ```
       {
         "message": "Success",
         "data": [
           {
             "id": 1,
             "name": "Makanan",
             "price": 10000,
             "stock": 80,
             "category": {
               "id": 1,
               "name": "Test"
             }
           }
         ]
       }
       ```
     - Struktur response yang sama diterapkan untuk `GetById`, `Create`, dan `Update`
   - `handler.CategoryHandler` — endpoints dasar kategori

4. `main.go`
   - Inisialisasi:
     - `db` (sql.DB) untuk koneksi PostgreSQL
     - `repo := repository.NewProductRepository(db)` dan `repository.NewCategoryRepository(db)`
     - `svc := service.NewProductService(repo)` dan `service.NewCategoryService(repo)`
     - `h := handler.NewProductHandler(svc)` dan `handler.NewCategoryHandler(svc)`
   - Routing Gin dan server run di `:8080`.

## Persiapan Lingkungan

- Go 1.20+ (disarankan)
- Database: PostgreSQL (atau kompatibel dengan sintaks `RETURNING`)
- `database/sql` dan driver PostgreSQL (mis. `github.com/lib/pq`) di `go.mod`
- Module path project: `simple-crud` (disesuaikan dengan `import` yang digunakan)

Pastikan `go.mod` sesuai dan import path di file-file:
- `simple-crud/handler`
- `simple-crud/models`
- `simple-crud/repository`
- `simple-crud/service`

## Menjalankan Aplikasi

1. Posisikan terminal di direktori `pertemuan-1`.
2. Jalankan:
   - `go run main.go`

Server akan berjalan di `http://localhost:8080`.

## API Endpoints

- Health
  - GET `/health`
    - Response: `{"status":"ok"}`

- Categories
  - GET `/categories`
    - Response: `[]Category`
  - GET `/categories/:id`
    - Params: `id` (int > 0)
    - Response: `Category` atau `404` saat tidak ditemukan
  - POST `/categories`
    - Body JSON:
      ```
      {
        "name": "New Cat"
      }
      ```
    - Response: `Category` yang dibuat
  - PUT `/categories/:id`
    - Params: `id` (int > 0)
    - Body JSON:
      ```
      {
        "name": "Updated Name"
      }
      ```
    - Response: `Category` yang diperbarui
  - DELETE `/categories/:id`
    - Params: `id` (int > 0)
    - Response: Status OK jika sukses, `404` jika tidak ditemukan

- Products
  - GET `/products`
    - Response: daftar produk dengan kategori nested
      ```
      {
        "message": "Success",
        "data": [
          {
            "id": 1,
            "name": "Makanan",
            "price": 10000,
            "stock": 80,
            "category": {
              "id": 1,
              "name": "Test"
            }
          }
        ]
      }
      ```
  - GET `/products/:id`
    - Params: `id` (int > 0)
    - Response: satu produk dengan kategori nested atau `404`
  - POST `/products`
    - Body JSON:
      ```
      {
        "category_id": 1,
        "name": "Minuman",
        "price": 5000,
        "stock": 30
      }
      ```
    - Proses: INSERT, lalu service akan `GetByID` untuk melengkapi `category.name`
    - Response: produk yang dibuat dengan kategori nested
  - PUT `/products/:id`
    - Params: `id` (int > 0)
    - Body JSON:
      ```
      {
        "category_id": 2,
        "name": "Minuman Segar",
        "price": 6000,
        "stock": 40
      }
      ```
    - Proses: UPDATE, lalu service akan `GetByID` untuk melengkapi `category.name`
    - Response: produk yang diperbarui dengan kategori nested
  - DELETE `/products/:id`
    - Params: `id` (int > 0)
    - Response: Status OK jika sukses, `404` jika tidak ditemukan

## Contoh curl

- List products
  - `curl -s http://localhost:8080/products | jq`

- Get product by id
  - `curl -s http://localhost:8080/products/1 | jq`

- Create product
  - `curl -s -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"category_id":1,"name":"Minuman","price":5000,"stock":30}' | jq`

- Update product
  - `curl -s -X PUT http://localhost:8080/products/1 -H "Content-Type: application/json" -d '{"category_id":2,"name":"Minuman Segar","price":6000,"stock":40}' | jq`

- Delete product
  - `curl -s -X DELETE http://localhost:8080/products/1 -w " HTTP %{http_code}\n"`

## Catatan

- Pastikan `categories` berisi data yang valid sebelum membuat `products`, karena `category_id` harus merujuk ke `categories.id`.
- Implementasi repository `products` menggunakan JOIN untuk mengisi `CategoryName`. Service `Create` dan `Update` akan memanggil `GetByID` setelah operasi tulis untuk memastikan respons memiliki `category.name` yang benar.
- Jika database bukan PostgreSQL, sesuaikan cara mendapatkan `ID` hasil insert (misalnya dengan `LastInsertId()` jika driver mendukung).
