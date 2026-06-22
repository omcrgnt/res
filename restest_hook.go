package res

// ResetGlobalForRestest replaces [Global] with an empty registry.
//
//go:internal github.com/omcrgnt/res/restest
func ResetGlobalForRestest() Registry {
	resetGlobalRegistry()
	return global
}
