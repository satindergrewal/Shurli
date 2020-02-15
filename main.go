package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func idx(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Sub-Atomic Swaps Rock!\n"))
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", idx)
	log.Fatal(http.ListenAndServe(":8080", r))
}
