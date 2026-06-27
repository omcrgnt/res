package res

import (
	"reflect"
	"testing"
)

type replaceWidget struct{ n int }

func TestTagRegular(t *testing.T) {
	resetGlobalRegistry()

	if err := Global().AddWithTags(&replaceWidget{n: 1}, TagRegular); err != nil {
		t.Fatal(err)
	}

	entries := Global().GetByType(reflect.TypeFor[*replaceWidget]())
	if len(entries) != 1 || !entries[0].Regular() {
		t.Fatalf("expected TagRegular entry, got %+v", entries)
	}
}

func TestReplaceAtType_replacesInPlace(t *testing.T) {
	reg := New()

	first := &replaceWidget{n: 1}
	second := &replaceWidget{n: 2}
	other := &replaceWidget{n: 99}

	if err := reg.AddWithTags(first, TagReplaceable); err != nil {
		t.Fatal(err)
	}
	if err := reg.Add(other); err != nil {
		t.Fatal(err)
	}

	if err := ReplaceAtType(reg, second, TagRegular); err != nil {
		t.Fatal(err)
	}

	var order []any
	reg.WalkEntries(func(e Entry) bool {
		order = append(order, e.Value())
		return true
	})
	if len(order) != 2 || order[0] != second || order[1] != other {
		t.Fatalf("order = %v, want [%p, %p]", order, second, other)
	}

	entries := reg.GetByType(reflect.TypeFor[*replaceWidget]())
	if len(entries) != 2 {
		t.Fatalf("expected 2 widget entries, got %d", len(entries))
	}
	if entries[0].Value() != second || !entries[0].Regular() {
		t.Fatalf("first widget entry = %+v", entries[0])
	}
}

func TestReplaceAtType_appendsWhenMissing(t *testing.T) {
	reg := New()
	w := &replaceWidget{n: 1}

	if err := ReplaceAtType(reg, w, TagRegular); err != nil {
		t.Fatal(err)
	}

	got, err := reg.GetOneByType(reflect.TypeFor[*replaceWidget]())
	if err != nil || got != w {
		t.Fatalf("GetOneByType = %v, %v", got, err)
	}
}

func TestReplaceAtType_nilValue(t *testing.T) {
	if err := ReplaceAtType(New(), nil, TagRegular); err == nil {
		t.Fatal("expected error for nil value")
	}
}
