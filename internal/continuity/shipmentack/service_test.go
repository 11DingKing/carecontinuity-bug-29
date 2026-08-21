package shipmentack_test

import (
	"carecontinuity/internal/continuity/shipmentack"
	"errors"
	"testing"
)

func TestMedicineShipmentAckPublicBehavior(t *testing.T) {
	c := shipmentack.NewCoordinator()
	calls := 0
	failure := errors.New("database unavailable")
	if err := c.Acknowledge("shipment-9", func() error { calls++; return failure }); !errors.Is(err, failure) {
		t.Fatalf("expected first failure, got %v", err)
	}
	if err := c.Acknowledge("shipment-9", func() error { calls++; return nil }); err != nil {
		t.Fatalf("retry failed: %v", err)
	}
	if calls != 2 {
		t.Fatalf("retry commit called %d times total", calls)
	}
	if !c.Persisted("shipment-9") {
		t.Fatal("successful retry was not persisted")
	}
}
