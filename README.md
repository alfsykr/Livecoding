## CRUD Produk Asinkron (Golang)

Proyek ini memenuhi permintaan tugas: membuat endpoint `POST /products` dengan MySQL, menjalankan proses asinkron setelah respons dikirim, serta menyiapkan instruksi run dan koleksi Postman.

### 1. Prasyarat
- Go 1.21+
- MySQL 8+ (atau kompatibel)
- Postman (opsional, untuk uji API)

### 2. Skema Database
Jalankan perintah berikut pada MySQL Anda:

```sql
SOURCE schema.sql;
```

File `schema.sql` akan membuat database `item` dan tabel `products`.

### 3. Konfigurasi Lingkungan
Set variabel lingkungan berikut (bisa menggunakan `.env`, PowerShell, atau terminal lain):

```bash
setx APP_PORT 8080
setx DB_DSN "user:pass@tcp(127.0.0.1:3306)/item?parseTime=true"

> Default aplikasi memakai `root:@tcp(127.0.0.1:3306)/item?parseTime=true` (root tanpa password). Bila kredensial Anda berbeda, wajib menimpa `DB_DSN`.
```

> Ganti `user` dan `pass` sesuai kredensial MySQL Anda.

### 4. Menjalankan Aplikasi
1. Pastikan modul Go telah diinisialisasi (`go.mod` tersedia).
2. Jalankan server:
   ```bash
   go run main.go
   ```
3. API akan aktif di `http://localhost:8080`.

### 5. Endpoint
- `POST /products`
  - **Body JSON**
    ```json
    {
      "name": "Keyboard Mekanik",
      "price": 950000,
      "stock": 12
    }
    ```
  - **Respons 201**
    ```json
    {
      "id": 1,
      "name": "Keyboard Mekanik",
      "price": 950000,
      "stock": 12
    }
    ```
  - Setelah respons dikirim, background worker akan:
    1. Mencatat log `Produk Keyboard Mekanik berhasil dibuat. Memicu notifikasi stok.`
    2. Tidur 5 detik lalu mencatat log selesai.
- `GET /`  
  Menampilkan pesan informasi singkat bahwa API berjalan dan contoh penggunaan `POST /products`. Gunakan endpoint ini untuk sanity check bila Anda hanya membuka `http://localhost:8080` di browser dan sebelumnya mendapat 404.

### 6. Koleksi Postman
Impor `postman_collection.json` ke Postman, lalu jalankan request `Create Product`. Koleksi sudah berisi contoh host (`http://localhost:8080`) dan body JSON default.

### 7. Logging & Error Handling
- Middleware sederhana `loggingMiddleware` mencatat setiap request dan lamanya pemrosesan.
- Handler akan mengembalikan kode status sesuai:
  - `405` untuk method selain POST.
  - `400` untuk payload tidak valid atau data kosong.
  - `500` bila gagal menyimpan ke database (cek log server untuk detail).

### 8. Testing Manual
Alternatif cURL:

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Mouse","price":250000,"stock":15}'
```

Server langsung merespons `201 Created`, sementara goroutine di background akan tetap berjalan melakukan logging setelah 5 detik.

