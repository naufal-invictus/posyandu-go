package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
	"text/template"
)

// Helper untuk mematikan cache browser
func disableCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func DataAnak(w http.ResponseWriter, r *http.Request) {
	disableCache(w)

	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionData := GetSessionDetails(r)
	role := sessionData["Role"]
	userID := sessionData["UserID"]

	query := `SELECT a.id_anak, a.id_orangtua, a.nama_anak, a.tempat_lahir, a.tanggal_lahir, a.jenis_kelamin, 
			  COALESCE(o.nama_ibu, 'Tidak Ada Data') as nama_ibu, COALESCE(o.alamat, '-') as alamat 
			  FROM anak a 
			  LEFT JOIN orang_tua o ON a.id_orangtua = o.id_orangtua`
	
	var rows *sql.Rows
	var err error

	if role == "orangtua" {
		query += " WHERE o.id_user = ?"
		rows, err = config.DB.Query(query, userID)
	} else {
		rows, err = config.DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Database Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var anakList []models.Anak
	for rows.Next() {
		var a models.Anak
		var idOrtu sql.NullInt64
		rows.Scan(&a.ID, &idOrtu, &a.NamaAnak, &a.TempatLahir, &a.TanggalLahir, &a.JenisKelamin, &a.NamaIbu, &a.Alamat)
		if idOrtu.Valid {
			a.IDOrangtua = int(idOrtu.Int64)
		}
		anakList = append(anakList, a)
	}

	data := map[string]interface{}{
		"Anak": anakList,
		"Role": role,
	}

	tmpl := template.Must(template.ParseFiles("views/anak/data.html"))
	tmpl.Execute(w, data)
}

// FIX: CreateAnak dengan Error Handling yang lebih baik
func CreateAnak(w http.ResponseWriter, r *http.Request) {
	disableCache(w)
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	
	// Ambil data orang tua untuk dropdown
	rows, err := config.DB.Query("SELECT id_orangtua, nama_ibu, alamat FROM orang_tua ORDER BY nama_ibu ASC")
	if err != nil {
		log.Println("Error fetching parents:", err) // Log error ke terminal
	}

	var parents []models.Parent
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var p models.Parent
			// Scan sesuai struct (ID, NamaIbu, Alamat)
			if err := rows.Scan(&p.ID, &p.NamaIbu, &p.Alamat); err == nil {
				parents = append(parents, p)
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("views/anak/create.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Parents": parents, 
		"Role": GetSessionDetails(r)["Role"],
	})
}

func StoreAnak(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var idOrtu interface{}
		idForm := r.FormValue("id_orangtua")
		
		if idForm != "" {
			idOrtu = idForm
		} else {
			sessionData := GetSessionDetails(r)
			var tempID int
			err := config.DB.QueryRow("SELECT id_orangtua FROM orang_tua WHERE id_user = ?", sessionData["UserID"]).Scan(&tempID)
			if err == nil {
				idOrtu = tempID
			} else {
				idOrtu = nil 
			}
		}

		_, err := config.DB.Exec("INSERT INTO anak (id_orangtua, nama_anak, tempat_lahir, tanggal_lahir, jenis_kelamin) VALUES (?, ?, ?, ?, ?)",
			idOrtu, r.FormValue("nama_anak"), r.FormValue("tempat_lahir"), r.FormValue("tanggal_lahir"), r.FormValue("jenis_kelamin"))
		
		if err != nil {
			http.Error(w, "Gagal Simpan: "+err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
	}
}

func EditAnak(w http.ResponseWriter, r *http.Request) {
	disableCache(w)
	
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id_anak")
	var a models.Anak
	var idOrtu sql.NullInt64

	query := `
		SELECT a.id_anak, a.id_orangtua, a.nama_anak, a.tempat_lahir, a.tanggal_lahir, a.jenis_kelamin, 
		COALESCE(o.nama_ibu, 'Tidak Ada Data') as nama_ibu, COALESCE(o.alamat, '-') as alamat 
		FROM anak a 
		LEFT JOIN orang_tua o ON a.id_orangtua = o.id_orangtua
		WHERE a.id_anak = ?`

	err := config.DB.QueryRow(query, id).Scan(
			&a.ID, &idOrtu, &a.NamaAnak, &a.TempatLahir, &a.TanggalLahir, &a.JenisKelamin, &a.NamaIbu, &a.Alamat)

	if err != nil {
		http.Error(w, "Data Anak Tidak Ditemukan", 404)
		return
	}

	if idOrtu.Valid {
		a.IDOrangtua = int(idOrtu.Int64)
	}

	rows, _ := config.DB.Query("SELECT id_orangtua, nama_ibu, alamat FROM orang_tua ORDER BY nama_ibu ASC")
	var parents []models.Parent
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var p models.Parent
			rows.Scan(&p.ID, &p.NamaIbu, &p.Alamat)
			parents = append(parents, p)
		}
	}

	data := map[string]interface{}{
		"Anak":    a,
		"Role":    GetSessionDetails(r)["Role"],
		"Parents": parents,
	}

	tmpl := template.Must(template.ParseFiles("views/anak/edit.html"))
	tmpl.Execute(w, data)
}

func UpdateAnak(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.URL.Query().Get("id_anak")
		
		_, err := config.DB.Exec("UPDATE anak SET nama_anak=?, tempat_lahir=?, tanggal_lahir=?, jenis_kelamin=? WHERE id_anak=?",
			r.FormValue("nama_anak"), r.FormValue("tempat_lahir"), r.FormValue("tanggal_lahir"), r.FormValue("jenis_kelamin"), id)

		if r.FormValue("id_orangtua") != "" {
			config.DB.Exec("UPDATE anak SET id_orangtua=? WHERE id_anak=?", r.FormValue("id_orangtua"), id)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
	}
}

func DeleteAnak(w http.ResponseWriter, r *http.Request) {
	disableCache(w)
	
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id_anak")
	
	_, err := config.DB.Exec("DELETE FROM anak WHERE id_anak=?", id)
	
	if err != nil {
		http.Error(w, "Gagal Hapus: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
}