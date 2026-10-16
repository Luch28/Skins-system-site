package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "index", map[string]interface{}{"Title": "Home"})
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "register", map[string]interface{}{"Title": "Registration"})
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	login := strings.TrimSpace(r.FormValue("login"))
	pass := strings.TrimSpace(r.FormValue("passwd"))
	repass := strings.TrimSpace(r.FormValue("repasswd"))

	msg := ""
	ok := false

	switch {
	case login == "" || pass == "" || repass == "":
		msg = "Not all fields are filled in."
	case !validChars(login):
		msg = "Login entered incorrectly."
	case !validChars(pass):
		msg = "Password entered incorrectly."
	case !validChars(repass):
		msg = "Password confirmation entered incorrectly."
	default:
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM "+dbTable+" WHERE "+colUser+"=?", login).Scan(&exists)
		if err != nil {
			msg = "Database query failed."
			break
		}
		if exists > 0 {
			msg = fmt.Sprintf("Account <b>%s</b> already exists.", login)
			break
		}
		if len(login) < 4 || len(login) > 8 {
			msg = "Login must contain at least 4 characters and no more than 8."
			break
		}
		if len(pass) < 4 || len(pass) > 15 {
			msg = "Password must contain at least 4 characters and no more than 15."
			break
		}
		if pass != repass {
			msg = "Passwords do not match."
			break
		}
		_, err = db.Exec("INSERT INTO "+dbTable+" ("+colUser+","+colPass+") VALUES(?,?)", login, md5Hex(pass))
		if err != nil {
			msg = "Database query failed."
			break
		}
		msg = fmt.Sprintf("Account <b>%s</b> successfully registered.", login)
		ok = true
	}

	render(w, "reg", map[string]interface{}{
		"Title": "Registration",
		"Msg":   template.HTML(msg),
		"OK":    ok,
	})
}