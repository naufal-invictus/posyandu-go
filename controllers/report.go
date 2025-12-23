package controllers

import (
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
	"text/template"
)

func HalamanLaporan(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var stats models.LaporanStats

	config.DB.QueryRow("SELECT COUNT(*) FROM anak").Scan(&stats.TotalAnak)
	config.DB.QueryRow("SELECT COUNT(*) FROM penimbangan").Scan(&stats.TotalPenimbangan)
	config.DB.QueryRow("SELECT COALESCE(SUM(jumlah), 0) FROM stok_vaksin").Scan(&stats.TotalVaksin)

	sessionData := GetSessionDetails(r)
	data := map[string]interface{}{
		"Stats": stats,
		"Role":  sessionData["Role"],
	}

	tmpl := template.Must(template.ParseFiles("views/laporan.html"))
	tmpl.Execute(w, data)
}