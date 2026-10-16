package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var allowedExts = map[string]bool{"png": true, "jpg": true, "jpeg": true, "gif": true}

func allowedFile(fn string) bool {
	idx := strings.LastIndex(fn, ".")
	if idx < 0 {
		return false
	}
	ext := strings.ToLower(fn[idx+1:])
	return allowedExts[ext]
}

func getFileVersion(filename, folder string) string {
	path := filepath.Join(folder, filename)
	if info, err := os.Stat(path); err == nil {
		return fmt.Sprintf("%d", info.ModTime().Unix())
	}
	return fmt.Sprintf("%d", time.Now().Unix())
}

type SkinUser struct {
	ID        int
	Username  string
	SkinPath  string
	CloakPath *string
}

func getCurrentUser(r *http.Request) *SkinUser {
	cookie, err := r.Cookie("skin_user")
	if err != nil {
		return nil
	}
	username := cookie.Value
	var u SkinUser
	var cloak *string
	err = db.QueryRow("SELECT id, username, IFNULL(skin_path,'default.png'), cloak_path FROM "+dbTable+" WHERE "+colUser+"=?", username).
		Scan(&u.ID, &u.Username, &u.SkinPath, &cloak)
	if err != nil {
		return nil
	}
	u.CloakPath = cloak
	return &u
}

func skinIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.FormValue("login") != "" {
			u := r.FormValue("username")
			p := r.FormValue("password")
			var storedPass string
			err := db.QueryRow("SELECT "+colPass+" FROM "+dbTable+" WHERE "+colUser+"=?", u).Scan(&storedPass)
			if err == nil && storedPass == md5Hex(p) {
				http.SetCookie(w, &http.Cookie{Name: "skin_user", Value: u, Path: "/"})
				http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
				return
			}
			render(w, "skin_index", map[string]interface{}{
				"Title": "Skin System",
				"Flash": "Incorrect login or password",
				"User":  nil,
			})
			return
		}

		user := getCurrentUser(r)

		if r.FormValue("upload_skin") != "" && user != nil {
			f, header, err := r.FormFile("skin_file")
			if err != nil {
				http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
				return
			}
			defer f.Close()
			if !allowedFile(header.Filename) {
				http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
				return
			}
			ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, ".")+1:])
			fn := user.Username + "." + ext
			os.MkdirAll("skins", 0755)
			dst, _ := os.Create(filepath.Join("skins", fn))
			io.Copy(dst, f)
			dst.Close()
			db.Exec("UPDATE "+dbTable+" SET skin_path=? WHERE id=?", fn, user.ID)
			http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
			return
		}

		if r.FormValue("upload_cloak") != "" && user != nil {
			f, header, err := r.FormFile("cloak_file")
			if err != nil {
				http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
				return
			}
			defer f.Close()
			if !allowedFile(header.Filename) {
				http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
				return
			}
			ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, ".")+1:])
			fn := user.Username + "." + ext
			os.MkdirAll("cloak", 0755)
			dst, _ := os.Create(filepath.Join("cloak", fn))
			io.Copy(dst, f)
			dst.Close()
			db.Exec("UPDATE "+dbTable+" SET cloak_path=? WHERE id=?", fn, user.ID)
			http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
			return
		}
	}

	if r.URL.Query().Get("logout") != "" {
		http.SetCookie(w, &http.Cookie{Name: "skin_user", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/SkinSystem/", http.StatusSeeOther)
		return
	}

	user := getCurrentUser(r)
	skinVersion := fmt.Sprintf("%d", time.Now().Unix())
	cloakVersion := fmt.Sprintf("%d", time.Now().Unix())

	if user != nil {
		if user.SkinPath != "" {
			skinVersion = getFileVersion(user.SkinPath, "skins")
		}
		if user.CloakPath != nil {
			cloakVersion = getFileVersion(*user.CloakPath, "cloak")
		}
	}

	render(w, "skin_index", map[string]interface{}{
		"Title":        "Skin System",
		"User":         user,
		"SkinVersion":  skinVersion,
		"CloakVersion": cloakVersion,
		"Flash":        "",
	})
}

func serveSkin(w http.ResponseWriter, r *http.Request) {
	fn := filepath.Base(r.URL.Path)
	http.ServeFile(w, r, filepath.Join("skins", fn))
}

func serveCloak(w http.ResponseWriter, r *http.Request) {
	fn := filepath.Base(r.URL.Path)
	http.ServeFile(w, r, filepath.Join("cloak", fn))
}