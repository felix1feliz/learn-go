package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	const PORT = ":3000"

	http.HandleFunc("GET /{$}", func(res http.ResponseWriter, req *http.Request) {
		var urlQueryValues url.Values = req.URL.Query()
		if urlQueryValues.Has("id") {
			var cookie http.Cookie
			cookie.Name = "id"
			cookie.Value = urlQueryValues.Get("id")
			setCookie(res, cookie)
		}
		res.WriteHeader(200)
	})

	listen(PORT)
}

func setCookie(res http.ResponseWriter, c http.Cookie) {
	var c_valid error = c.Valid()
	if c_valid != nil {
		log.Fatal(c_valid)
	}
	res.Header().Add("Set-Cookie", c.String())
}

func listen(PORT string) {
	fmt.Printf("Server listening on port '%s'\n", PORT)
	var listen_err error = http.ListenAndServe(PORT, nil)
	if listen_err != nil {
		log.Fatal(listen_err)
	}
}
