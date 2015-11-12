package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

func upload(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "GET" {
		t, err := template.ParseFiles("template/upload.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, time.Now().Unix())
	} else {
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("uploadfile")
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		newFile, err := os.Create("upload/" + header.Filename)
		defer newFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		io.Copy(newFile, file)
		http.Redirect(w, r, "/upload", 302)
	}
}

func main() {
	addr := ":8090"
	fmt.Println("FileUpload Listening on ", addr)
	http.HandleFunc("/upload", upload)
	http.Handle("/", http.FileServer(http.Dir("upload")))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
