package controllers

import (
	"database/sql"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
	"text/template"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("secret-key"))

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/login.html"))
	tmpl.Execute(w, nil)
}

func LoginProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user models.User
	err := config.DB.QueryRow("SELECT id_user, username, role FROM users WHERE username = ? AND password = ?", username, password).Scan(&user.ID, &user.Username, &user.Role)

	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	session.Values["id_user"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func IsAuthenticated(r *http.Request) bool {
	session, _ := Store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
}

func GetSessionDetails(r *http.Request) map[string]interface{} {
	session, _ := Store.Get(r, "session-name")
	return map[string]interface{}{
		"Role":     session.Values["role"],
		"Username": session.Values["username"],
		"UserID":   session.Values["id_user"],
	}
}
