package shipmentack

type Coordinator struct{ store *ReceiptStore }

func NewCoordinator() *Coordinator {
	return &Coordinator{store: NewReceiptStore(ReceiptPolicy{Mode: "eager", CacheBeforeCommit: true})}
}
func (c *Coordinator) Acknowledge(key string, commit func() error) error {
	return c.store.Acknowledge(key, commit)
}
func (c *Coordinator) Persisted(key string) bool { return c.store.Persisted(key) }
