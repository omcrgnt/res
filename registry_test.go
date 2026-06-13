package res

import (
	"reflect"
	"testing"
)

func TestRegistry_AddGet(t *testing.T) {
	r := newRegistry()
	val := "hello"

	if err := addUser(r, val); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	res, ok := get[string](r)
	if !ok || res != "hello" {
		t.Errorf("Get failed: expected 'hello', got %v", res)
	}

	if err := addUser(r, "world"); err == nil {
		t.Error("expected error on duplicate type registration")
	}
}

func TestAdd_nil(t *testing.T) {
	r := newRegistry()
	if err := addUser(r, nil); err == nil {
		t.Fatal("expected error for nil resource")
	}
	if err := addBuiltin(r, nil); err == nil {
		t.Fatal("expected error for nil builtin")
	}
}

func TestAddAll_returnsErrorOnDuplicate(t *testing.T) {
	resetGlobalRegistry()

	err := AddAll("a", "b")
	if err == nil {
		t.Fatal("expected duplicate type error")
	}
	if countGlobal() != 1 {
		t.Fatalf("expected 1 resource after partial AddAll, got %d", countGlobal())
	}
	val, ok := Get[string]()
	if !ok || val != "a" {
		t.Fatalf("expected first resource retained, got %v ok=%v", val, ok)
	}
}

func TestAddAll_returnsErrorOnNil(t *testing.T) {
	resetGlobalRegistry()

	err := AddAll(nil, "ok")
	if err == nil {
		t.Fatal("expected nil resource error")
	}
	if _, ok := Get[string](); ok {
		t.Error("resource must not be added after nil in AddAll")
	}
}

func TestRegistry_AddAll(t *testing.T) {
	resetGlobalRegistry()

	err := AddAll(10, "string")
	if err != nil {
		t.Fatalf("AddAll failed: %v", err)
	}

	if countGlobal() != 2 {
		t.Fatalf("expected 2 resources, got %d", countGlobal())
	}
}

func TestDefault_Add(t *testing.T) {
	resetGlobalRegistry()

	if err := Default.Add("via-default"); err != nil {
		t.Fatalf("Default.Add failed: %v", err)
	}

	got, ok := Get[string]()
	if !ok || got != "via-default" {
		t.Errorf("unexpected resource: %v ok=%v", got, ok)
	}
}

func TestAddBuiltin(t *testing.T) {
	resetGlobalRegistry()

	if err := AddBuiltin("system"); err != nil {
		t.Fatal(err)
	}
	if err := AddBuiltin("again"); err == nil {
		t.Fatal("expected duplicate system error")
	}

	var origin Origin
	WalkEntries(func(e Entry) bool {
		origin = e.Origin
		return true
	})
	if origin != System {
		t.Fatalf("origin: got %v want System", origin)
	}
}

func TestAdd_replacesSystem(t *testing.T) {
	resetGlobalRegistry()

	type widget struct{ n int }

	if err := AddBuiltin(&widget{n: 1}); err != nil {
		t.Fatal(err)
	}
	if err := Add(&widget{n: 2}); err != nil {
		t.Fatal(err)
	}

	got, ok := Get[*widget]()
	if !ok || got.n != 2 {
		t.Fatalf("got %+v ok=%v", got, ok)
	}

	var origin Origin
	WalkEntries(func(e Entry) bool {
		origin = e.Origin
		return true
	})
	if origin != User {
		t.Fatalf("origin: got %v want User", origin)
	}
	if countGlobal() != 1 {
		t.Fatalf("expected 1 resource after replace, got %d", countGlobal())
	}
}

func TestAddBuiltin_afterUser(t *testing.T) {
	resetGlobalRegistry()

	if err := Add("user"); err != nil {
		t.Fatal(err)
	}
	if err := AddBuiltin("system"); err == nil {
		t.Fatal("expected error when AddBuiltin after user same type")
	}
}

func TestDefault_Walk(t *testing.T) {
	resetGlobalRegistry()
	_ = AddAll("a", 2)

	var seen []any
	Default.Walk(func(t reflect.Type, res any) bool {
		seen = append(seen, res)
		return true
	})

	if len(seen) != 2 {
		t.Fatalf("Default.Walk expected 2, got %d", len(seen))
	}
}

func TestRegistry_Walk(t *testing.T) {
	r := newRegistry()
	_ = addUser(r, 10)
	_ = addUser(r, "string")

	count := 0
	r.walk(func(t reflect.Type, res any) bool {
		count++
		return true
	})

	if count != 2 {
		t.Errorf("Walk failed: expected 2 items, got %d", count)
	}
}

func TestRegistry_Walk_stopsOnFalse(t *testing.T) {
	r := newRegistry()
	_ = addUser(r, "one")
	_ = addUser(r, 2)
	_ = addUser(r, true)

	var seen []any
	r.walk(func(t reflect.Type, res any) bool {
		seen = append(seen, res)
		return res != 2
	})

	if len(seen) != 2 {
		t.Fatalf("walk must stop after fn returns false, got %v", seen)
	}
	if seen[0] != "one" || seen[1] != 2 {
		t.Fatalf("unexpected visit order: %v", seen)
	}
}

