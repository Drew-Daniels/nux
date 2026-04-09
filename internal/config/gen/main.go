package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/invopop/jsonschema"
)

func main() {
	outDir := "schemas"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatal(err)
	}

	generate(filepath.Join(outDir, "project.schema.json"), &config.ProjectConfig{})
	generate(filepath.Join(outDir, "global.schema.json"), &config.GlobalConfig{})
}

func generate(path string, v any) {
	r := new(jsonschema.Reflector)
	schema := r.Reflect(v)

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatalf("marshal %s: %v", path, err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
}
