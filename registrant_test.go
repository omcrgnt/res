package res

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

type stubNew struct{}

func (stubNew) NewResource() (any, error) { return "built", nil }

type stubBuild struct{}

func (stubBuild) BuildConfig() (BuildSpec, error) {
	return stubSpec{}, nil
}

type stubSpec struct{}

func (stubSpec) Build() (any, error) { return "built", nil }

type stubDual struct{}

func (stubDual) NewResource() (any, error) { return nil, nil }
func (stubDual) BuildConfig() (BuildSpec, error) {
	return stubSpec{}, nil
}

type stubBare struct{}

func TestValidateRegistrant_rejectsInvalid(t *testing.T) {
	cases := []struct {
		name string
		v    any
	}{
		{"string", "x"},
		{"bare struct", stubBare{}},
		{"builder only", stubBuilderOnly{}},
		{"dual", stubDual{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := New().Add(tc.v)
			if !errors.Is(err, ErrInvalidRegistrant) {
				t.Fatalf("Add() = %v, want ErrInvalidRegistrant", err)
			}
		})
	}
}

type stubBuilderOnly struct{}

func (stubBuilderOnly) Build() (any, error) { return nil, nil }

var _ BuildSpec = stubBuilderOnly{}

func TestValidateRegistrant_acceptsWire(t *testing.T) {
	if err := New().Add(stubNew{}); err != nil {
		t.Fatal(err)
	}
	if err := New().Add(stubBuild{}); err != nil {
		t.Fatal(err)
	}
}

func TestMustAddToGlobalWithTags_panicsOnInvalid(t *testing.T) {
	resetGlobalRegistry()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		} else if s, ok := r.(string); !ok || !strings.Contains(s, "MustAddToGlobalWithTags") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	MustAddToGlobalWithTags("invalid", TagReplaceable)
}

func TestMustAddToGlobalWithTags_ok(t *testing.T) {
	resetGlobalRegistry()
	MustAddToGlobalWithTags(stubNew{}, TagReplaceable)
	if len(Global().GetByType(reflect.TypeFor[stubNew]())) != 1 {
		t.Fatal("expected wire on global")
	}
}
