package catalog

// ModuleOption represents a user-configurable option for a module
type ModuleOption struct {
	Key      string
	Question string
	Choices  []string
	Default  string
}

// Module represents one go-tools library
type Module struct {
	Name        string
	ImportPath  string
	Description string
	Version     string
	Tags        []string
	Options     []ModuleOption
}

// Preset bundles commonly used modules
type Preset struct {
	Name        string
	Description string
	Modules     []string
}

var Modules = []Module{
	{Name: "go-foundation", ImportPath: "github.com/mirkobrombin/go-foundation", Description: "Shared primitives: DI, tags, resiliency", Tags: []string{"core"}},
	{Name: "go-logger", ImportPath: "github.com/mirkobrombin/go-logger", Description: "Structured logging with multiple sinks", Tags: []string{"observability"}},
	{Name: "go-metrics", ImportPath: "github.com/mirkobrombin/go-metrics", Description: "Counter abstraction for application metrics", Tags: []string{"observability"}},
	{Name: "go-cli-builder/v2", ImportPath: "github.com/mirkobrombin/go-cli-builder/v2", Description: "Declarative struct-tag-driven CLI builder", Version: "v2", Tags: []string{"cli"}},
	{Name: "go-conf-builder/v2", ImportPath: "github.com/mirkobrombin/go-conf-builder/v2", Description: "Multi-source config loader (env, file, flags)", Version: "v2", Tags: []string{"cli"}},
	{Name: "go-struct-flags/v2", ImportPath: "github.com/mirkobrombin/go-struct-flags/v2", Description: "Flag-to-struct binding via struct tags", Version: "v2", Tags: []string{"cli"}},
	{Name: "go-module-router/v2", ImportPath: "github.com/mirkobrombin/go-module-router/v2", Description: "Transport-agnostic router with DI", Version: "v2", Tags: []string{"http"}},
	{Name: "go-httpx", ImportPath: "github.com/mirkobrombin/go-httpx", Description: "Middleware-enhanced HTTP client", Tags: []string{"http"}},
	{Name: "go-auth", ImportPath: "github.com/mirkobrombin/go-auth", Description: "HMAC-SHA256 token signing and verification", Tags: []string{"auth"}},
	{Name: "go-guard", ImportPath: "github.com/mirkobrombin/go-guard", Description: "Declarative attribute-based access control (ABAC)", Tags: []string{"auth"}},
	{Name: "go-secrets", ImportPath: "github.com/mirkobrombin/go-secrets", Description: "Secret store abstraction", Tags: []string{"auth"}},
	{Name: "go-signal/v2", ImportPath: "github.com/mirkobrombin/go-signal/v2", Description: "Type-safe in-process event bus", Version: "v2", Tags: []string{"async"}},
	{
		Name: "go-relay/v2", ImportPath: "github.com/mirkobrombin/go-relay/v2", Description: "Durable async job processing (broker-agnostic)", Version: "v2", Tags: []string{"async"},
		Options: []ModuleOption{
			{Key: "broker", Question: "Message broker backend", Choices: []string{"Memory", "Redis", "NATS", "Warp"}, Default: "Memory"},
		},
	},
	{Name: "go-worker", ImportPath: "github.com/mirkobrombin/go-worker", Description: "Fixed-size worker pool", Tags: []string{"async"}},
	{Name: "go-retry", ImportPath: "github.com/mirkobrombin/go-retry", Description: "Compact retry helper with backoff", Tags: []string{"async"}},
	{Name: "go-revert/v2", ImportPath: "github.com/mirkobrombin/go-revert/v2", Description: "Saga / compensation workflows (panic-safe rollback)", Version: "v2", Tags: []string{"async"}},
	{Name: "go-state-flow", ImportPath: "github.com/mirkobrombin/go-state-flow", Description: "Declarative finite-state machine", Tags: []string{"state"}},
	{Name: "go-lock", ImportPath: "github.com/mirkobrombin/go-lock", Description: "Distributed locking abstraction", Tags: []string{"infra"}},
	{Name: "go-warp", ImportPath: "github.com/mirkobrombin/go-warp", Description: "L1/L2 cache with distributed sync bus", Tags: []string{"infra"}},
	{Name: "go-slipstream", ImportPath: "github.com/mirkobrombin/go-slipstream", Description: "Embedded Bitcask+Raft database", Tags: []string{"infra"}},
	{Name: "go-plugin", ImportPath: "github.com/mirkobrombin/go-plugin", Description: "Plugin registry with lifecycle management", Tags: []string{"plugin"}},
	{
		Name: "go-wormhole", ImportPath: "github.com/fabricatorsltd/go-wormhole", Description: "EF-style ORM with code-first migrations and Unit of Work", Tags: []string{"data"},
		Options: []ModuleOption{
			{Key: "provider", Question: "Database provider", Choices: []string{"SQLite", "PostgreSQL", "MySQL", "MSSQL"}, Default: "SQLite"},
		},
	},
}

var Presets = []Preset{
	{
		Name:        "API server",
		Description: "HTTP API with routing, auth, and structured logging",
		Modules:     []string{"go-foundation", "go-logger", "go-conf-builder/v2", "go-module-router/v2", "go-auth", "go-guard", "go-httpx"},
	},
	{
		Name:        "Worker service",
		Description: "Background job processing with async messaging",
		Modules:     []string{"go-foundation", "go-logger", "go-conf-builder/v2", "go-relay/v2", "go-worker", "go-signal/v2", "go-retry"},
	},
	{
		Name:        "Full-stack",
		Description: "All 22 go-tools modules",
		Modules: func() []string {
			names := make([]string, len(Modules))
			for i, m := range Modules {
				names[i] = m.Name
			}
			return names
		}(),
	},
	{
		Name:        "Minimal",
		Description: "Foundation, logger, config — clean starting point",
		Modules:     []string{"go-foundation", "go-logger", "go-conf-builder/v2"},
	},
}

// FindModule returns the Module with the given name, or nil
func FindModule(name string) *Module {
	for i := range Modules {
		if Modules[i].Name == name {
			return &Modules[i]
		}
	}
	return nil
}

// GroupByTag returns modules grouped by their first tag
func GroupByTag() map[string][]Module {
	groups := make(map[string][]Module)
	for _, m := range Modules {
		if len(m.Tags) > 0 {
			tag := m.Tags[0]
			groups[tag] = append(groups[tag], m)
		}
	}
	return groups
}
