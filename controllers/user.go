package controllers

import (
	"net/http"
	"sipograf-go/config"
	"text/template"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionData := GetSessionDetails(r)
	if sessionData["Role"] != "admin" {
		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("views/user/create.html"))
	tmpl.Execute(w, nil)
}

func StoreUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/create_user", http.StatusSeeOther)
		return
	}

	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	_, err := config.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", username, password, role)

	if err != nil {
		http.Error(w, "Gagal membuat user (Username mungkin sudah ada)", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
}
