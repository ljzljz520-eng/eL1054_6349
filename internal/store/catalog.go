package store

import (
	bolt "go.etcd.io/bbolt"
	"schoolbusauth/internal/domain"
)

func (s *Store) SaveStudent(value domain.Student) error {
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(studentsBucket), value.ID, value) })
}

func (s *Store) Student(id string) (domain.Student, error) {
	var value domain.Student
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(studentsBucket), id, &value) })
	return value, missing("student", id, err)
}

func (s *Store) Students() ([]domain.Student, error) {
	var values []domain.Student
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		values, err = scanJSON[domain.Student](tx.Bucket(studentsBucket))
		return err
	})
	return values, err
}

func (s *Store) SaveRoute(value domain.Route) error {
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(routesBucket), value.ID, value) })
}

func (s *Store) Route(id string) (domain.Route, error) {
	var value domain.Route
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(routesBucket), id, &value) })
	return value, missing("route", id, err)
}

func (s *Store) Routes() ([]domain.Route, error) {
	var values []domain.Route
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		values, err = scanJSON[domain.Route](tx.Bucket(routesBucket))
		return err
	})
	return values, err
}

func (s *Store) SaveGuardian(value domain.Guardian) error {
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(guardiansBucket), value.ID, value) })
}

func (s *Store) Guardian(id string) (domain.Guardian, error) {
	var value domain.Guardian
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(guardiansBucket), id, &value) })
	return value, missing("guardian", id, err)
}

func (s *Store) Guardians() ([]domain.Guardian, error) {
	var values []domain.Guardian
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		values, err = scanJSON[domain.Guardian](tx.Bucket(guardiansBucket))
		return err
	})
	return values, err
}
