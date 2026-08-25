package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	studentsBucket       = []byte("students")
	routesBucket         = []byte("routes")
	guardiansBucket      = []byte("guardians")
	authorizationsBucket = []byte("authorizations")
	auditsBucket         = []byte("audit_events")
	printsBucket         = []byte("print_records")
)

type Store struct{ db *bolt.DB }

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("store path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	s := &Store{db: db}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{studentsBucket, routesBucket, guardiansBucket, authorizationsBucket, auditsBucket, printsBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", name, err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func putJSON(bucket *bolt.Bucket, key string, value any) error {
	if key == "" {
		return fmt.Errorf("entity id is required")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode entity: %w", err)
	}
	return bucket.Put([]byte(key), data)
}

func getJSON(bucket *bolt.Bucket, key string, destination any) error {
	data := bucket.Get([]byte(key))
	if data == nil {
		return os.ErrNotExist
	}
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("decode entity: %w", err)
	}
	return nil
}

func scanJSON[T any](bucket *bolt.Bucket) ([]T, error) {
	result := make([]T, 0)
	err := bucket.ForEach(func(_, value []byte) error {
		var item T
		if err := json.Unmarshal(value, &item); err != nil {
			return err
		}
		result = append(result, item)
		return nil
	})
	return result, err
}

func missing(name, id string, err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s %q not found", name, id)
	}
	return err
}
