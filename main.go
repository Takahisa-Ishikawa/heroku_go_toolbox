package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	shukujitsu "github.com/Takahisa-Ishikawa/shukujistu-go"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("template/index.gohtml", "template/_menu.gohtml"))
	if err := t.Execute(w, struct {
		UserName string
		Time     time.Time
	}{
		"ゲスト",
		time.Now(),
	}); err != nil {
		log.Printf("テンプレート %s の実行に失敗!: %v", t.Name(), err)
		http.Error(w, "内部エラーです", http.StatusInternalServerError)
	}
}

func handleSecret(w http.ResponseWriter, r *http.Request) {
	user, password, _ := r.BasicAuth()
	if user != "user" || password != "password" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restrected"`)
		http.Error(w, "認証に失敗しました", http.StatusUnauthorized)
		return
	}
	log.Printf("%s %s", r.Method, r.RequestURI)
	w.Write([]byte("秘密のページです"))
}

func handleHoliday(w http.ResponseWriter, r *http.Request) {
	entries, err := shukujitsu.AllEntries()
	if err != nil {
		log.Printf("祝日の取得に失敗!: %v", err)
		http.Error(w, "内部エラーです", http.StatusInternalServerError)
	}

	t := template.Must(template.ParseFiles("template/holiday.gohtml", "template/_menu.gohtml"))
	if err := t.Execute(w, struct {
		Time     time.Time
		Holidays []shukujitsu.Entry
	}{
		time.Now(),
		entries,
	}); err != nil {
		log.Printf("テンプレート %s の実行に失敗!: %v", t.Name(), err)
		http.Error(w, "内部エラーです", http.StatusInternalServerError)
	}
}

func main() {
	port := os.Getenv("PORT") // 実行時に Heroku が指定するポート番号を取得
	if len(port) == 0 {
		port = "8080" // ローカルで実行するときのポート番号を指定
	}
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/secret", handleSecret)
	http.HandleFunc("/holiday", handleHoliday)
	log.Printf("ポート %s で待ち受けを開始します...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("サーバーが異常終了しました: %v", err)
	}
}
