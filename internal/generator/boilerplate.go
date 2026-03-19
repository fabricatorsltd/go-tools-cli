package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fabricatorsltd/go-tools-cli/internal/catalog"
	"github.com/fabricatorsltd/go-tools-cli/internal/wizard"
)

func generateBoilerplate(cfg wizard.Config) ([]string, error) {
	var files []string

	internalDir := filepath.Join(cfg.OutputDir, "internal")

	hasModule := func(name string) bool {
		for _, s := range cfg.SelectedModules {
			if s == name {
				return true
			}
		}
		return false
	}

	type stub struct {
		moduleName string
		dir        string
		filename   string
		content    func() string
	}

	stubs := []stub{
		{
			moduleName: "go-conf-builder/v2",
			dir:        "config",
			filename:   "config.go",
			content: func() string {
				return fmt.Sprintf(`package config

import (
	confbuilder "github.com/mirkobrombin/go-conf-builder/v2/pkg/conf"
)

// Config holds all application configuration.
type Config struct {
	Host string `+"`"+`conf:"host" default:"0.0.0.0"`+"`"+`
	Port int    `+"`"+`conf:"port" default:"8080"`+"`"+`
}

// Load loads configuration from all sources (env, file, flags).
func Load() (*Config, error) {
	cfg := &Config{}
	loader := confbuilder.New(cfg)
	if err := loader.Load(); err != nil {
		return nil, err
	}
	return cfg, nil
}
`)
			},
		},
		{
			moduleName: "go-logger",
			dir:        "logger",
			filename:   "logger.go",
			content: func() string {
				return `package logger

import (
	gologger "github.com/mirkobrombin/go-logger"
)

var log *gologger.Logger

// Init initialises the application logger.
func Init() {
	log = gologger.NewLogger(gologger.Options{
		Level:  gologger.LevelInfo,
		Format: gologger.FormatText,
	})
}

// Get returns the shared logger instance.
func Get() *gologger.Logger {
	if log == nil {
		Init()
	}
	return log
}
`
			},
		},
		{
			moduleName: "go-module-router/v2",
			dir:        "router",
			filename:   "router.go",
			content: func() string {
				return `package router

import (
	"github.com/mirkobrombin/go-module-router/v2/pkg/router"
)

// New creates and configures the application router.
func New() *router.Router {
	r := router.New()
	// Register your modules here
	// r.Register(&mymodule.Module{})
	return r
}
`
			},
		},
		{
			moduleName: "go-relay/v2",
			dir:        "jobs",
			filename:   "jobs.go",
			content: func() string {
				broker := "Memory"
				if opts, ok := cfg.ModuleOptions["go-relay/v2"]; ok {
					if b, ok := opts["broker"]; ok {
						broker = b
					}
				}
				brokerComment := fmt.Sprintf("// Using %s broker", broker)
				return fmt.Sprintf(`package jobs

import (
	"github.com/mirkobrombin/go-relay/v2/pkg/relay"
	"github.com/mirkobrombin/go-relay/v2/pkg/broker"
)

%s

// Init sets up the relay broker and registers job handlers.
func Init() (*relay.Relay, error) {
	b := broker.NewMemoryBroker() // TODO: swap broker if needed
	r := relay.New(b)
	// Register handlers:
	// r.Register("my-job", handleMyJob)
	return r, nil
}
`, brokerComment)
			},
		},
		{
			moduleName: "go-wormhole",
			dir:        "db",
			filename:   "db.go",
			content: func() string {
				provider := "SQLite"
				if opts, ok := cfg.ModuleOptions["go-wormhole"]; ok {
					if p, ok := opts["provider"]; ok {
						provider = p
					}
				}
				return fmt.Sprintf(`package db

import (
	"github.com/fabricatorsltd/go-wormhole/pkg/wormhole"
)

// Provider: %s

// Connect opens a database connection using go-wormhole.
func Connect(dsn string) (*wormhole.DB, error) {
	return wormhole.Open(dsn)
}
`, provider)
			},
		},
	}

	for _, s := range stubs {
		if !hasModule(s.moduleName) {
			continue
		}
		dir := filepath.Join(internalDir, s.dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating dir %s: %w", dir, err)
		}
		content := s.content()
		outPath := filepath.Join(dir, s.filename)
		if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", outPath, err)
		}
		files = append(files, outPath)
	}

	// Generic stubs for other modules
	otherStubs := map[string]struct {
		dir, file, pkg, importPath string
	}{
		"go-auth":       {"auth", "auth.go", "auth", "github.com/mirkobrombin/go-auth"},
		"go-guard":      {"guard", "guard.go", "guard", "github.com/mirkobrombin/go-guard"},
		"go-secrets":    {"secrets", "secrets.go", "secrets", "github.com/mirkobrombin/go-secrets"},
		"go-signal/v2":  {"signals", "signals.go", "signals", "github.com/mirkobrombin/go-signal/v2"},
		"go-worker":     {"worker", "worker.go", "worker", "github.com/mirkobrombin/go-worker"},
		"go-state-flow": {"stateflow", "stateflow.go", "stateflow", "github.com/mirkobrombin/go-state-flow"},
		"go-lock":       {"lock", "lock.go", "lock", "github.com/mirkobrombin/go-lock"},
		"go-warp":       {"cache", "cache.go", "cache", "github.com/mirkobrombin/go-warp"},
		"go-plugin":     {"plugin", "plugin.go", "plugin", "github.com/mirkobrombin/go-plugin"},
		"go-slipstream": {"slipstream", "slipstream.go", "slipstream", "github.com/mirkobrombin/go-slipstream"},
		"go-retry":      {"retry", "retry.go", "retry", "github.com/mirkobrombin/go-retry"},
		"go-revert/v2":  {"revert", "revert.go", "revert", "github.com/mirkobrombin/go-revert/v2"},
		"go-httpx":      {"httpx", "httpx.go", "httpx", "github.com/mirkobrombin/go-httpx"},
		"go-metrics":    {"metrics", "metrics.go", "metrics", "github.com/mirkobrombin/go-metrics"},
	}

	alreadyHandled := map[string]bool{
		"go-conf-builder/v2":  true,
		"go-logger":           true,
		"go-module-router/v2": true,
		"go-relay/v2":         true,
		"go-wormhole":         true,
		"go-foundation":       true,
	}

	for modName, info := range otherStubs {
		if alreadyHandled[modName] || !hasModule(modName) {
			continue
		}
		dir := filepath.Join(internalDir, info.dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating dir %s: %w", dir, err)
		}
		content := fmt.Sprintf("package %s\n\n// TODO: implement %s logic here.\n// Import: %s\n", info.pkg, modName, info.importPath)
		outPath := filepath.Join(dir, info.file)
		if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", outPath, err)
		}
		files = append(files, outPath)
	}

	return files, nil
}

