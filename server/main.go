package main

import (
	"context"
	"net/http"
)

const port int = 3333

func main() {
	setUpLogger()

	mux := http.NewServeMux()
	mux.HandleFunc("/list", list)
	mux.HandleFunc("/delete", delete)
	mux.HandleFunc("/rename", renameFile)
	// mux.HandleFunc("/move", move) endpoint in work
	mux.HandleFunc("/copy", copy)
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/download", download)

	ctx := context.Background()
	setUpServer(mux)
	<-ctx.Done()

}
