package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form) //prints server side information
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //sends data to client side

}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method of request:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, token)
	} else {
		//max memory on server
		r.ParseMultipartForm(32 << 20)
		// get file handle
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		//copy to your file system
		io.Copy(f, file)
	}

}

func main() {
	http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/upload", upload)
	// nil --> DefaultServeMux, router variable which can cal handler functions fo specified URLs
	err := http.ListenAndServe(":9090", nil) //set listening port and initialize a server object
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
