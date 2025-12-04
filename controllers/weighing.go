package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
	"text/template"
)

func DataPenimbangan(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idAnak := r.URL.Query().Get("id_anak")
	sessionData := GetSessionDetails(r)

	rows, err := config.DB.Query("SELECT id_penimbangan, id_anak, tgl_penimbangan, umur, berat_badan, tinggi_badan, jenis_imunisasi, keterangan, petugas FROM t_penimbangan WHERE id_anak = ?", idAnak)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []models.Penimbangan
	for rows.Next() {
		var p models.Penimbangan
		rows.Scan(&p.ID, &p.IDAnak, &p.TglPenimbangan, &p.Umur, &p.BeratBadan, &p.TinggiBadan, &p.JenisImunisasi, &p.Keterangan, &p.Petugas)
		list = append(list, p)
	}

	data := map[string]interface{}{
		"Penimbangan": list,
		"IDAnak":      idAnak,
		"Role":        sessionData["Role"],
	}

	tmpl := template.Must(template.ParseFiles("views/penimbangan/data.html"))
	tmpl.Execute(w, data)
}

func CreatePenimbangan(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	idAnak := r.URL.Query().Get("id_anak")
	tmpl := template.Must(template.ParseFiles("views/penimbangan/create.html"))
	tmpl.Execute(w, map[string]string{"IDAnak": idAnak})
}

func StorePenimbangan(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idAnak := r.FormValue("id_anak")
		_, err := config.DB.Exec("INSERT INTO t_penimbangan (id_anak, tgl_penimbangan, umur, berat_badan, tinggi_badan, jenis_imunisasi, keterangan, petugas, posyandu) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			idAnak, r.FormValue("tgl_penimbangan"), r.FormValue("umur"), r.FormValue("berat_badan"),
			r.FormValue("tinggi_badan"), r.FormValue("jenis_imunisasi"), r.FormValue("keterangan"),
			r.FormValue("petugas"), r.FormValue("posyandu"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/data_penimbangan?id_anak="+idAnak, http.StatusSeeOther)
	}
}

func EditPenimbangan(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id_penimbangan")
	var p models.Penimbangan
	err := config.DB.QueryRow("SELECT id_penimbangan, id_anak, tgl_penimbangan, umur, berat_badan, tinggi_badan, jenis_imunisasi, keterangan, petugas, posyandu FROM t_penimbangan WHERE id_penimbangan = ?", id).Scan(
		&p.ID, &p.IDAnak, &p.TglPenimbangan, &p.Umur, &p.BeratBadan, &p.TinggiBadan, &p.JenisImunisasi, &p.Keterangan, &p.Petugas, &p.Posyandu)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("views/penimbangan/edit.html"))
	tmpl.Execute(w, p)
}

func UpdatePenimbangan(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.URL.Query().Get("id_penimbangan")
		idAnak := r.FormValue("id_anak")
		_, err := config.DB.Exec("UPDATE t_penimbangan SET tgl_penimbangan=?, umur=?, berat_badan=?, tinggi_badan=?, jenis_imunisasi=?, keterangan=?, petugas=?, posyandu=? WHERE id_penimbangan=?",
			r.FormValue("tgl_penimbangan"), r.FormValue("umur"), r.FormValue("berat_badan"), r.FormValue("tinggi_badan"),
			r.FormValue("jenis_imunisasi"), r.FormValue("keterangan"), r.FormValue("petugas"), r.FormValue("posyandu"), id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/data_penimbangan?id_anak="+idAnak, http.StatusSeeOther)
	}
}

func DeletePenimbangan(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id_penimbangan")
	idAnak := r.URL.Query().Get("id_anak")
	config.DB.Exec("DELETE FROM t_penimbangan WHERE id_penimbangan=?", id)
	http.Redirect(w, r, "/data_penimbangan?id_anak="+idAnak, http.StatusSeeOther)
}

func KMS(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	idAnak := r.URL.Query().Get("id_anak")

	rows, err := config.DB.Query("SELECT umur, berat_badan, tinggi_badan, jenis_imunisasi FROM t_penimbangan WHERE id_anak = ? ORDER BY umur ASC", idAnak)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []models.KMSData
	for rows.Next() {
		var k models.KMSData
		var imunisasi string
		rows.Scan(&k.Umur, &k.BeratBadan, &k.TinggiBadan, &imunisasi)

		k.Tooltip = fmt.Sprintf("Umur %d bln", k.Umur)
		if imunisasi != "" {
			k.Tooltip += fmt.Sprintf(" (Imunisasi: %s)", imunisasi)
		}
		data = append(data, k)
	}

	jsonData, _ := json.Marshal(data)

	normalRanges := []map[string]interface{}{
		{"umur": 0, "lower": 3.2, "upper": 3.3},
		{"umur": 1, "lower": 4.2, "upper": 4.5},
		{"umur": 2, "lower": 5.1, "upper": 5.6},
		{"umur": 3, "lower": 5.8, "upper": 6.4},
		{"umur": 4, "lower": 6.4, "upper": 7},
		{"umur": 5, "lower": 6.9, "upper": 7.5},
		{"umur": 6, "lower": 7.3, "upper": 7.9},
		{"umur": 7, "lower": 7.6, "upper": 8.3},
		{"umur": 8, "lower": 7.9, "upper": 8.6},
		{"umur": 9, "lower": 8.2, "upper": 8.9},
		{"umur": 10, "lower": 8.5, "upper": 9.2},
		{"umur": 11, "lower": 8.7, "upper": 9.4},
		{"umur": 12, "lower": 8.9, "upper": 9.6},
		{"umur": 24, "lower": 11.5, "upper": 12.2},
	}
	jsonRanges, _ := json.Marshal(normalRanges)

	tmpl := template.Must(template.ParseFiles("views/penimbangan/kms.html"))
	tmpl.Execute(w, map[string]interface{}{
		"JsonData":   string(jsonData),
		"JsonRanges": string(jsonRanges),
		"IDAnak":     idAnak,
	})
}
