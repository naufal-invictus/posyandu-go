package models

import (
	"database/sql"
)

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

type Anak struct {
	ID           int
	IDOrangtua   sql.NullInt64
	NamaAnak     string
	NamaIbu      string
	TempatLahir  string
	TanggalLahir string
	JenisKelamin string
	Alamat       string
}

type Penimbangan struct {
	ID             int
	IDAnak         int
	TglPenimbangan string
	Umur           int
	BeratBadan     float64
	TinggiBadan    float64
	JenisImunisasi string
	Keterangan     string
	Petugas        string
	Posyandu       string
}

type KMSData struct {
	Umur        int     `json:"umur"`
	BeratBadan  float64 `json:"berat_badan"`
	TinggiBadan float64 `json:"tinggi_badan"`
	Tooltip     string  `json:"tooltip"`
}
