package main

import (
	"fmt"
	"orchestra"
	"net/http"
	"io"
)

func main() {
	fmt.Println("Welcome to Orchestra!")
	var o *orchestra.Orchestra = orchestra.NewOrchestra("127.0.0.1", 8000)
	o.HandleFunc("/", func (w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Add("Content-type", "text/html")
		for k, v := range p {
			io.WriteString(w, k + ": " + v + "<br>")
		}
	})
	o.HandleFunc("/test", func (w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Add("Content-type", "text/html")
		io.WriteString(w, "<input />")
	})
	o.HandleFunc("/test/:id/:edit", func (w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Add("Content-type", "text/html")
		for k, v := range p {
			io.WriteString(w, k + ": " + v + "<br>")
		}
	})
	o.ListenAndServe()
}

