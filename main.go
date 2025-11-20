package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Product mewakili data produk yang tersimpan di database.
type Product struct {
	ID    int64   `json:"id,omitempty"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// App menyimpan dependensi utama aplikasi.
type App struct {
	db             *sql.DB
	notificationCh chan Product
}

func main() {
	cfg := loadConfig()

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		log.Fatalf("gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}

	app := &App{
		db:             db,
		notificationCh: make(chan Product, 100),
	}

	// Worker goroutine untuk simulasi tugas asynchronous.
	go app.backgroundNotifier()

	mux := http.NewServeMux()
	mux.HandleFunc("/products", app.handleCreateProduct)
	mux.HandleFunc("/", handleRoot)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown.
	go func() {
		log.Printf("server berjalan di http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	waitForShutdown(server, app)
}

// Config menyimpan konfigurasi runtime.
type Config struct {
	Port string
	DSN  string
}

func loadConfig() Config {
	port := os.Getenv("APP_PORT")
	if strings.TrimSpace(port) == "" {
		port = "8080"
	}

	dsn := os.Getenv("DB_DSN")
	if strings.TrimSpace(dsn) == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/item?parseTime=true"
	}

	return Config{
		Port: port,
		DSN:  dsn,
	}
}

// handleCreateProduct menangani POST /products.
func (a *App) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "payload tidak valid", http.StatusBadRequest)
		return
	}

	if err := validateProduct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertedID, err := a.insertProduct(r.Context(), req)
	if err != nil {
		log.Printf("gagal menyimpan produk: %v", err)
		http.Error(w, "gagal menyimpan produk", http.StatusInternalServerError)
		return
	}

	req.ID = insertedID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(req); err != nil {
		log.Printf("gagal menulis respons: %v", err)
	}

	select {
	case a.notificationCh <- req:
	default:
		// Jika channel penuh, jangan blokir request utama.
		go func(p Product) { a.notificationCh <- p }(req)
	}
}

func validateProduct(p Product) error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("nama produk wajib diisi")
	}
	if p.Price < 0 {
		return errors.New("harga tidak boleh bernilai negatif")
	}
	if p.Stock < 0 {
		return errors.New("stok tidak boleh bernilai negatif")
	}
	return nil
}

func (a *App) insertProduct(ctx context.Context, p Product) (int64, error) {
	query := `INSERT INTO products (name, price, stock) VALUES (?, ?, ?)`
	result, err := a.db.ExecContext(ctx, query, p.Name, p.Price, p.Stock)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (a *App) backgroundNotifier() {
	for product := range a.notificationCh {
		log.Printf("Produk %s berhasil dibuat. Memicu notifikasi stok.", product.Name)
		time.Sleep(5 * time.Second)
		log.Printf("Notifikasi stok untuk produk %s (ID: %d) selesai diproses.", product.Name, product.ID)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s selesai dalam %s", r.Method, r.URL.Path, duration)
	})
}

func waitForShutdown(server *http.Server, app *App) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("memulai proses shutdown...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("gagal shutdown server: %v", err)
	}

	close(app.notificationCh)
}

// parsePrice membantu membaca string harga ke float64 (opsional untuk future use).
func parsePrice(priceStr string) (float64, error) {
	priceStr = strings.TrimSpace(priceStr)
	return strconv.ParseFloat(priceStr, 64)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{
		"message":       "Produk API berjalan",
		"usage_example": "POST /products dengan JSON {\"name\":\"Mouse\",\"price\":250000,\"stock\":15}",
	}
	_ = json.NewEncoder(w).Encode(resp)
}

