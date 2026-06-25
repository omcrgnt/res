package restest_test

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/restest"
)

func TestWith_registersInOrder(t *testing.T) {
	reg := restest.With("a", 1, "b")

	var values []any
	reg.WalkEntries(func(e res.Entry) bool {
		values = append(values, e.Value)
		return true
	})
	if len(values) != 3 {
		t.Fatalf("got len %d", len(values))
	}
	for i, want := range []any{"a", 1, "b"} {
		w, ok := values[i].(restest.Wire)
		if !ok || w.V != want {
			t.Fatalf("values[%d] = %T(%v), want wire %v", i, values[i], values[i], want)
		}
	}
}

func TestWithTagged(t *testing.T) {
	reg := restest.WithTagged("x", res.TagReplaceable)

	reg.WalkEntries(func(e res.Entry) bool {
		if !e.Replaceable() {
			t.Fatal("expected replaceable")
		}
		return false
	})
}

func TestResetGlobal_isolatedFromPriorGlobal(t *testing.T) {
	_ = res.Global().Add(restest.Wire{V: "leak"})
	reg := restest.ResetGlobal()
	_ = reg.Add(restest.Wire{V: "only"})

	n := 0
	reg.WalkEntries(func(res.Entry) bool {
		n++
		return true
	})
	if n != 1 {
		t.Fatalf("expected 1 entry, got %d", n)
	}
}

func TestAddAll_stopsOnError(t *testing.T) {
	reg := restest.Registry()
	if err := restest.AddAll(reg, "ok"); err != nil {
		t.Fatal(err)
	}
	if err := reg.Add(nil); err == nil {
		t.Fatal("expected nil add error")
	}
	if got := reg.GetByType(reflect.TypeFor[restest.Wire]()); len(got) != 1 {
		t.Fatalf("expected one entry before failure, got %d", len(got))
	}
}
