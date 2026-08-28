package service

import (
	"fmt"
	"schoolbusauth/internal/document"
)

type Preview struct {
	View document.View
	HTML []byte
}

func (s *Service) Preview(id string) (Preview, error) {
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return Preview{}, err
	}
	view := document.BuildView(bundle, s.clock.Now())
	html, err := s.html.Render(view)
	if err != nil {
		return Preview{}, err
	}
	return Preview{View: view, HTML: html}, nil
}

func (s *Service) DownloadHTML(id string) ([]byte, string, error) {
	preview, err := s.Preview(id)
	if err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("authorization-%s.html", id)
	return preview.HTML, filename, nil
}

func (s *Service) DownloadText(id string) ([]byte, string, error) {
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return nil, "", err
	}
	view := document.BuildView(bundle, s.clock.Now())
	content, err := s.text.Render(view)
	if err != nil {
		return nil, "", err
	}
	return content, fmt.Sprintf("authorization-%s.txt", id), nil
}
