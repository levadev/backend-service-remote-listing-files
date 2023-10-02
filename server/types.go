package main

import "time"

type fileSystem struct {
	Name  string `json:"name"`
	Id    int    `json:"id"`
	IsDir bool   `json:"isDir"`
}

type responseJSON struct {
	Success bool `json:"success"`
	// MainDir string       `json:"mainDir"`
	Date  time.Time    `json:"date"`
	Files []fileSystem `json:"files"`
}
