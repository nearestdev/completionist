package config

import (
	"os"
	"path/filepath"
)

// ResolveAppPath finds a backend-relative path across the supported run modes:
// - `go run` from apps/api
// - compiled binary launched from the monorepo root
func ResolveAppPath(rel string) string {
	candidates := []string{
		rel,
		filepath.Join("apps", "api", rel),
	}

	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		candidates = append(candidates,
			filepath.Join(execDir, "..", rel),
			filepath.Join(execDir, "..", "..", rel),
		)
	}

	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}

		if _, err := os.Stat(abs); err == nil {
			return abs
		}
	}

	return filepath.Clean(rel)
}
