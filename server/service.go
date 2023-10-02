package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const DATA string = "path/to/your/home/dir"

func setUpLogger() {
	logFile, err := os.OpenFile("logger.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(logFile)
}

func setUpServer(mux *http.ServeMux) {
	// timeout := time.Duration(5) * time.Second

	// transport := &http.Transport{
	// 	ResponseHeaderTimeout: timeout,
	// 	Dial: func(network, addr string) (net.Conn, error) {
	// 		return net.DialTimeout(network, addr, timeout)
	// 	},
	// 	DisableKeepAlives: true,
	// }

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server closed\n")
	} else if err != nil {
		log.Fatalf("error listening for server: %s\n", err)
	}

	log.Printf("Server started on http://localhost:%d", port)
}

// listed folder files
func list(w http.ResponseWriter, r *http.Request) {

	hasDir := r.URL.Query().Has("dir")
	var mainData *responseJSON
	var filesList []fileSystem
	tm := time.Now()
	falseSucces := &responseJSON{
		Success: false,
		Date:    tm,
		Files:   []fileSystem{},
	}
	// if dir exist in params
	if hasDir {

		dir := r.URL.Query().Get("dir")
		// if dir is not empty
		if dir != "" {

			folder, err := os.Open(filepath.Join(DATA, dir))

			if err != nil {
				log.Panicf("could not open dir: %s\n", err)
			}

			files, err := folder.Readdir(0)

			if err != nil {
				log.Panicf("could not read dir: %s\n", err)
			}

			for i, file := range files {
				currentFile := fileSystem{
					Id:    i,
					Name:  file.Name(),
					IsDir: file.IsDir(),
				}

				filesList = append(filesList, currentFile)
			}

			mainData = &responseJSON{
				Success: true,
				Date:    tm,
				Files:   filesList,
			}
			defer folder.Close()
		} else {
			mainData = falseSucces
		}

	} else {
		mainData = falseSucces
	}
	// make json from struct
	response, err := json.Marshal(mainData)

	if err != nil {
		log.Panicf("could not marshal json: %s\n", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(response))
}

// download file
// TODO: download folders as zip
func download(w http.ResponseWriter, r *http.Request) {
	hasIsDir := r.URL.Query().Has("isDir")
	hasName := r.URL.Query().Has("name")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if hasName && hasIsDir {
		name := r.URL.Query().Get("name")

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		file, err := os.Open(filepath.Join(DATA, name))
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()

		io.Copy(w, file)
	}
}

// delete file
// TODO: delete folders
func delete(w http.ResponseWriter, r *http.Request) {
	hasIsDir := r.URL.Query().Has("isDir")
	hasName := r.URL.Query().Has("name")
	tm := time.Now()
	falseSucces := &responseJSON{
		Success: false,
		Date:    tm,
		Files:   []fileSystem{},
	}
	var mainData *responseJSON
	if hasName && hasIsDir {
		name := r.URL.Query().Get("name")

		isDir, err := strconv.ParseBool(r.URL.Query().Get("isDir"))
		if err != nil {
			log.Panic(err)
		}

		if name != "" {
			err := os.Remove(filepath.Join(DATA, name))
			if err != nil {
				log.Panic(err)
			}
			mainData = &responseJSON{
				Success: true,
				Date:    tm,
				Files: []fileSystem{{
					Name:  name,
					IsDir: isDir,
				}},
			}
		} else {
			mainData = falseSucces
		}
	} else {
		mainData = falseSucces
	}
	response, err := json.Marshal(mainData)
	if err != nil {
		log.Panicf("could not marshal json: %s\n", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(response))
}

// rename file
// TODO: rename folder
func renameFile(w http.ResponseWriter, r *http.Request) {

	hasName := r.URL.Query().Has("name")
	hasNewName := r.URL.Query().Has("newName")
	var mainData *responseJSON
	var filesList []fileSystem
	tm := time.Now()
	falseSucces := &responseJSON{
		Success: false,
		Date:    tm,
		Files:   []fileSystem{},
	}
	// if exist params
	if hasName && hasNewName {

		name := r.URL.Query().Get("name")
		newName := r.URL.Query().Get("newName")

		// if dir is not empty
		if newName != "" && name != "" {
			oldPath := filepath.Join(DATA, name)
			dir, _ := filepath.Split(oldPath)

			err := os.Rename(oldPath, filepath.Join(dir, newName))
			if err != nil {
				log.Panicf("could not rename : %s\n", err)
				mainData = falseSucces
			}

			filesList = append(filesList, fileSystem{
				Name:  newName,
				IsDir: false,
			})

			mainData = &responseJSON{
				Success: true,
				Date:    tm,
				Files:   filesList,
			}

		} else {
			mainData = falseSucces
		}

	} else {
		mainData = falseSucces
	}
	// make json from struct
	response, err := json.Marshal(mainData)

	if err != nil {
		log.Panicf("could not marshal json: %s\n", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(response))
}

// upload files
func upload(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	dirForUpload := DATA
	hasDir := r.URL.Query().Has("dir")

	if hasDir {
		dir := r.URL.Query().Get("dir")

		if dir != "" {
			dirForUpload = filepath.Join(DATA, dir)
		}
	}
	err := r.ParseMultipartForm(200000)
	if err != nil {
		log.Panicf("could not ParseMultipartForm : %s\n", err)
		return
	}

	files := r.MultipartForm.File

	for _, item := range files {

		file := item[0]

		openedFile, err := file.Open()
		if err != nil {
			log.Panicf("could not open file : %s\n", err)
			return
		}
		defer openedFile.Close()

		out, err := os.Create(filepath.Join(dirForUpload, file.Filename))
		if err != nil {
			log.Panicf("Unable to create the file for writing. Check your write access privilege : %s\n", err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, openedFile)
		if err != nil {
			log.Panicf("could not copy file : %s\n", err)
			return
		}
	}
}

// copy file to
// TODO: copy folders
func copy(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	hasDir := r.URL.Query().Has("dir")
	hasNewDir := r.URL.Query().Has("newDir")
	hasFileName := r.URL.Query().Has("name")

	if hasDir && hasNewDir && hasFileName {

		dir := r.URL.Query().Get("dir")
		newDir := r.URL.Query().Get("newDir")
		fileName := r.URL.Query().Get("name")

		if dir != "" && newDir != "" && fileName != "" {
			//косяк открывается только директория надо файл
			file, err := os.Open(filepath.Join(DATA, dir))
			if err != nil {
				log.Panic(err)
			}
			defer file.Close()
			//здесь только директория неверно
			out, err := os.Create(filepath.Join(DATA, newDir))
			if err != nil {
				log.Panicf("Unable to create the file for writing. Check your write access privilege : %s\n", err)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				log.Panicf("could not copy file : %s\n", err)
				return
			}

		} else {
			log.Panicf("Empty params : %v, %v, %v\n", dir, newDir, fileName)
		}

	} else {
		log.Panicf("Missing params : %v, %v, %v\n", hasDir, hasNewDir, hasFileName)
	}
}
