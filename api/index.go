package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.URL.Query().Get("u")
	if targetUrl == "" {
		http.Error(w, "url ?u= is required", http.StatusBadRequest)
		return
	}

	url, err := url.Parse(targetUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &http.Request{
		URL:    url,
		Method: r.Method,
		Body:   r.Body,
		Header: r.Header,
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	addCorsHeaders(w)
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func addCorsHeaders(w http.ResponseWriter) {
	delete(w.Header(), "Access-Control-Allow-Origin")
	delete(w.Header(), "Access-Control-Allow-Methods")
	delete(w.Header(), "Access-Control-Allow-Headers")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		dst.Add(k, strings.Join(vv, ","))
	}
}
