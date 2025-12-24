package models

// Struct User untuk Login
type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

// Struct Anak (Relasi ke OrangTua)
type Anak struct {
	ID           int
	IDOrangtua   int
	NamaAnak     string
	TempatLahir  string
	TanggalLahir string
	JenisKelamin string
	NamaIbu      string 
	NamaAyah     string
	Alamat       string
}

// Struct Penimbangan
type Penimbangan struct {
	ID             int
	IDAnak         int
	TanggalTimbang string
	UmurBulan      int
	BeratBadan     float64
	TinggiBadan    float64
	LingkarKepala  float64
	Petugas        string
}

// Struct untuk Dropdown/Data Orang Tua
type Parent struct {
    ID int
    Username string
    NamaIbu string
    Alamat string
}

// Struct untuk Grafik KMS (JSON)
type KMSData struct {
	Umur        int     `json:"umur"`
	BeratBadan  float64 `json:"berat_badan"`
	TinggiBadan float64 `json:"tinggi_badan"`
	Tooltip     string  `json:"tooltip"`
}


// Struct Jadwal Kegiatan (Dipakai di API & Report)
type Jadwal struct {
	ID           int    `json:"id"`
	NamaKegiatan string `json:"nama_kegiatan"`
	Tanggal      string `json:"tanggal"`
	Lokasi       string `json:"lokasi"`
	Keterangan   string `json:"keterangan"`
}

// Struct Stok Vaksin (Dipakai di API & Report)
type StokVaksin struct {
	ID          int    `json:"id"`
	NamaVaksin  string `json:"nama_vaksin"`
	Jumlah      int    `json:"jumlah"`
	TanggalUpd  string `json:"tanggal_update"`
}

// Struct Statistik Dashboard Laporan
type LaporanStats struct {
	TotalAnak        int
	TotalPenimbangan int
	TotalVaksin      int
}