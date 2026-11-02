package tun

import (
	"errors"
	"testing"
)

func TestMemory(t *testing.T) {
	m := NewMemory()
	on, err := m.Current()
	if err != nil {
		t.Fatal(err)
	}
	if on {
		t.Fatal("TUN starts enabled")
	}
	if err := m.Set(true); err != nil {
		t.Fatal(err)
	}
	on, err = m.Current()
	if err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Fatal("TUN did not turn on")
	}
	m.Fail(errors.New("denied"))
	if err := m.Set(false); err == nil {
		t.Fatal("Set succeeded after Fail")
	}
	on, _ = m.Current()
	if !on {
		t.Fatal("failed Set cleared TUN")
	}
}
