package rules

import "sync"

// registry holds all registered rules.
var (
	mu       sync.RWMutex
	registry []Rule
)

// Register adds a rule to the global registry. Called by rule init() functions.
func Register(r Rule) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, r)
}

// All returns all registered rules.
func All() []Rule {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Rule, len(registry))
	copy(out, registry)
	return out
}

// ForPlatform returns all rules applicable to the given platform.
func ForPlatform(p Platform) []Rule {
	mu.RLock()
	defer mu.RUnlock()
	var out []Rule
	for _, r := range registry {
		for _, rp := range r.Platforms() {
			if rp == p {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

// ByID returns a specific rule by its ID, or nil if not found.
func ByID(id string) Rule {
	mu.RLock()
	defer mu.RUnlock()
	for _, r := range registry {
		if r.ID() == id {
			return r
		}
	}
	return nil
}
