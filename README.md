# Sipograf (Golang)

Sistem Informasi Posyandu

 Go (Golang) 1.25, MySQL

##  Cara Jalanin
1.  **Clone / Download** repo ini.
2.  **Setup Database:**
    * Buka phpMyAdmin.
    * Hapus database kalau ada: `dbsipograf`.
    * Copy dan Paste di SQL Phpmyadmin (buka file sqltemplate.txt)
3. **Install Dependency:**
    Buka cmd/terminal di folder project, ketik:
    ```bash
    go mod tidy
    ```
5.  **Gas Jalanin:**
    ```bash
    go run main.go
    ```
6.  Buka browser: `http://localhost:8080`

### 1. Login Admin
* **Username:** `admin`
* **Password:** `Abc123`
### 2. Tambah Akun Ortu (Admin Only)
### 3. Login Orang Tua
