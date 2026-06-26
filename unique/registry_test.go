package unique

import (
	"errors"
	"reflect"
	"testing"

	"github.com/omcrgnt/res"
)

type testWidget struct{ n int }

func TestAdd_empty(t *testing.T) {
	reg := New()
	w := &testWidget{n: 1}

	if err := reg.Add(w); err != nil {
		t.Fatal(err)
	}

	entries := reg.GetByType(reflect.TypeFor[*testWidget]())
	if len(entries) != 1 || !entries[0].Regular() || entries[0].Value != w {
		t.Fatalf("entry = %+v", entries)
	}
}

func TestAdd_overReplaceable(t *testing.T) {
	reg := New()
	def := &testWidget{n: 1}
	app := &testWidget{n: 2}

	if err := reg.AddReplaceable(def); err != nil {
		t.Fatal(err)
	}
	if err := reg.Add(app); err != nil {
		t.Fatal(err)
	}

	entries := reg.GetByType(reflect.TypeFor[*testWidget]())
	if len(entries) != 1 || !entries[0].Regular() || entries[0].Value != app {
		t.Fatalf("entry = %+v", entries)
	}
}

func TestAdd_fixedConflict(t *testing.T) {
	reg := New()
	if err := reg.AddFixed(&testWidget{n: 1}); err != nil {
		t.Fatal(err)
	}
	if err := reg.Add(&testWidget{n: 2}); !errors.Is(err, ErrFixed) {
		t.Fatalf("Add = %v, want ErrFixed", err)
	}
}

func TestAdd_regularExists(t *testing.T) {
	reg := New()
	if err := reg.Add(&testWidget{n: 1}); err != nil {
		t.Fatal(err)
	}
	if err := reg.Add(&testWidget{n: 2}); !errors.Is(err, ErrRegularExists) {
		t.Fatalf("Add = %v, want ErrRegularExists", err)
	}
}

func TestAdd_duplicateType(t *testing.T) {
	reg := newRegistry(res.New())
	// Bypass unique write rules to simulate broken bag.
	_ = reg.reg.Add(&testWidget{n: 1})
	_ = reg.reg.Add(&testWidget{n: 2})

	if err := reg.Add(&testWidget{n: 3}); !errors.Is(err, ErrDuplicateType) {
		t.Fatalf("Add = %v, want ErrDuplicateType", err)
	}
}

func TestAddReplaceable_and_AddFixed(t *testing.T) {
	reg := New()

	if err := reg.AddReplaceable(&testWidget{n: 1}); err != nil {
		t.Fatal(err)
	}
	if err := reg.Add(&testWidget{n: 2}); err != nil {
		t.Fatalf("override replaceable: %v", err)
	}
	if err := reg.AddReplaceable(&testWidget{n: 3}); !errors.Is(err, ErrTypeOccupied) {
		t.Fatalf("AddReplaceable on regular = %v", err)
	}

	reg2 := New()
	if err := reg2.AddFixed(&testWidget{n: 1}); err != nil {
		t.Fatal(err)
	}
	if err := reg2.AddReplaceable(&testWidget{n: 2}); !errors.Is(err, ErrTypeOccupied) {
		t.Fatalf("AddReplaceable on fixed = %v", err)
	}
	if err := reg2.AddFixed(&testWidget{n: 3}); !errors.Is(err, ErrTypeOccupied) {
		t.Fatalf("AddFixed on fixed = %v", err)
	}
}

func TestMerge(t *testing.T) {
	dst := New()
	src := New()

	if err := dst.AddReplaceable(&testWidget{n: 1}); err != nil {
		t.Fatal(err)
	}
	if err := src.Add(&testWidget{n: 2}); err != nil {
		t.Fatal(err)
	}
	if err := Merge(dst, src); err != nil {
		t.Fatal(err)
	}

	entries := dst.GetByType(reflect.TypeFor[*testWidget]())
	if len(entries) != 1 || !entries[0].Regular() || entries[0].Value.(*testWidget).n != 2 {
		t.Fatalf("after merge: %+v", entries)
	}
}

func TestGlobal_isolatedFromNew(t *testing.T) {
	g := Global()
	n := New()

	w1 := &testWidget{n: 1}
	w2 := &testWidget{n: 2}
	if err := g.Add(w1); err != nil {
		t.Fatal(err)
	}
	if err := n.Add(w2); err != nil {
		t.Fatal(err)
	}

	gotG, _ := g.GetOneByType(reflect.TypeFor[*testWidget]())
	gotN, _ := n.GetOneByType(reflect.TypeFor[*testWidget]())
	if gotG != w1 || gotN != w2 {
		t.Fatalf("Global=%p New=%p", gotG, gotN)
	}
}

func TestMerge_nilAndSame(t *testing.T) {
	reg := New()
	if err := Merge(nil, reg); !errors.Is(err, errNilRegistry) {
		t.Fatalf("nil dst = %v", err)
	}
	if err := Merge(reg, nil); !errors.Is(err, errNilRegistry) {
		t.Fatalf("nil src = %v", err)
	}
	if err := Merge(reg, reg); !errors.Is(err, errSameRegistry) {
		t.Fatalf("same reg = %v", err)
	}
}

func TestAdd_preservesOrderOnReplace(t *testing.T) {
	reg := New()
	other := "other"
	if err := reg.AddReplaceable(&testWidget{n: 0}); err != nil {
		t.Fatal(err)
	}
	if err := reg.reg.Add(other); err != nil {
		t.Fatal(err)
	}
	if err := reg.Add(&testWidget{n: 1}); err != nil {
		t.Fatal(err)
	}

	var order []any
	reg.WalkEntries(func(e res.Entry) bool {
		order = append(order, e.Value)
		return true
	})
	if len(order) != 2 {
		t.Fatalf("order len = %d", len(order))
	}
	if _, ok := order[0].(*testWidget); !ok || order[1] != other {
		t.Fatalf("order = %v", order)
	}
}
