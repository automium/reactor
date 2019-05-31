package service

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
)

const (
	serviceName = "example"
)

var (
	address string
	path string
)

func Run(_address, _path string) error {
	address = _address
	path = _path

	a := strings.Split(address, ":")
	port, err := strconv.Atoi(a[1])
	if err != nil {
		return err
	}

	s, err := NewService(serviceName, port)
	if err != nil {
		return err
	}

	fs := http.FileServer(http.Dir("static"))
  	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.Handle("/", s)
	return http.ListenAndServe(address, nil)
}