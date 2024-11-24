package main

import (
	"encoding/json"
	"net/http"
)

type responseMessage struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// GETリクエスト以外は受け付けない
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// クエリパラメータからnameを取得
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// レスポンスヘッダーにContent-Typeを設定
	w.Header().Set("Content-Type", "application/json")

	// レスポンスメッセージを作成
	message := responseMessage{
		Message: "Hello, " + name + "-san!",
	}

	// JSONにエンコード
	bytes, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// レスポンスを書き込む
	w.Write(bytes)
}

func main() {
	http.HandleFunc("/hello", handler)
	http.ListenAndServe(":8080", nil)
}
