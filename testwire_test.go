package res

// wireValue is a NewResourceer test helper that materializes to v.
type wireValue struct{ v any }

func (w wireValue) NewResource() (any, error) { return w.v, nil }

func addWire(reg Registry, v any) error {
	return reg.Add(wireValue{v: v})
}

func addBuilt(r *registry, v any) error {
	w := wireValue{v: v}
	if err := r.Add(w); err != nil {
		return err
	}
	return r.replaceValue(w, v)
}
