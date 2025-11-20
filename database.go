package main

import (
    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
    var err error

    DB, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        log.Fatal("gagal membuat koneksi:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("gagal terhubung ke database:", err)
    }

    log.Println("Database connected!")
}
