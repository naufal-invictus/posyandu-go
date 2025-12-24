package controllers

import (
	"net/http"
	"sipograf-go/config"
	"text/template"
)

// Menampilkan Form Tambah Orang Tua (User + Profile)
func CreateOrangTua(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if GetSessionDetails(r)["Role"] != "admin" {
		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("views/user/create_parent.html"))
	tmpl.Execute(w, nil)
}

func StoreOrangTua(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/create_orangtua", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	namaIbu := r.FormValue("nama_ibu")
	namaAyah := r.FormValue("nama_ayah")
	alamat := r.FormValue("alamat")
	noHP := r.FormValue("no_hp")

	// 1. Insert User
	res, err := config.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, 'orangtua')", username, password)
	if err != nil {
		http.Error(w, "Gagal membuat user (Username mungkin sudah ada): "+err.Error(), 500)
		return
	}

	// 2. Ambil ID User yang baru dibuat
	lastID, _ := res.LastInsertId()

	// 3. Insert Profil Orang Tua
	_, err = config.DB.Exec("INSERT INTO orang_tua (id_user, nama_ibu, nama_ayah, alamat, no_hp) VALUES (?, ?, ?, ?, ?)",
		lastID, namaIbu, namaAyah, alamat, noHP)

	if err != nil {
		http.Error(w, "User dibuat tapi gagal simpan profil ortu: "+err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
}