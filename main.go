package main

import (
	"log"
	"net/http"
	"os"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
	w.Write([]byte("こんにちは！"))
}

func main() {
	port := os.Getenv("PORT") // 実行時に Heroku が指定するポート番号を取得
	if len(port) == 0 {
		port = "8080" // ローカルで実行するときのポート番号を指定
	}
	http.HandleFunc("/", handleHome)
	log.Printf("ポート %s で待ち受けを開始します...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("サーバーが異常終了しました: %v", err)
	}
}
