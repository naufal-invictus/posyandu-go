package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
	"text/template"
)

func disableCacheW(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func DataPenimbangan(w http.ResponseWriter, r *http.Request) {
	disableCacheW(w)
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idAnak := r.URL.Query().Get("id_anak")
	sessionData := GetSessionDetails(r)

	query := `SELECT id_penimbangan, id_anak, DATE_FORMAT(tanggal_timbang, '%Y-%m-%d'), 
			  umur_bulan, berat_badan, tinggi_badan, lingkar_kepala, petugas_pemeriksa 
			  FROM penimbangan WHERE id_anak = ? ORDER BY tanggal_timbang DESC`

	rows, err := config.DB.Query(query, idAnak)
	if err != nil {
		http.Error(w, "Error Database: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var list []models.Penimbangan
	for rows.Next() {
		var p models.Penimbangan
		// Scan 
		rows.Scan(&p.ID, &p.IDAnak, &p.TanggalTimbang, &p.UmurBulan, &p.BeratBadan, &p.TinggiBadan, &p.LingkarKepala, &p.Petugas)
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
	disableCacheW(w)
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
		
		tb := r.FormValue("tinggi_badan")
		lk := r.FormValue("lingkar_kepala")
		if tb == "" { tb = "0" }
		if lk == "" { lk = "0" }

		_, err := config.DB.Exec(`INSERT INTO penimbangan 
			(id_anak, tanggal_timbang, umur_bulan, berat_badan, tinggi_badan, lingkar_kepala, petugas_pemeriksa) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			idAnak, r.FormValue("tgl_penimbangan"), r.FormValue("umur"), r.FormValue("berat_badan"),
			tb, lk, r.FormValue("petugas"))

		if err != nil {
			http.Error(w, "Gagal Simpan: "+err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/data_penimbangan?id_anak="+idAnak, http.StatusSeeOther)
	}
}

func EditPenimbangan(w http.ResponseWriter, r *http.Request) {
	disableCacheW(w)
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id_penimbangan")
	var p models.Penimbangan
	
	
	query := `SELECT id_penimbangan, id_anak, DATE_FORMAT(tanggal_timbang, '%Y-%m-%d'), 
			  umur_bulan, berat_badan, tinggi_badan, lingkar_kepala, petugas_pemeriksa 
			  FROM penimbangan WHERE id_penimbangan = ?`

	err := config.DB.QueryRow(query, id).Scan(
		&p.ID, &p.IDAnak, &p.TanggalTimbang, &p.UmurBulan, &p.BeratBadan, &p.TinggiBadan, &p.LingkarKepala, &p.Petugas)

	if err != nil {
		http.Error(w, "Data tidak ditemukan: "+err.Error(), 404)
		return
	}

	tmpl := template.Must(template.ParseFiles("views/penimbangan/edit.html"))
	tmpl.Execute(w, p)
}

func UpdatePenimbangan(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.URL.Query().Get("id_penimbangan")
		idAnak := r.FormValue("id_anak")

		_, err := config.DB.Exec(`UPDATE penimbangan SET 
			tanggal_timbang=?, umur_bulan=?, berat_badan=?, tinggi_badan=?, lingkar_kepala=?, petugas_pemeriksa=? 
			WHERE id_penimbangan=?`,
			r.FormValue("tgl_penimbangan"), r.FormValue("umur"), r.FormValue("berat_badan"), 
			r.FormValue("tinggi_badan"), r.FormValue("lingkar_kepala"), r.FormValue("petugas"), id)

		if err != nil {
			http.Error(w, "Gagal Update: "+err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/data_penimbangan?id_anak="+idAnak, http.StatusSeeOther)
	}
}

func DeletePenimbangan(w http.ResponseWriter, r *http.Request) {
	disableCacheW(w)
	id := r.URL.Query().Get("id_penimbangan")
	idAnak := r.URL.Query().Get("id_anak")
	
	config.DB.Exec("DELETE FROM penimbangan WHERE id_penimbangan=?", id)
	http.Redirect(w, r, "/data_penimbangan?id_anak="+idAnak, http.StatusSeeOther)
}

func KMS(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	idAnak := r.URL.Query().Get("id_anak")

	rows, err := config.DB.Query("SELECT umur_bulan, berat_badan, tinggi_badan FROM penimbangan WHERE id_anak = ? ORDER BY umur_bulan ASC", idAnak)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var data []models.KMSData
	for rows.Next() {
		var k models.KMSData
		rows.Scan(&k.Umur, &k.BeratBadan, &k.TinggiBadan)
		k.Tooltip = fmt.Sprintf("Umur %d bln", k.Umur)
		data = append(data, k)
	}

	jsonData, _ := json.Marshal(data)
	
	normalRanges := []map[string]interface{}{
		{"umur": 0, "lower": 2.5, "upper": 3.9},
		{"umur": 12, "lower": 7.7, "upper": 10.8},
		{"umur": 24, "lower": 10.0, "upper": 13.0},
	}
	jsonRanges, _ := json.Marshal(normalRanges)

	tmpl := template.Must(template.ParseFiles("views/penimbangan/kms.html"))
	tmpl.Execute(w, map[string]interface{}{
		"JsonData":   string(jsonData),
		"JsonRanges": string(jsonRanges),
		"IDAnak":     idAnak,
	})
}