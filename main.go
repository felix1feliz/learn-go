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

	http.HandleFunc("GET /theme.css", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Del("Content-Type")
		res.Header().Add("Content-Type", "text/css")
		var defaultTheme http.Cookie
		defaultTheme.Name = "theme"
		defaultTheme.Value = "dark-theme"
		defaultTheme.SameSite = http.SameSiteNoneMode
		defaultTheme.Secure = true
		theme, getCookieErr := req.Cookie("theme")

		if getCookieErr != nil {
			if getCookieErr != http.ErrNoCookie {
				log.Fatal(getCookieErr)
			}

			setCookie(res, defaultTheme)
			http.ServeFile(res, req, "dark-theme.css")
			return
		}

		if validateTheme(theme.Value) {
			res.WriteHeader(200)
			http.ServeFile(res, req, theme.Value+".css")
			return
		}

		setCookie(res, defaultTheme)
		res.WriteHeader(422)
	})

	http.HandleFunc("GET /dark-theme.css", getTheme("dark-theme"))
	http.HandleFunc("GET /light-theme.css", getTheme("light-theme"))
	http.HandleFunc("GET /hacker-theme.css", getTheme("hacker-theme"))

	listen(PORT)
}

func getTheme(theme string) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Del("Content-Type")
		res.Header().Add("Content-Type", "text/css")

		if req.URL.Query().Has("set-theme") {
			var themeCookie http.Cookie
			themeCookie.Name = "theme"
			themeCookie.Value = theme
			themeCookie.SameSite = http.SameSiteNoneMode
			themeCookie.Secure = true
			setCookie(res, themeCookie)
			return
		}

		http.ServeFile(res, req, theme+".css")
	}
}

func validateTheme(theme string) bool {
	var themes [3]string = [...]string{"dark-theme", "light-theme", "hacker-theme"}

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
