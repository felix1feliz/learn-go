package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const PORT = ":3000"

	http.HandleFunc("GET /{$}", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "themes-test.html")
		res.WriteHeader(200)
	})

	http.HandleFunc("GET /styles.css", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "styles.css")
		res.WriteHeader(200)
	})

	http.HandleFunc("GET /theme.css", getTheme)

	listen(PORT)
}

func getTheme(res http.ResponseWriter, req *http.Request) {
	res.Header().Del("Content-Type")
	res.Header().Add("Content-Type", "text/css")

	if req.URL.Query().Has("set-theme") {
		if validateTheme(req.URL.Query().Get("set-theme")) {
			var chosenTheme http.Cookie
			chosenTheme.Name = "theme"
			chosenTheme.Value = req.URL.Query().Get("set-theme")

			setCookie(res, chosenTheme)

			res.WriteHeader(200)
			http.ServeFile(res, req, chosenTheme.Value+".css")
			return
		}

		var defaultTheme http.Cookie
		defaultTheme.Name = "theme"
		defaultTheme.Value = "dark-theme"

		setCookie(res, defaultTheme)
		res.WriteHeader(422)
		http.ServeFile(res, req, "dark-theme.css")
		return
	}
	theme, getCookieErr := req.Cookie("theme")

	if getCookieErr != nil {
		if getCookieErr != http.ErrNoCookie {
			log.Fatal(getCookieErr)
		}

		var defaultTheme http.Cookie
		defaultTheme.Name = "theme"
		defaultTheme.Value = "dark-theme"

		setCookie(res, defaultTheme)
		http.ServeFile(res, req, "dark-theme.css")
		return
	}

	if validateTheme(theme.Value) {
		res.WriteHeader(200)
		http.ServeFile(res, req, theme.Value+".css")
		return
	}

	var defaultTheme http.Cookie
	defaultTheme.Name = "theme"
	defaultTheme.Value = "dark-theme"

	setCookie(res, defaultTheme)
	res.WriteHeader(422)
	http.ServeFile(res, req, "dark-theme.css")
}

func validateTheme(theme string) bool {
	var themes [2]string = [...]string{"dark-theme", "light-theme"}

	for i := range themes {
		if themes[i] == theme {
			return true
		}
	}
	return false
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
