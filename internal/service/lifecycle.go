package service

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/validation"
)

func (s *Service) Issue(id, actor string) (domain.AuthorizationBundle, error) {
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if err := bundle.Authorization.CanIssue(); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if err := s.policy.EvaluateIssue(bundle, s.clock.Now()).BlockingError(); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if bundle.Authorization.IsExpired(s.clock.Now()) {
		return domain.AuthorizationBundle{}, fmt.Errorf("cannot issue an expired authorization")
	}
	now := s.clock.Now().UTC()
	bundle.Authorization.Status = domain.StatusIssued
	bundle.Authorization.IssuedAt = &now
	bundle.Authorization.Revision++
	audit := domain.AuditEvent{ID: eventID(id, "issued", bundle.Authorization.Revision), AuthorizationID: id, Action: "issued", Actor: actor, OccurredAt: now, Details: "authorization approved for gate use"}
	if err := s.store.UpdateAuthorization(bundle.Authorization, audit); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	return bundle, nil
}

func (s *Service) Revoke(id, actor, reason string) (domain.AuthorizationBundle, error) {
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if err := bundle.Authorization.CanRevoke(); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if reason == "" {
		return domain.AuthorizationBundle{}, fmt.Errorf("revocation reason is required")
	}
	now := s.clock.Now().UTC()
	bundle.Authorization.Status = domain.StatusRevoked
	bundle.Authorization.RevokedAt = &now
	bundle.Authorization.Revision++
	audit := domain.AuditEvent{ID: eventID(id, "revoked", bundle.Authorization.Revision), AuthorizationID: id, Action: "revoked", Actor: actor, OccurredAt: now, Details: reason}
	if err := s.store.UpdateAuthorization(bundle.Authorization, audit); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	return bundle, nil
}

func (s *Service) Complete(id string) error {
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return err
	}
	if bundle.Authorization.Status != domain.StatusIssued {
		return fmt.Errorf("only an issued authorization can be completed")
	}
	return s.onComplete(bundle)
}

func (s *Service) RecordPrint(id, printer string, copies int) (domain.PrintRecord, error) {
	if err := validation.Copies(copies); err != nil {
		return domain.PrintRecord{}, err
	}
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return domain.PrintRecord{}, err
	}
	if !bundle.Authorization.IsUsable(s.clock.Now()) {
		return domain.PrintRecord{}, fmt.Errorf("authorization is not usable today")
	}
	if printer == "" {
		return domain.PrintRecord{}, fmt.Errorf("printer is required")
	}
	record := domain.PrintRecord{ID: fmt.Sprintf("print-%s-%03d", id, bundle.Authorization.Revision+copies), AuthorizationID: id, Printer: printer, Copies: copies, PrintedAt: s.clock.Now().UTC()}
	if err := s.store.SavePrintRecord(record); err != nil {
		return domain.PrintRecord{}, err
	}
	return record, nil
}
