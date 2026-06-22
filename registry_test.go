package res

import (
	"errors"
	"reflect"
	"testing"
)

func TestRegistry_Add(t *testing.T) {
	r := New().(*registry)
	val := "hello"

	if err := r.Add(val); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	entries := r.GetByType(reflect.TypeFor[string]())
	if len(entries) != 1 || entries[0].Value != "hello" {
		t.Errorf("GetByType failed: got %v", entries)
	}

	if err := r.Add("world"); err != nil {
		t.Fatal("duplicate Add should succeed")
	}
	if len(r.GetByType(reflect.TypeFor[string]())) != 2 {
		t.Error("expected 2 entries for same type")
	}
}

func TestAdd_nil(t *testing.T) {
	r := New().(*registry)
	if err := r.Add(nil); err == nil {
		t.Fatal("expected error for nil resource")
	}
	if err := r.AddWithTags(nil, TagReplaceable); err == nil {
		t.Fatal("expected error for nil tagged resource")
	}
}

func TestAddWithTags_requiresTag(t *testing.T) {
	resetGlobalRegistry()
	if err := Global().AddWithTags("x"); err == nil {
		t.Fatal("expected error when no tags")
	}
}

func TestAddWithTags(t *testing.T) {
	resetGlobalRegistry()

	if err := Global().AddWithTags("default", TagReplaceable); err != nil {
		t.Fatal(err)
	}
	if err := Global().AddWithTags("again", TagReplaceable); err != nil {
		t.Fatal("duplicate AddWithTags should succeed")
	}

	var tagged bool
	Global().WalkEntries(func(e Entry) bool {
		tagged = e.Has(TagReplaceable)
		return false
	})
	if !tagged {
		t.Fatalf("expected TagReplaceable")
	}
}

func TestAddWithTags_dedupesTags(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().AddWithTags("x", TagReplaceable, TagReplaceable)

	entries := Global().GetByType(reflect.TypeFor[string]())
	if len(entries[0].Tags()) != 1 {
		t.Fatalf("expected 1 unique tag, got %v", entries[0].Tags())
	}
}

func TestAdd_and_AddWithTags_sameType(t *testing.T) {
	resetGlobalRegistry()

	type widget struct{ n int }

	if err := Global().AddWithTags(&widget{n: 1}, TagReplaceable); err != nil {
		t.Fatal(err)
	}
	if err := Global().Add(&widget{n: 2}); err != nil {
		t.Fatal(err)
	}

	entries := Global().GetByType(reflect.TypeFor[*widget]())
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].Replaceable() || entries[1].Replaceable() {
		t.Fatalf("tag flags: %+v", entries)
	}
}

func TestGlobal_Add(t *testing.T) {
	resetGlobalRegistry()

	if err := Global().Add("via-default"); err != nil {
		t.Fatalf("Global().Add failed: %v", err)
	}

	entries := Global().GetByType(reflect.TypeFor[string]())
	if len(entries) != 1 || entries[0].Value != "via-default" {
		t.Errorf("unexpected entries: %v", entries)
	}
}

func TestRemove(t *testing.T) {
	resetGlobalRegistry()

	s := "remove-me"
	if err := Global().Add(s); err != nil {
		t.Fatal(err)
	}
	if err := Global().Remove("other"); err == nil {
		t.Fatal("expected not found")
	}
	if err := Global().Remove(s); err != nil {
		t.Fatal(err)
	}
	if len(Global().GetByType(reflect.TypeFor[string]())) != 0 {
		t.Fatal("expected removed")
	}
}

func TestGlobal_WalkEntries(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().Add("a")
	_ = Global().Add(2)

	var seen []any
	Global().WalkEntries(func(e Entry) bool {
		seen = append(seen, e.Value)
		return true
	})

	if len(seen) != 2 {
		t.Fatalf("WalkEntries expected 2, got %d", len(seen))
	}
}

func TestRegistry_WalkEntries(t *testing.T) {
	r := New().(*registry)
	_ = r.Add(10)
	_ = r.Add("string")

	count := 0
	r.WalkEntries(func(Entry) bool {
		count++
		return true
	})

	if count != 2 {
		t.Errorf("WalkEntries failed: expected 2 items, got %d", count)
	}
}

