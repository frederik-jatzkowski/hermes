package util

import (
	"net/http"
	"time"
)

type HTTPSRedirect struct{}

func (r *HTTPSRedirect) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	url := "https://" + req.Host + req.URL.Path
	http.Redirect(res, req, url, http.StatusMovedPermanently)
}

func ServeHTTPSRedirect() {
	s := http.Server{
		Addr:           "0.0.0.0:http",
		Handler:        &HTTPSRedirect{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
