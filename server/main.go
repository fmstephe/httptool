package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/", handle)

	http.ListenAndServe(":3333", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error: %s", err)
		return
	}
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("Error: %s", err)
		return
	}

	fmt.Printf("METHOD: %s\n\n", r.Method)
	fmt.Printf("URL: %s\n\n", r.URL)
	fmt.Printf("Headers:\n%v\n\n", r.Header)
	fmt.Printf("Form:\n%v\n\n", r.Form)
	fmt.Printf("URL Params:\n%v\n\n", queryValues)
	fmt.Printf("Body: %s\n\n", bs)
}
