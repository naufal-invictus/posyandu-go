package main

import (
	"log"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/controllers"
)

func main() {
	config.ConnectDB()

	// Static
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("public/img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r) 
			return
		}
		http.ServeFile(w, r, "views/index.html")
	})

	// 2. AUTH
	http.HandleFunc("/login", controllers.LoginPage)
	http.HandleFunc("/login_process", controllers.LoginProcess)
	http.HandleFunc("/logout", controllers.Logout)

	// 3. MANAJEMEN ANAK (Pastikan semua route terdaftar)
	http.HandleFunc("/data_anak", controllers.DataAnak)
	http.HandleFunc("/create_anak", controllers.CreateAnak)
	http.HandleFunc("/store_anak", controllers.StoreAnak)
	http.HandleFunc("/edit_anak", controllers.EditAnak)     
	http.HandleFunc("/update_anak", controllers.UpdateAnak) 
	http.HandleFunc("/delete_anak", controllers.DeleteAnak) 

	// 4. MANAJEMEN PENIMBANGAN
	http.HandleFunc("/data_penimbangan", controllers.DataPenimbangan)
	http.HandleFunc("/create_penimbangan", controllers.CreatePenimbangan)
	http.HandleFunc("/store_penimbangan", controllers.StorePenimbangan)
	http.HandleFunc("/edit_penimbangan", controllers.EditPenimbangan)
	http.HandleFunc("/update_penimbangan", controllers.UpdatePenimbangan)
	http.HandleFunc("/delete_penimbangan", controllers.DeletePenimbangan)
	http.HandleFunc("/kms", controllers.KMS)

	// 5. MANAJEMEN USER & ORANG TUA (Fitur Admin)
	http.HandleFunc("/create_orangtua", controllers.CreateOrangTua) 
	http.HandleFunc("/store_orangtua", controllers.StoreOrangTua)   

	// 6. LAPORAN & API
	http.HandleFunc("/laporan", controllers.HalamanLaporan)
	http.HandleFunc("/api/jadwal", controllers.ApiJadwal)
	http.HandleFunc("/api/stok_vaksin", controllers.ApiStokVaksin)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}