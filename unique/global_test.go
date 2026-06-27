package unique

import (
	"reflect"
	"testing"
)

type fixedWidget struct{ n int }

func TestMustAddReplaceable(t *testing.T) {
	g := Global()
	w := &testWidget{n: 1}

	MustAddReplaceable(w)

	entries := g.GetByType(reflect.TypeFor[*testWidget]())
	if len(entries) != 1 || !entries[0].Replaceable() || entries[0].Value() != w {
		t.Fatalf("entry = %+v", entries)
	}
}

func TestMustAddFixed(t *testing.T) {
	g := Global()
	w := &fixedWidget{n: 2}

	MustAddFixed(w)

	entries := g.GetByType(reflect.TypeFor[*fixedWidget]())
	if len(entries) != 1 || !entries[0].Fixed() || entries[0].Value() != w {
		t.Fatalf("entry = %+v", entries)
	}
}
