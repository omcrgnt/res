package res

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/sdi"
)

// Dedup runs policy for each interface port before sdi.Resolve inject phase.
func (globalAccessor) Dedup(interfaces []reflect.Type, policy sdi.DedupPolicy) error {
	return globalRegistry.dedup(interfaces, policy)
}

func Dedup(interfaces []reflect.Type, policy sdi.DedupPolicy) error {
	return Default.Dedup(interfaces, policy)
}

func (r *registry) dedup(interfaces []reflect.Type, policy sdi.DedupPolicy) error {
	for _, iface := range interfaces {
		entries := r.listImplementors(iface)
		dedupEntries := make([]sdi.DedupEntry, len(entries))
		for i, e := range entries {
			dedupEntries[i] = sdi.DedupEntry{
				Value:     e.value,
				Removable: e.origin == System,
			}
		}
		ctx := sdi.DedupContext{
			Interface: iface,
			Entries:   dedupEntries,
			Remove: func(v any) error {
				return r.removeSystem(v)
			},
		}
		if err := policy(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *registry) listImplementors(iface reflect.Type) []entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []entry
	for _, e := range r.resources {
		if reflect.TypeOf(e.value).Implements(iface) {
			matches = append(matches, e)
		}
	}
	return matches
}

func (r *registry) removeSystem(v any) error {
	if v == nil {
		return fmt.Errorf("cannot remove nil resource")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	idx := -1
	var key reflect.Type
	for i, e := range r.resources {
		if e.value == v {
			if e.origin != System {
				return fmt.Errorf("cannot remove user resource")
			}
			idx = i
			key = reflect.TypeOf(v)
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("resource not found")
	}

	delete(r.byType, key)
	r.resources = append(r.resources[:idx], r.resources[idx+1:]...)
	return nil
}
