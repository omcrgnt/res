package unique

import (
	"errors"
	"reflect"
	"testing"

	"github.com/omcrgnt/res"
)

func TestAddWithCustomTag_readBack(t *testing.T) {
	reg := New()
	if err := reg.AddWithCustomTag(&testWidget{n: 1}, "ecfg", "APP"); err != nil {
		t.Fatal(err)
	}

	entries := reg.GetByType(reflect.TypeFor[*testWidget]())
	val, ok := entries[0].GetCustomTag("ecfg")
	if !ok || val != "APP" {
		t.Fatalf("custom tag = %v, ok=%v", val, ok)
	}
}

func TestAddWithCustomTag_emptyKey(t *testing.T) {
	reg := New()
	if err := reg.AddWithCustomTag(&testWidget{}, "", "x"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestChangeValue_specToBuilt(t *testing.T) {
	reg := New()

	type spec struct{ n int }
	type built struct{ n int }

	if err := reg.reg.Add(&spec{n: 1}); err != nil {
		t.Fatal(err)
	}

	var entry res.Entry
	reg.WalkEntries(func(e res.Entry) bool {
		entry = e
		return false
	})

	if err := entry.ChangeValue(&built{n: 2}); err != nil {
		t.Fatal(err)
	}
	if entry.Value().(*built).n != 2 {
		t.Fatalf("got %v", entry.Value())
	}
}

func TestChangeValue_duplicateBuiltType(t *testing.T) {
	reg := New()

	type specA struct{}
	type specB struct{}
	type built struct{}

	if err := reg.reg.Add(&specA{}); err != nil {
		t.Fatal(err)
	}
	if err := reg.reg.Add(&specB{}); err != nil {
		t.Fatal(err)
	}

	var specs []res.Entry
	reg.WalkEntries(func(e res.Entry) bool {
		specs = append(specs, e)
		return true
	})

	if err := specs[0].ChangeValue(&built{}); err != nil {
		t.Fatal(err)
	}
	if err := specs[1].ChangeValue(&built{}); !errors.Is(err, ErrRegularExists) {
		t.Fatalf("second ChangeValue = %v, want ErrRegularExists", err)
	}
}

func TestMerge_preservesCustomTag(t *testing.T) {
	dst := New()
	src := New()

	if err := src.AddWithCustomTag(&testWidget{n: 1}, "ecfg", "ITEM"); err != nil {
		t.Fatal(err)
	}
	if err := Merge(dst, src); err != nil {
		t.Fatal(err)
	}

	entries := dst.GetByType(reflect.TypeFor[*testWidget]())
	val, ok := entries[0].GetCustomTag("ecfg")
	if !ok || val != "ITEM" {
		t.Fatalf("custom tag = %v, ok=%v", val, ok)
	}
}
