package main

import (
	"fmt"
	"net/http"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("user")
	postPass := r.FormValue("password")
	ver := r.FormValue("version")

	launcherVer := getGameInfo("launcher")
	if launcherVer != ver {
		fmt.Fprint(w, "Old version")
		return
	}

	var realPass string
	err := db.QueryRow("SELECT "+colPass+" FROM "+dbTable+" WHERE "+colUser+"=?", login).Scan(&realPass)
	if err != nil {
		fmt.Fprint(w, "Bad login")
		return
	}

	if realPass != md5Hex(postPass) {
		fmt.Fprint(w, "Bad login")
		return
	}

	sessid := generateSessionID()
	gameBuild := getGameInfo("build")

	_, err = db.Exec("UPDATE "+dbTable+" SET "+colSesID+"=? WHERE "+colUser+"=?", sessid, login)
	if err != nil {
		fmt.Fprint(w, "Bad login")
		return
	}

	dlticket := md5Hex(login)
	fmt.Fprintf(w, "%s:%s:%s:%s:", gameBuild, dlticket, login, sessid)
}

func joinServerHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := r.URL.Query().Get("sessionId")
	user := r.URL.Query().Get("user")
	serverid := r.URL.Query().Get("serverId")

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM "+dbTable+" WHERE "+colSesID+"=? AND "+colUser+"=? AND "+colServer+"=?", sessionid, user, serverid).Scan(&count)
	if err == nil && count == 1 {
		fmt.Fprint(w, "OK")
		return
	}

	res, err := db.Exec("UPDATE "+dbTable+" SET "+colServer+"=? WHERE "+colSesID+"=? AND "+colUser+"=?", serverid, sessionid, user)
	if err != nil {
		fmt.Fprint(w, "Bad login")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 1 {
		fmt.Fprint(w, "OK")
	} else {
		fmt.Fprint(w, "Bad login")
	}
}

func checkServerHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	serverid := r.URL.Query().Get("serverId")

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM "+dbTable+" WHERE "+colUser+"=? AND "+colServer+"=?", user, serverid).Scan(&count)
	if err != nil || count != 1 {
		fmt.Fprint(w, "NO")
		return
	}
	fmt.Fprint(w, "YES")
}