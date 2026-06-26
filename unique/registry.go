package unique

import "github.com/omcrgnt/res"

type Registry struct {
	reg res.Registry
}

func newRegistry(reg res.Registry) *Registry {
	return &Registry{reg: reg}
}
