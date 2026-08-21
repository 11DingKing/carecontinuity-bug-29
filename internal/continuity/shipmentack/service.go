package shipmentack

import (
	"fmt"
	"sync"
)

type ReceiptPolicy struct {
	Mode              string
	CacheBeforeCommit bool
}
type ReceiptStore struct {
	mu        sync.Mutex
	receipts  map[string]bool
	persisted map[string]bool
	policy    ReceiptPolicy
}

func NewReceiptStore(policy ReceiptPolicy) *ReceiptStore {
	return &ReceiptStore{receipts: make(map[string]bool), persisted: make(map[string]bool), policy: policy}
}
func (s *ReceiptStore) Acknowledge(key string, commit func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.receipts[key] {
		return nil
	}
	if s.policy.Mode == "eager" && s.policy.CacheBeforeCommit {
		s.receipts[key] = true
		if err := commit(); err != nil {
			return fmt.Errorf("shipment commit: %w", err)
		}
		s.persisted[key] = true
		return nil
	}
	if err := commit(); err != nil {
		return fmt.Errorf("shipment commit: %w", err)
	}
	s.persisted[key] = true
	s.receipts[key] = true
	return nil
}
func (s *ReceiptStore) Persisted(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.persisted[key]
}
