package service

import (
	"net/http"
	"html/template"

	consul "github.com/hashicorp/consul/api"
)

type Service struct {
	Name        	string
	Port			int
	ConsulAgent		*consul.Agent
	ConsulHealth	*consul.Health
}

func NewService(name string, port int) (*Service, error) {
	s := new(Service)
	s.Name = name
	s.Port = port

	s.registerConsul()

	return s, nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	num := s.countConsulServices()

	t := template.Must(template.ParseFiles("tmpl/homepage.html"))
	t.Execute(w, num)
}