# Simple CRUD - Pertemuan 1

Proyek ini adalah aplikasi sederhana CRUD kategori menggunakan Gin (web framework), dengan arsitektur terpisah: repository, service, dan handler. Aplikasi berjalan sebagai HTTP server di port 8080.

## Fitur Utama

- CRUD kategori:
  - List semua kategori
  - Get kategori by ID
  - Create kategori
  - Update kategori
  - Delete kategori
- Health check endpoint untuk monitoring
- Struktur kode modular:
  - `repository` untuk akses data (in-memory array)
  - `service` untuk business logic
  - `handler` untuk HTTP layer (Gin)

## Struktur Proyek

```
pertemuan-1/
  ├─ handler/
  │  └─ category.go           // HTTP handlers (Gin) untuk kategori
  ├─ model/
  │  └─ Category.go           // Definisi struct entity Category (ID, Name)
  ├─ repository/
  │  └─ category.go           // Implementasi in-memory repository Category
  ├─ service/
  │  └─ category.go           // Business logic dan interface CategoryService
  └─ main.go                  // Bootstrap server, wiring repo -> service -> handler, routes
```

Catatan:
- Nama file di `model` diharapkan mendefinisikan `struct` `Category`:
  - Field minimal: `ID int`, `Name string`

## Arsitektur dan Alur

1. `repository.CategoryRepository`
   - Menyediakan operasi data:
     - `GetAllCategories() []model.Category`
     - `GetCategoryByID(id int) (model.Category, error)`
     - `CreateCategory(c model.Category) model.Category`
     - `UpdateCategory(id int, c model.Category) error`
     - `DeleteCategory(id int) error`
   - Implementasi saat ini: in-memory array dengan seed data.

2. `service.CategoryService`
   - Abstraksi business logic:
     - `GetAll() []model.Category`
     - `GetByID(id int) (model.Category, error)`
     - `Create(category model.Category) model.Category`
     - `Update(id int, category model.Category) error`
     - `Delete(id int) error`
   - Menggunakan `repository.CategoryRepository` di balik layar.

3. `handler.CategoryHandler`
   - HTTP endpoints (Gin) yang memanggil `service`:
     - `GetAll(c *gin.Context)`
     - `GetByID(c *gin.Context)`
     - `Create(c *gin.Context)`
     - `Update(c *gin.Context)`
     - `Delete(c *gin.Context)`

4. `main.go`
   - Inisialisasi:
     - `repo := repository.NewCategoryRepository()`
     - `svc := service.NewCategoryService(repo)`
     - `h := handler.NewCategoryHandler(svc)`
   - Routing Gin dan server run di `:8080`.

## Persiapan Lingkungan

- Go 1.20+ (disarankan)
- Module path project: `simple-crud` (disesuaikan dengan `import` yang digunakan)

Jika module path berbeda, pastikan `go.mod` sesuai dan import path di file-file:
- `simple-crud/handler`
- `simple-crud/model`
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
        "id": 4,            // opsional; jika tidak disediakan, repository saat ini tidak auto-generate, jadi pastikan konsisten
        "name": "New Cat"   // wajib
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
    - Response: Status `204 No Content` jika sukses, `400/404` jika gagal
  - DELETE `/categories/:id`
    - Params: `id` (int > 0)
    - Response: Status `204 No Content` jika sukses, `404` jika tidak ditemukan

Contoh curl:

- List categories
  - `curl -s http://localhost:8080/categories | jq`

- Get by id
  - `curl -s http://localhost:8080/categories/1 | jq`

- Create
  - `curl -s -X POST http://localhost:8080/categories -H "Content-Type: application/json" -d '{"id":4,"name":"Sports"}' | jq`

- Update
  - `curl -s -X PUT http://localhost:8080/categories/4 -H "Content-Type: application/json" -d '{"name":"Outdoors"}' -w " HTTP %{http_code}\n"`

- Delete
  - `curl -s -X DELETE http://localhost:8080/categories/4 -w " HTTP %{http_code}\n"`
