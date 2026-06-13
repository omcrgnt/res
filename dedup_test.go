package res

import (
	"errors"
	"reflect"
	"testing"

	"github.com/omcrgnt/sdi"
)

type dedupPort interface {
	Port() string
}

type systemOut struct{}

func (systemOut) Port() string { return "system" }

type userOut struct{}

func (userOut) Port() string { return "user" }

type altSystemOut struct{}

func (altSystemOut) Port() string { return "alt" }

func TestDedup_systemUser(t *testing.T) {
	resetGlobalRegistry()

	if err := AddBuiltin(systemOut{}); err != nil {
		t.Fatal(err)
	}
	if err := Add(userOut{}); err != nil {
		t.Fatal(err)
	}

	iface := reflect.TypeFor[dedupPort]()
	if err := Dedup([]reflect.Type{iface}, sdi.DefaultDedupPolicy); err != nil {
		t.Fatal(err)
	}

	found := Find[dedupPort]()
	if len(found) != 1 {
		t.Fatalf("expected 1 implementor, got %d", len(found))
	}
	if _, ok := found[0].(userOut); !ok {
		t.Fatalf("expected userOut, got %T", found[0])
	}
}

func TestDedup_removeSystemRejectsUser(t *testing.T) {
	resetGlobalRegistry()

	u := userOut{}
	if err := Add(u); err != nil {
		t.Fatal(err)
	}

	err := Dedup([]reflect.Type{reflect.TypeFor[dedupPort]()}, func(ctx sdi.DedupContext) error {
		return ctx.Remove(u)
	})
	if err == nil {
		t.Fatal("expected error removing user resource")
	}
}

func TestDedup_twoSystemSamePort(t *testing.T) {
	resetGlobalRegistry()

	if err := AddBuiltin(systemOut{}); err != nil {
		t.Fatal(err)
	}
	if err := AddBuiltin(altSystemOut{}); err != nil {
		t.Fatal(err)
	}

	err := Dedup([]reflect.Type{reflect.TypeFor[dedupPort]()}, sdi.DefaultDedupPolicy)
	if !errors.Is(err, sdi.ErrMultipleSystemDefaults) {
		t.Fatalf("expected multiple system defaults, got %v", err)
	}

	if n := len(Find[dedupPort]()); n != 2 {
		t.Fatalf("pool should be unchanged, got %d implementors", n)
	}
}

func TestDedup_executorCallsPolicy(t *testing.T) {
	resetGlobalRegistry()
	_ = Add("x")

	called := false
	err := Dedup([]reflect.Type{reflect.TypeFor[dedupPort]()}, func(ctx sdi.DedupContext) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("policy not called")
	}
}
