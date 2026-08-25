package query

import (
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/store"
	"sort"
)

type Service struct{ store *store.Store }

func New(repository *store.Store) *Service { return &Service{store: repository} }

func (s *Service) Search(filter domain.SearchFilter) ([]domain.AuthorizationBundle, error) {
	values, err := s.store.Bundles()
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuthorizationBundle, 0)
	for _, value := range values {
		if value.Matches(filter) {
			result = append(result, value)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Authorization.AuthorizedDate == result[j].Authorization.AuthorizedDate {
			return result[i].Student.Name < result[j].Student.Name
		}
		return result[i].Authorization.AuthorizedDate > result[j].Authorization.AuthorizedDate
	})
	return result, nil
}

func (s *Service) Detail(id string) (Detail, error) {
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return Detail{}, err
	}
	audits, err := s.store.AuditEvents(id)
	if err != nil {
		return Detail{}, err
	}
	prints, err := s.store.PrintRecords(id)
	if err != nil {
		return Detail{}, err
	}
	return Detail{Bundle: bundle, AuditEvents: audits, PrintRecords: prints}, nil
}

type Detail struct {
	Bundle       domain.AuthorizationBundle `json:"bundle"`
	AuditEvents  []domain.AuditEvent        `json:"audit_events"`
	PrintRecords []domain.PrintRecord       `json:"print_records"`
}
