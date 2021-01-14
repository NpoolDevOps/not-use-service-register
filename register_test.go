package srvreg

import (
	"testing"
)

func TestRegister(t *testing.T) {
	Register("auth.npool.com", "127.0.0.1:18899")
}

func TestQuery(t *testing.T) {
	Query("auth.npool.com")
}

func TestWatch(t *testing.T) {
	Watch("auth.npool.com", func(ev Event) {
		t.Logf("event: %v", ev)
	})
}
