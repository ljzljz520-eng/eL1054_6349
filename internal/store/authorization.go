package store

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"schoolbusauth/internal/domain"
)

func (s *Store) CreateBundle(bundle domain.AuthorizationBundle, audit domain.AuditEvent) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(authorizationsBucket).Get([]byte(bundle.Authorization.ID)) != nil {
			return fmt.Errorf("authorization already exists")
		}
		items := []struct {
			bucket []byte
			id     string
			value  any
		}{
			{studentsBucket, bundle.Student.ID, bundle.Student},
			{routesBucket, bundle.Route.ID, bundle.Route},
			{guardiansBucket, bundle.Guardian.ID, bundle.Guardian},
			{authorizationsBucket, bundle.Authorization.ID, bundle.Authorization},
			{auditsBucket, audit.ID, audit},
		}
		for _, item := range items {
			if err := putJSON(tx.Bucket(item.bucket), item.id, item.value); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Authorization(id string) (domain.Authorization, error) {
	var value domain.Authorization
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(authorizationsBucket), id, &value) })
	return value, missing("authorization", id, err)
}

func (s *Store) Bundle(id string) (domain.AuthorizationBundle, error) {
	var result domain.AuthorizationBundle
	err := s.db.View(func(tx *bolt.Tx) error {
		if err := getJSON(tx.Bucket(authorizationsBucket), id, &result.Authorization); err != nil {
			return err
		}
		if err := getJSON(tx.Bucket(studentsBucket), result.Authorization.StudentID, &result.Student); err != nil {
			return err
		}
		if err := getJSON(tx.Bucket(routesBucket), result.Authorization.RouteID, &result.Route); err != nil {
			return err
		}
		return getJSON(tx.Bucket(guardiansBucket), result.Authorization.GuardianID, &result.Guardian)
	})
	return result, missing("authorization bundle", id, err)
}

func (s *Store) Bundles() ([]domain.AuthorizationBundle, error) {
	values := make([]domain.AuthorizationBundle, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		auths, err := scanJSON[domain.Authorization](tx.Bucket(authorizationsBucket))
		if err != nil {
			return err
		}
		for _, auth := range auths {
			bundle := domain.AuthorizationBundle{Authorization: auth}
			if err := getJSON(tx.Bucket(studentsBucket), auth.StudentID, &bundle.Student); err != nil {
				return err
			}
			if err := getJSON(tx.Bucket(routesBucket), auth.RouteID, &bundle.Route); err != nil {
				return err
			}
			if err := getJSON(tx.Bucket(guardiansBucket), auth.GuardianID, &bundle.Guardian); err != nil {
				return err
			}
			values = append(values, bundle)
		}
		return nil
	})
	return values, err
}

func (s *Store) UpdateAuthorization(value domain.Authorization, audit domain.AuditEvent) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(authorizationsBucket)
		if bucket.Get([]byte(value.ID)) == nil {
			return fmt.Errorf("authorization not found")
		}
		if err := putJSON(bucket, value.ID, value); err != nil {
			return err
		}
		return putJSON(tx.Bucket(auditsBucket), audit.ID, audit)
	})
}