func generateFull(cfg wizard.Config) ([]string, error) {
	var files []string

	preset := catalog.Presets[cfg.PresetIndex]

	// Determine the command type based on preset
	var cmdDir, cmdFile, cmdContent string
	switch preset.Name {
	case "Worker service":
		cmdDir = "cmd/worker"
		cmdFile = "worker.go"
		cmdContent = generateWorkerCmd(cfg)
	default:
		cmdDir = "cmd/serve"
		cmdFile = "serve.go"
		cmdContent = generateServeCmd(cfg)
	}

	dir := filepath.Join(cfg.OutputDir, cmdDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating dir %s: %w", dir, err)
	}

	outPath := filepath.Join(dir, cmdFile)
	if err := os.WriteFile(outPath, []byte(cmdContent), 0644); err != nil {
		return nil, fmt.Errorf("writing %s: %w", outPath, err)
	}
	files = append(files, outPath)

	// Makefile
	makefileContent := generateMakefile(cfg)
	makefilePath := filepath.Join(cfg.OutputDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return nil, fmt.Errorf("writing Makefile: %w", err)
	}
	files = append(files, makefilePath)

	return files, nil
}

func generateServeCmd(cfg wizard.Config) string {
	return fmt.Sprintf(`package serve

import "fmt"

// Serve starts the HTTP server.
// Wire your router, middleware, and services here.
func Run() error {
	fmt.Println("Starting %s server...")
	// TODO: initialise config, logger, router
	return nil
}
`, cfg.ProjectName)
}

func generateWorkerCmd(cfg wizard.Config) string {
	return fmt.Sprintf(`package worker

import "fmt"

// Run starts the background worker service.
func Run() error {
	fmt.Println("Starting %s worker...")
	// TODO: initialise relay, register handlers, start processing
	return nil
}
`, cfg.ProjectName)
}

func generateMakefile(cfg wizard.Config) string {
	return fmt.Sprintf(`.PHONY: build run test tidy

build:
	go build -o bin/%s .

run:
	go run .

test:
	go test ./...

tidy:
	go mod tidy
`, cfg.ProjectName)
}
