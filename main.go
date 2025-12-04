package main

import (
	"log"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/controllers"
)

func main() {
	config.ConnectDB()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("public/img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/index.html")
	})

	http.HandleFunc("/login", controllers.LoginPage)
	http.HandleFunc("/login_process", controllers.LoginProcess)
	http.HandleFunc("/logout", controllers.Logout)

	http.HandleFunc("/data_anak", controllers.DataAnak)
	http.HandleFunc("/create_anak", controllers.CreateAnak)
	http.HandleFunc("/store_anak", controllers.StoreAnak)
	http.HandleFunc("/edit_anak", controllers.EditAnak)
	http.HandleFunc("/update_anak", controllers.UpdateAnak)
	http.HandleFunc("/delete_anak", controllers.DeleteAnak)

	http.HandleFunc("/data_penimbangan", controllers.DataPenimbangan)
	http.HandleFunc("/create_penimbangan", controllers.CreatePenimbangan)
	http.HandleFunc("/store_penimbangan", controllers.StorePenimbangan)
	http.HandleFunc("/edit_penimbangan", controllers.EditPenimbangan)
	http.HandleFunc("/update_penimbangan", controllers.UpdatePenimbangan)
	http.HandleFunc("/delete_penimbangan", controllers.DeletePenimbangan)
	http.HandleFunc("/kms", controllers.KMS)
http.HandleFunc("/create_user", controllers.CreateUser)
	http.HandleFunc("/store_user", controllers.StoreUser)
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
