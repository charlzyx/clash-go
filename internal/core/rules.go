package core

import (
	"fmt"
	"sort"
	"time"

	"github.com/metacubex/mihomo/tunnel"
)

// Rule is one compiled rule from the kernel. Index is not stored: the
// rules page derives the line number as index+1 and does not write it back.
type Rule struct {
	Type    string
	Payload string
	Proxy   string
}

// RuleProvider is one rule-provider the kernel has loaded.
type RuleProvider struct {
	Name        string
	VehicleType string
	Behavior    string
	RuleCount   int
	UpdatedAt   time.Time
}

// Rules returns the kernel's rule list. The slice is nil when the core
// is stopped. The order is the kernel's order.
func (c *Core) Rules() []Rule {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return nil
	}
	raw := tunnel.Rules()
	out := make([]Rule, 0, len(raw))
	for _, rule := range raw {
		if rule == nil {
			continue
		}
		out = append(out, Rule{
			Type:    rule.RuleType().String(),
			Payload: rule.Payload(),
			Proxy:   rule.Adapter(),
		})
	}
	return out
}

// RuleProviders returns the loaded rule providers, sorted by name.
// The slice is nil when the core is stopped.
func (c *Core) RuleProviders() []RuleProvider {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return nil
	}
	raw := tunnel.RuleProviders()
	out := make([]RuleProvider, 0, len(raw))
	for name, p := range raw {
		if p == nil {
			continue
		}
		item := RuleProvider{
			Name:        p.Name(),
			VehicleType: p.VehicleType().String(),
			Behavior:    p.Behavior().String(),
			RuleCount:   p.Count(),
		}
		if item.Name == "" {
			item.Name = name
		}
		if stamped, ok := p.(interface{ UpdatedAt() time.Time }); ok {
			item.UpdatedAt = stamped.UpdatedAt()
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// UpdateRuleProvider asks one rule provider to reload. It does not hold
// the core lock across the provider's own update, which may read a file.
func (c *Core) UpdateRuleProvider(name string) error {
	if c == nil {
		return fmt.Errorf("core is nil")
	}
	if name == "" {
		return fmt.Errorf("rule provider name is empty")
	}
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return fmt.Errorf("core is not running")
	}
	p := tunnel.RuleProviders()[name]
	c.mu.Unlock()
	if p == nil {
		return fmt.Errorf("rule provider %q not found", name)
	}
	return p.Update()
}
