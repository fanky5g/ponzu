package storage

import "net/http"

func (s *Service) Open(name string) (http.File, error) {
	return s.client.Open(name)
}