func TestRegistry_WalkEntries_stopsOnFalse(t *testing.T) {
	r := New().(*registry)
	_ = r.Add("one")
	_ = r.Add(2)
	_ = r.Add(true)

	var seen []any
	r.WalkEntries(func(e Entry) bool {
		seen = append(seen, e.Value)
		return e.Value != 2
	})

	if len(seen) != 2 {
		t.Fatalf("WalkEntries must stop after fn returns false, got %v", seen)
	}
	if seen[0] != "one" || seen[1] != 2 {
		t.Fatalf("unexpected visit order: %v", seen)
	}
}

func TestGlobal_shared(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().Add(42)

	entries := Global().GetByType(reflect.TypeFor[int]())
	if len(entries) != 1 || entries[0].Value != 42 {
		t.Fatalf("Global failed: got %v", entries)
	}
}

type Shaper interface {
	Area() int
}

type Square struct {
	Side int
}

func (s *Square) Area() int {
	return s.Side * s.Side
}

type wrappedSquare struct {
	*Square
}

func TestGetByType(t *testing.T) {
	resetGlobalRegistry()

	sq := &Square{Side: 10}
	_ = Global().Add(sq)

	entries := Global().GetByType(reflect.TypeFor[*Square]())
	if len(entries) != 1 {
		t.Errorf("GetByType[*Square] failed: expected 1 match, got %d", len(entries))
	}

	strings := Global().GetByType(reflect.TypeFor[string]())
	if len(strings) != 0 {
		t.Errorf("GetByType[string] failed: expected 0 matches, got %d", len(strings))
	}
}

