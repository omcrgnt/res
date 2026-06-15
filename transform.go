package res

// TransformFunc transforms a single resource during [Registry.Transform].
type TransformFunc func(any) any

func (r *registry) Transform(fns ...TransformFunc) error {
	if len(fns) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i, it := range r.items {
		res := it.value
		for _, fn := range fns {
			res = fn(res)
		}
		it.value = res
		r.items[i] = it
	}

	return nil
}
