package controllers

import (
	"database/sql"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
	"text/template"
)

func DataAnak(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionData := GetSessionDetails(r)
	role := sessionData["Role"]
	userID := sessionData["UserID"]

	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = config.DB.Query("SELECT t.id_anak, t.id_orangtua, t.nama_anak, t.nama_ibu, t.tempat_lahir, t.tanggal_lahir, t.jenis_kelamin, t.alamat, u.username FROM t_anak t LEFT JOIN users u ON t.id_orangtua = u.id_user")
	} else {
		rows, err = config.DB.Query("SELECT id_anak, id_orangtua, nama_anak, nama_ibu, tempat_lahir, tanggal_lahir, jenis_kelamin, alamat, '' as username FROM t_anak WHERE id_orangtua = ?", userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type AnakWithUser struct {
		models.Anak
		NamaOrangTua string
	}

	var anakList []AnakWithUser
	for rows.Next() {
		var a AnakWithUser
		var namaOrangTua sql.NullString
		rows.Scan(&a.ID, &a.IDOrangtua, &a.NamaAnak, &a.NamaIbu, &a.TempatLahir, &a.TanggalLahir, &a.JenisKelamin, &a.Alamat, &namaOrangTua)
		if namaOrangTua.Valid {
			a.NamaOrangTua = namaOrangTua.String
		} else {
			a.NamaOrangTua = "-"
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

func CreateAnak(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionData := GetSessionDetails(r)
	role := sessionData["Role"]

	var parents []models.User
	if role == "admin" {
		rows, _ := config.DB.Query("SELECT id_user, username FROM users WHERE role = 'orangtua'")
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var u models.User
				rows.Scan(&u.ID, &u.Username)
				parents = append(parents, u)
			}
		}
	}

	data := map[string]interface{}{
		"Role":    role,
		"Parents": parents,
	}

	tmpl := template.Must(template.ParseFiles("views/anak/create.html"))
	tmpl.Execute(w, data)
}

func StoreAnak(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		namaAnak := r.FormValue("nama_anak")
		namaIbu := r.FormValue("nama_ibu")
		tempatLahir := r.FormValue("tempat_lahir")
		tanggalLahir := r.FormValue("tanggal_lahir")
		jenisKelamin := r.FormValue("jenis_kelamin")
		alamat := r.FormValue("alamat")

		sessionData := GetSessionDetails(r)
		role := sessionData["Role"]
		var idOrangtua interface{} = nil

		if role == "orangtua" {
			idOrangtua = sessionData["UserID"]
		} else if role == "admin" {
			formParent := r.FormValue("id_orangtua")
			if formParent != "" {
				idOrangtua = formParent
			}
		}

		_, err := config.DB.Exec("INSERT INTO t_anak (id_orangtua, nama_anak, nama_ibu, tempat_lahir, tanggal_lahir, jenis_kelamin, alamat) VALUES (?, ?, ?, ?, ?, ?, ?)",
			idOrangtua, namaAnak, namaIbu, tempatLahir, tanggalLahir, jenisKelamin, alamat)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
	}
}

func EditAnak(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id_anak")
	var a models.Anak
	err := config.DB.QueryRow("SELECT id_anak, id_orangtua, nama_anak, nama_ibu, tempat_lahir, tanggal_lahir, jenis_kelamin, alamat FROM t_anak WHERE id_anak = ?", id).Scan(&a.ID, &a.IDOrangtua, &a.NamaAnak, &a.NamaIbu, &a.TempatLahir, &a.TanggalLahir, &a.JenisKelamin, &a.Alamat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionData := GetSessionDetails(r)
	role := sessionData["Role"]

	var parents []models.User
	if role == "admin" {
		rows, _ := config.DB.Query("SELECT id_user, username FROM users WHERE role = 'orangtua'")
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var u models.User
				rows.Scan(&u.ID, &u.Username)
				parents = append(parents, u)
			}
		}
	}

	data := map[string]interface{}{
		"Anak":    a,
		"Role":    role,
		"Parents": parents,
	}

	tmpl := template.Must(template.ParseFiles("views/anak/edit.html"))
	tmpl.Execute(w, data)
}

func UpdateAnak(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.URL.Query().Get("id_anak")
		namaAnak := r.FormValue("nama_anak")
		namaIbu := r.FormValue("nama_ibu")
		tempatLahir := r.FormValue("tempat_lahir")
		tanggalLahir := r.FormValue("tanggal_lahir")
		jenisKelamin := r.FormValue("jenis_kelamin")
		alamat := r.FormValue("alamat")

		sessionData := GetSessionDetails(r)

		var err error
		if sessionData["Role"] == "admin" {
			idOrangtua := r.FormValue("id_orangtua")
			var idParent interface{} = nil
			if idOrangtua != "" {
				idParent = idOrangtua
			}
			_, err = config.DB.Exec("UPDATE t_anak SET id_orangtua=?, nama_anak=?, nama_ibu=?, tempat_lahir=?, tanggal_lahir=?, jenis_kelamin=?, alamat=? WHERE id_anak=?",
			idParent, namaAnak, namaIbu, tempatLahir, tanggalLahir, jenisKelamin, alamat, id)
		} else {
			_, err = config.DB.Exec("UPDATE t_anak SET nama_anak=?, nama_ibu=?, tempat_lahir=?, tanggal_lahir=?, jenis_kelamin=?, alamat=? WHERE id_anak=?",
			namaAnak, namaIbu, tempatLahir, tanggalLahir, jenisKelamin, alamat, id)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
	}
}

func DeleteAnak(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id_anak")
	_, err := config.DB.Exec("DELETE FROM t_anak WHERE id_anak=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
}
