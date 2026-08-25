package service

import (
	"fmt"
	"schoolbusauth/internal/document"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/policy"
	"schoolbusauth/internal/store"
	"time"
)

type Clock interface{ Now() time.Time }
type FixedClock struct{ Value time.Time }

func (c FixedClock) Now() time.Time { return c.Value }

type CompletionCallback func(domain.AuthorizationBundle) error

type Service struct {
	store      *store.Store
	clock      Clock
	html       *document.HTMLRenderer
	text       document.TextRenderer
	onComplete CompletionCallback
	policy     policy.Engine
}

func New(repository *store.Store, clock Clock, callback CompletionCallback) (*Service, error) {
	if repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	if clock == nil {
		return nil, fmt.Errorf("clock is required")
	}
	html, err := document.NewHTMLRenderer()
	if err != nil {
		return nil, err
	}
	return &Service{store: repository, clock: clock, html: html, text: document.TextRenderer{}, onComplete: callback, policy: policy.New("本校")}, nil
}

func (s *Service) Store() *store.Store { return s.store }
func (s *Service) Now() time.Time      { return s.clock.Now() }

func entityID(prefix, value string) string {
	cleaned := make([]rune, 0, len(value))
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z':
			cleaned = append(cleaned, char)
		case char >= 'A' && char <= 'Z':
			cleaned = append(cleaned, char+('a'-'A'))
		case char >= '0' && char <= '9':
			cleaned = append(cleaned, char)
		case char >= 0x4e00 && char <= 0x9fff:
			cleaned = append(cleaned, char)
		}
		if len(cleaned) >= 24 {
			break
		}
	}
	if len(cleaned) == 0 {
		cleaned = append(cleaned, 'x')
	}
	return prefix + "-" + string(cleaned)
}

func eventID(authID, action string, revision int) string {
	return fmt.Sprintf("event-%s-%s-%03d", authID, action, revision)
}
