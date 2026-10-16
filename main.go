package main

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

func main() {
	initDB()

	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/register", registerPageHandler)
	http.HandleFunc("/joinserver", joinServerHandler)
	http.HandleFunc("/checkserver", checkServerHandler)

	http.HandleFunc("/SkinSystem/", skinIndexHandler)
	http.HandleFunc("/skins/", serveSkin)
	http.HandleFunc("/cloak/", serveCloak)
	http.HandleFunc("/SkinSystem/skins/", serveSkin)
	http.HandleFunc("/SkinSystem/cloak/", serveCloak)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func render(w http.ResponseWriter, name string, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Println("Template error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hashPassword(password string) string {
	var h uint32 = 0
	for _, ch := range password {
		h = (31*h + uint32(ch)) & 0xFFFFFFFF
	}
	return hex.EncodeToString([]byte{
		byte(h >> 24), byte(h >> 16), byte(h >> 8), byte(h),
	})
}

func validChars(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func generateSessionID() string {
	a := rand.Intn(1147483647) + 1000000000
	b := rand.Intn(1147483647) + 1000000000
	c := rand.Intn(10)
	return formatInt(a) + formatInt(b) + formatInt(c)
}

func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func stripPrefix(path, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}