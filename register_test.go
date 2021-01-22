package srvreg

import (
	"testing"
)

func TestRegister(t *testing.T) {
	Register("domain_auth.npool.com", "131.0.0.1:18899")
}

func TestQuery(t *testing.T) {
	Query("domain_auth.npool.com")
}

func TestBatchQuery(t *testing.T) {
	BatchQuery([]string{"domain_auth.npool.com"})
}

func TestWatch(t *testing.T) {
	Watch("domain_auth.npool.com", func(ev Event) {
		t.Logf("event: %v", ev)
	})
}

