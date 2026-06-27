package res

import (
	"reflect"
	"testing"
)

func TestAddWithTagsAndCustomTags_roundtrip(t *testing.T) {
	reg := New()
	if err := AddWithTagsAndCustomTags(reg, "spec", map[string]any{"ecfg": "SERVICE_ITEM"}, TagRegular); err != nil {
		t.Fatal(err)
	}

	var got Entry
	reg.WalkEntries(func(e Entry) bool {
		got = e
		return false
	})
	if got == nil {
		t.Fatal("no entry")
	}
	val, ok := got.GetCustomTag("ecfg")
	if !ok || val != "SERVICE_ITEM" {
		t.Fatalf("custom tag = %v, ok=%v", val, ok)
	}
	if !got.Regular() {
		t.Fatal("expected regular")
	}
}

func TestEntry_ChangeValue_preservesTagsAndCustomTags(t *testing.T) {
	reg := New()
	spec := &testChangeSpec{n: 1}
	if err := AddWithTagsAndCustomTags(reg, spec, map[string]any{"ecfg": "BLOCK"}, TagRegular); err != nil {
		t.Fatal(err)
	}

	var entry Entry
	reg.WalkEntries(func(e Entry) bool {
		entry = e
		return false
	})

	built := &testChangeBuilt{n: 2}
	if err := entry.ChangeValue(built); err != nil {
		t.Fatal(err)
	}
	if entry.Value() != built {
		t.Fatalf("value = %v", entry.Value())
	}
	if entry.Type() != reflect.TypeFor[*testChangeBuilt]() {
		t.Fatalf("type = %v", entry.Type())
	}
	val, ok := entry.GetCustomTag("ecfg")
	if !ok || val != "BLOCK" {
		t.Fatalf("custom tag after change = %v", val)
	}
	if !entry.Regular() {
		t.Fatal("regular tag lost")
	}
}

func TestReplaceAtType_preservesCustomTags(t *testing.T) {
	reg := New()
	if err := AddWithTagsAndCustomTags(reg, "old", map[string]any{"k": "v"}, TagReplaceable); err != nil {
		t.Fatal(err)
	}
	if err := ReplaceAtType(reg, "new", TagRegular); err != nil {
		t.Fatal(err)
	}

	entries := reg.GetByType(reflect.TypeFor[string]())
	if len(entries) != 1 {
		t.Fatalf("entries = %d", len(entries))
	}
	val, ok := entries[0].GetCustomTag("k")
	if !ok || val != "v" {
		t.Fatalf("custom tag = %v, ok=%v", val, ok)
	}
}

func TestTransform_preservesCustomTags(t *testing.T) {
	reg := New()
	if err := AddWithTagsAndCustomTags(reg, &testChangeBuilt{n: 1}, map[string]any{"k": 1}, TagRegular); err != nil {
		t.Fatal(err)
	}
	if err := reg.Transform(func(v any) any {
		return &testChangeBuilt{n: v.(*testChangeBuilt).n + 1}
	}); err != nil {
		t.Fatal(err)
	}

	entries := reg.GetByType(reflect.TypeFor[*testChangeBuilt]())
	val, ok := entries[0].GetCustomTag("k")
	if !ok || val != 1 {
		t.Fatalf("custom tag = %v", val)
	}
	if entries[0].Value().(*testChangeBuilt).n != 2 {
		t.Fatalf("transform value = %v", entries[0].Value())
	}
}

type testChangeSpec struct{ n int }

type testChangeBuilt struct{ n int }
