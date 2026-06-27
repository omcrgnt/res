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
		values = append(values, e.Value())
		return true
	})
	if len(values) != 3 || values[0] != "a" || values[1] != 1 || values[2] != "b" {
		t.Fatalf("got %v", values)
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
	_ = res.Global().Add("leak")
	reg := restest.ResetGlobal()
	_ = reg.Add("only")

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
	if err := restest.AddAll(reg, "ok", nil); err == nil {
		t.Fatal("expected nil add error")
	}
	if got := reg.GetByType(reflect.TypeFor[string]()); len(got) != 1 {
		t.Fatalf("expected one entry before failure, got %d", len(got))
	}
}
