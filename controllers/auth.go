package controllers

import (
	"database/sql"
	"net/http"
	"sipograf-go/config"
	"sipograf-go/models"
	"time"
)

// Global Key untuk session (Sederhana)
var sessionStore = make(map[string]map[string]interface{})

func LoginPage(w http.ResponseWriter, r *http.Request) {
	// Jika sudah login, lempar ke dashboard sesuai role
	if IsAuthenticated(r) {
		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "views/login.html")
}

func LoginProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		err := config.DB.QueryRow("SELECT id_user, username, password, role FROM users WHERE username = ? AND password = ?", username, password).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?err=failed", http.StatusSeeOther)
			return
		}

		sessionID := "session_" + username + "_" + time.Now().Format("20060102150405")
		
		// Simpan data session di memory server
		sessionStore[sessionID] = map[string]interface{}{
			"UserID":   user.ID,
			"Username": user.Username,
			"Role":     user.Role,
		}

		// SET COOKIE DENGAN PATH
		http.SetCookie(w, &http.Cookie{
			Name:     "sipograf_session",
			Value:    sessionID,
			Path:     "/",            
			HttpOnly: true,
			MaxAge:   3600 * 24,    
		})

		http.Redirect(w, r, "/data_anak", http.StatusSeeOther)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("sipograf_session")
	if err == nil {
		delete(sessionStore, c.Value)
	}
	// Hapus Cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "sipograf_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func IsAuthenticated(r *http.Request) bool {
	c, err := r.Cookie("sipograf_session")
	if err != nil {
		return false
	}
	_, ok := sessionStore[c.Value]
	return ok
}

func GetSessionDetails(r *http.Request) map[string]interface{} {
	c, err := r.Cookie("sipograf_session")
	if err != nil {
		return nil
	}
	return sessionStore[c.Value]
}