func TestGetOneByType(t *testing.T) {
	resetGlobalRegistry()

	first := &Square{Side: 1}
	second := &Square{Side: 2}
	_ = Global().Add(first)
	_ = Global().Add(second)

	got, err := Global().GetOneByType(reflect.TypeFor[*Square]())
	if err != nil {
		t.Fatal(err)
	}
	if got != first {
		t.Fatalf("GetOneByType want first added %p, got %p", first, got)
	}

	_, err = Global().GetOneByType(reflect.TypeFor[string]())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

type writePort interface {
	Write([]byte) (int, error)
}

type stdoutSink struct{}

func (stdoutSink) Write(p []byte) (int, error) { return len(p), nil }

type fileSink struct{}

func (fileSink) Write(p []byte) (int, error) { return len(p), nil }

func TestGetOneByInterface(t *testing.T) {
	resetGlobalRegistry()

	first := stdoutSink{}
	second := fileSink{}
	_ = Global().AddWithTags(first, TagReplaceable)
	_ = Global().Add(second)

	got, err := Global().GetOneByInterface(reflect.TypeFor[writePort]())
	if err != nil {
		t.Fatal(err)
	}
	if got != first {
		t.Fatalf("GetOneByInterface want first registered %v, got %v", first, got)
	}

	_, err = Global().GetOneByInterface(reflect.TypeFor[Shaper]())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetByInterface(t *testing.T) {
	resetGlobalRegistry()

	_ = Global().AddWithTags(stdoutSink{}, TagReplaceable)
	_ = Global().Add(fileSink{})

	entries := Global().GetByInterface(reflect.TypeFor[writePort]())
	if len(entries) != 2 {
		t.Fatalf("expected 2 implementors, got %d", len(entries))
	}

	shapes := Global().GetByInterface(reflect.TypeFor[Shaper]())
	_ = Global().Add(&Square{Side: 1})
	shapes = Global().GetByInterface(reflect.TypeFor[Shaper]())
	if len(shapes) != 1 {
		t.Errorf("GetByInterface[Shaper] expected 1, got %d", len(shapes))
	}
}

func TestTransform_noop(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().Add(&Square{Side: 5})

	err := Global().Transform(func(r any) any { return r })
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	entries := Global().GetByType(reflect.TypeFor[*Square]())
	if len(entries) != 1 || entries[0].Value.(*Square).Side != 5 {
		t.Errorf("expected square after noop transform, got %v", entries)
	}
}

func TestTransform_preservesTags(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().AddWithTags(&Square{Side: 1}, TagReplaceable)

	if err := Global().Transform(func(r any) any {
		if sq, ok := r.(*Square); ok {
			return &Square{Side: sq.Side + 1}
		}
		return r
	}); err != nil {
		t.Fatal(err)
	}

	var replaceable bool
	Global().WalkEntries(func(e Entry) bool {
		replaceable = e.Replaceable()
		return true
	})
	if !replaceable {
		t.Fatalf("TagReplaceable after transform: got false want true")
	}
}

func TestAddWithTags_fixed(t *testing.T) {
	resetGlobalRegistry()
	if err := Global().AddWithTags(&Square{Side: 1}, TagFixed); err != nil {
		t.Fatal(err)
	}

	var fixed bool
	Global().WalkEntries(func(e Entry) bool {
		fixed = e.Fixed()
		return true
	})
	if !fixed {
		t.Fatal("expected TagFixed on entry")
	}
}

func TestTransform_preservesFixedTag(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().AddWithTags(&Square{Side: 1}, TagFixed)

	if err := Global().Transform(func(r any) any {
		if sq, ok := r.(*Square); ok {
			return &Square{Side: sq.Side + 1}
		}
		return r
	}); err != nil {
		t.Fatal(err)
	}

	var fixed bool
	Global().WalkEntries(func(e Entry) bool {
		fixed = e.Fixed()
		return true
	})
	if !fixed {
		t.Fatalf("TagFixed after transform: got false want true")
	}
}

func TestTransform_empty(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().Add(&Square{Side: 3})

	if err := Global().Transform(); err != nil {
		t.Fatalf("empty Transform failed: %v", err)
	}

	entries := Global().GetByType(reflect.TypeFor[*Square]())
	if len(entries) != 1 || entries[0].Value.(*Square).Side != 3 {
		t.Errorf("empty Transform must not change resources, got %v", entries)
	}
}

func TestTransform_updatesSliceInPlace(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().Add(&Square{Side: 5})

	before := Global().GetByType(reflect.TypeFor[*Square]())[0].Value.(*Square)

	err := Global().Transform(func(r any) any {
		if sq, ok := r.(*Square); ok {
			return &Square{Side: sq.Side + 10}
		}
		return r
	})
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	after := Global().GetByType(reflect.TypeFor[*Square]())[0].Value.(*Square)
	if after == before {
		t.Fatal("registry must hold transformed resource instance")
	}
	if after.Side != 15 {
		t.Fatalf("transformed resource expected Side=15, got %v", after.Side)
	}
}

func TestTransform_typeChange(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().Add(&Square{Side: 10})

	err := Global().Transform(func(r any) any {
		if sq, ok := r.(*Square); ok {
			return &wrappedSquare{Square: sq}
		}
		return r
	})
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if len(Global().GetByType(reflect.TypeFor[*Square]())) != 0 {
		t.Error("GetByType[*Square] should be empty after wrap")
	}

	shapes := Global().GetByInterface(reflect.TypeFor[Shaper]())
	if len(shapes) != 1 {
		t.Errorf("GetByInterface[Shaper] expected 1, got %d", len(shapes))
	}
}

type Circle struct {
	Radius int
}

func (c *Circle) Area() int { return c.Radius * c.Radius }

func TestTransform_duplicateType(t *testing.T) {
	resetGlobalRegistry()
	_ = Global().Add(&Square{Side: 1})
	_ = Global().Add(&Circle{Radius: 2})

	err := Global().Transform(func(r any) any {
		switch v := r.(type) {
		case *Square:
			return &wrappedSquare{Square: v}
		case *Circle:
			return &wrappedSquare{Square: &Square{Side: v.Radius}}
		}
		return r
	})
	if err != nil {
		t.Fatalf("Transform with duplicate resulting type should succeed: %v", err)
	}

	if len(Global().GetByType(reflect.TypeFor[*wrappedSquare]())) != 2 {
		t.Fatalf("expected 2 wrappedSquare entries, got %d", len(Global().GetByType(reflect.TypeFor[*wrappedSquare]())))
	}
}

func TestRegistry_implementsRegistry(t *testing.T) {
	var _ Registry = Global()
}
