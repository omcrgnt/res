package respolicy_test

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/respolicy"
)

type stub struct{ n int }

type prefixPolicy struct{ prefix string }

func (p prefixPolicy) PrepareAdd(v any) (any, error) {
	s, ok := v.(stub)
	if !ok {
		return v, nil
	}
	return stub{n: len(p.prefix) + s.n}, nil
}

func TestWrap_Add_delegates(t *testing.T) {
	reg := respolicy.Wrap(res.New(), respolicy.AcceptAll{})
	want := stub{n: 1}
	if err := reg.Add(want); err != nil {
		t.Fatal(err)
	}
	got, err := reg.GetOneByType(reflect.TypeFor[stub]())
	if err != nil {
		t.Fatal(err)
	}
	if got.(stub).n != 1 {
		t.Fatalf("got %#v", got)
	}
}

func TestWrap_AddPolicy_transforms(t *testing.T) {
	reg := respolicy.Wrap(res.New(), prefixPolicy{prefix: "xx"})
	if err := reg.Add(stub{n: 3}); err != nil {
		t.Fatal(err)
	}
	got, err := reg.GetOneByType(reflect.TypeFor[stub]())
	if err != nil {
		t.Fatal(err)
	}
	if got.(stub).n != 5 {
		t.Fatalf("got n=%d, want 5", got.(stub).n)
	}
}

func TestWrap_AddWithTags(t *testing.T) {
	reg := respolicy.Wrap(res.New(), respolicy.AcceptAll{})
	if err := reg.AddWithTags(stub{n: 1}, res.TagReplaceable); err != nil {
		t.Fatal(err)
	}
	var replaceable bool
	reg.WalkEntries(func(e res.Entry) bool {
		replaceable = e.Replaceable()
		return true
	})
	if !replaceable {
		t.Fatal("expected TagReplaceable on entry")
	}
}

type recordPolicy struct {
	order *[]string
	name  string
}

func (p recordPolicy) PrepareAdd(v any) (any, error) {
	*p.order = append(*p.order, p.name)
	return v, nil
}

func TestWrap_chain(t *testing.T) {
	var order []string
	reg := respolicy.Wrap(res.New(),
		recordPolicy{order: &order, name: "first"},
		recordPolicy{order: &order, name: "second"},
	)
	if err := reg.Add(stub{}); err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Fatalf("order=%v", order)
	}
}

func TestRejectInvalid_nil(t *testing.T) {
	reg := respolicy.Wrap(res.New(), respolicy.RejectInvalid{})
	if err := reg.Add(nil); err == nil {
		t.Fatal("expected error for nil")
	}
}
