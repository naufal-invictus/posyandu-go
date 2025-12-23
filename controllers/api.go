package controllers

import (
	"encoding/json"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
)

func ApiJadwal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	rows, err := config.DB.Query("SELECT id_jadwal, nama_kegiatan, tanggal, lokasi, keterangan FROM jadwal_kegiatan")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jadwal []models.Jadwal
	for rows.Next() {
		var j models.Jadwal
		rows.Scan(&j.ID, &j.NamaKegiatan, &j.Tanggal, &j.Lokasi, &j.Keterangan)
		jadwal = append(jadwal, j)
	}

	json.NewEncoder(w).Encode(jadwal)
}

func ApiStokVaksin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT s.id_stok, j.nama_vaksin, s.jumlah, s.tanggal_update 
			  FROM stok_vaksin s 
			  JOIN jenis_vaksin j ON s.id_vaksin = j.id_vaksin`
	
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stok []models.StokVaksin
	for rows.Next() {
		var s models.StokVaksin
		rows.Scan(&s.ID, &s.NamaVaksin, &s.Jumlah, &s.TanggalUpd)
		stok = append(stok, s)
	}

	json.NewEncoder(w).Encode(stok)
}