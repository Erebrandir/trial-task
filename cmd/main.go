package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"trial_task/internal/download"
	"trial_task/internal/upload"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/upload", upload.UploadFile)
	r.HandleFunc("/download", download.DownloadFile)

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
