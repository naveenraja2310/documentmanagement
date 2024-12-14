package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

type Result struct {
	Status  bool
	Data    interface{}
	Message error
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	files, err := os.ReadDir(path)
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: err})
		return
	}

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		})
	}

	json.NewEncoder(w).Encode(Result{Status: true, Data: fileInfos, Message: nil})
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("path is undefined")})
		return
	}

	http.ServeFile(w, r, path)
}

func renameFileOrFolder(w http.ResponseWriter, r *http.Request) {
	oldPath := r.URL.Query().Get("oldPath")
	newPath := r.URL.Query().Get("newPath")

	if oldPath == "" || newPath == "" {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("both path are required")})
		return
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("file not found")})
		return
	}

	json.NewEncoder(w).Encode(Result{Status: true, Data: nil, Message: nil})
}

func createTextFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fileName := r.URL.Query().Get("fileName")
	if path == "" {
		path = "."
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("fail to read body")})
		return
	}
	filepath := filepath.Join(path, fileName)

	file, err := os.Create(filepath)
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("failed to create file")})
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("failed to write to file")})
		return
	}

	json.NewEncoder(w).Encode(Result{Status: true, Data: nil, Message: nil})
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("path is undefined")})
		return
	}

	err := os.Remove(filepath.Clean(path))
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("failed to delete file")})
		return
	}

	json.NewEncoder(w).Encode(Result{Status: true, Data: nil, Message: nil})
}

func readTextFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("path is undefined")})
		return
	}

	content, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("failed to read file")})
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(Result{Status: true, Data: string(content), Message: nil})
}

func updateTextFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("parse error")})
		return
	}

	path := r.FormValue("path")
	if path == "" {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("path is undefined")})
		return
	}

	content := r.FormValue("content")
	if content == "" {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: errors.New("content is required")})
		return
	}

	err = ioutil.WriteFile(filepath.Clean(path), []byte(content), 0644)
	if err != nil {
		json.NewEncoder(w).Encode(Result{Status: false, Data: nil, Message: err})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Result{Status: true, Data: "", Message: nil})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/list", listFiles)
	mux.HandleFunc("/download", downloadFile)
	mux.HandleFunc("/rename", renameFileOrFolder)

	mux.HandleFunc("/createtextfile", createTextFile)
	mux.HandleFunc("/deletefile", deleteFile)
	mux.HandleFunc("/readfile", readTextFile)
	mux.HandleFunc("/updatefile", updateTextFile)

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/", fs)
	corsMux := corsMiddleware(mux)

	http.ListenAndServe(":8080", corsMux)

}
