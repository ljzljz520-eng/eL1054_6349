package store

import (
	bolt "go.etcd.io/bbolt"
	"schoolbusauth/internal/domain"
	"sort"
)

func (s *Store) SavePrintRecord(value domain.PrintRecord) error {
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(printsBucket), value.ID, value) })
}

func (s *Store) PrintRecords(authorizationID string) ([]domain.PrintRecord, error) {
	var values []domain.PrintRecord
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		values, err = scanJSON[domain.PrintRecord](tx.Bucket(printsBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	filtered := values[:0]
	for _, value := range values {
		if value.AuthorizationID == authorizationID {
			filtered = append(filtered, value)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].PrintedAt.Before(filtered[j].PrintedAt) })
	return filtered, nil
}

func (s *Store) AuditEvents(authorizationID string) ([]domain.AuditEvent, error) {
	var values []domain.AuditEvent
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		values, err = scanJSON[domain.AuditEvent](tx.Bucket(auditsBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	filtered := values[:0]
	for _, value := range values {
		if value.AuthorizationID == authorizationID {
			filtered = append(filtered, value)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].OccurredAt.Before(filtered[j].OccurredAt) })
	return filtered, nil
}

func (s *Store) DeleteDraft(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		var auth domain.Authorization
		if err := getJSON(tx.Bucket(authorizationsBucket), id, &auth); err != nil {
			return err
		}
		if auth.Status != domain.StatusDraft {
			return bolt.ErrTxNotWritable
		}
		return tx.Bucket(authorizationsBucket).Delete([]byte(id))
	})
}