func TestWalk_global(t *testing.T) {
	resetGlobalRegistry()
	_ = AddAll("a", 2)

	count := 0
	Walk(func(t reflect.Type, res any) bool {
		count++
		return true
	})

	if count != 2 {
		t.Fatalf("global Walk expected 2 visits, got %d", count)
	}
}

func TestInstance_isDefault(t *testing.T) {
	resetGlobalRegistry()
	_ = Default.Add(42)

	got, ok := Get[int]()
	if !ok || got != 42 {
		t.Fatalf("Instance alias failed: got %v ok=%v", got, ok)
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

func TestFind(t *testing.T) {
	resetGlobalRegistry()

	sq := &Square{Side: 10}
	_ = Add(sq)

	squares := Find[*Square]()
	if len(squares) != 1 {
		t.Errorf("Find[*Square] failed: expected 1 match, got %d", len(squares))
	}

	shapes := Find[Shaper]()
	if len(shapes) != 1 {
		t.Errorf("Find[Shaper] failed: expected 1 match, got %d", len(shapes))
	}

	strings := Find[string]()
	if len(strings) != 0 {
		t.Errorf("Find[string] failed: expected 0 matches, got %d", len(strings))
	}
}

type writePort interface {
	Write([]byte) (int, error)
}

type stdoutSink struct{}

func (stdoutSink) Write(p []byte) (int, error) { return len(p), nil }

type fileSink struct{}

func (fileSink) Write(p []byte) (int, error) { return len(p), nil }

func TestAddBuiltin_and_Add_sameInterface(t *testing.T) {
	resetGlobalRegistry()

	if err := AddBuiltin(stdoutSink{}); err != nil {
		t.Fatal(err)
	}
	if err := Add(fileSink{}); err != nil {
		t.Fatal(err)
	}

	if n := len(Find[writePort]()); n != 2 {
		t.Fatalf("expected 2 implementors, got %d", n)
	}
}

func TestTransform_noop(t *testing.T) {
	resetGlobalRegistry()
	_ = AddAll(&Square{Side: 5})

	err := Transform(func(r any) any { return r })
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	sq, ok := Get[*Square]()
	if !ok || sq.Side != 5 {
		t.Errorf("expected square after noop transform, got %v ok=%v", sq, ok)
	}
}

func TestTransform_preservesOrigin(t *testing.T) {
	resetGlobalRegistry()
	_ = AddBuiltin(&Square{Side: 1})

	if err := Transform(func(r any) any {
		if sq, ok := r.(*Square); ok {
			return &Square{Side: sq.Side + 1}
		}
		return r
	}); err != nil {
		t.Fatal(err)
	}

	var origin Origin
	WalkEntries(func(e Entry) bool {
		origin = e.Origin
		return true
	})
	if origin != System {
		t.Fatalf("origin after transform: got %v want System", origin)
	}
}

func TestTransform_empty(t *testing.T) {
	resetGlobalRegistry()
	_ = AddAll(&Square{Side: 3})

	if err := Transform(); err != nil {
		t.Fatalf("empty Transform failed: %v", err)
	}

	sq, ok := Get[*Square]()
	if !ok || sq.Side != 3 {
		t.Errorf("empty Transform must not change resources, got %v ok=%v", sq, ok)
	}
}

func TestTransform_updatesSliceInPlace(t *testing.T) {
	resetGlobalRegistry()
	_ = Add(&Square{Side: 5})

	before, ok := Get[*Square]()
	if !ok {
		t.Fatal("expected Square before transform")
	}

	err := Transform(func(r any) any {
		if sq, ok := r.(*Square); ok {
			return &Square{Side: sq.Side + 10}
		}
		return r
	})
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	after, ok := Get[*Square]()
	if !ok {
		t.Fatal("expected Square after transform")
	}
	if after == before {
		t.Fatal("registry must hold transformed resource instance")
	}
	if after.Side != 15 {
		t.Fatalf("transformed resource expected Side=15, got %v", after.Side)
	}
}

func countGlobal() int {
	n := 0
	Walk(func(reflect.Type, any) bool {
		n++
		return true
	})
	return n
}

func TestTransform_typeChange(t *testing.T) {
	resetGlobalRegistry()
	_ = Add(&Square{Side: 10})

	err := Transform(func(r any) any {
		if sq, ok := r.(*Square); ok {
			return &wrappedSquare{Square: sq}
		}
		return r
	})
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if _, ok := Get[*Square](); ok {
		t.Error("Get[*Square] should fail after wrap")
	}

	shapes := Find[Shaper]()
	if len(shapes) != 1 {
		t.Errorf("Find[Shaper] expected 1, got %d", len(shapes))
	}
}

type Circle struct {
	Radius int
}

func (c *Circle) Area() int { return c.Radius * c.Radius }

func TestTransform_duplicateType(t *testing.T) {
	resetGlobalRegistry()
	_ = AddAll(&Square{Side: 1}, &Circle{Radius: 2})

	err := Transform(func(r any) any {
		switch v := r.(type) {
		case *Square:
			return &wrappedSquare{Square: v}
		case *Circle:
			return &wrappedSquare{Square: &Square{Side: v.Radius}}
		}
		return r
	})
	if err == nil {
		t.Fatal("expected duplicate type error")
	}
}
