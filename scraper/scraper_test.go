package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
)

func TestResultsConformToSchema(t *testing.T) {
	schema := jsonschema.Schema{}
	schemaFile, err := os.ReadFile("./schema.json")
	if err != nil {
		t.Fatalf("failed to read schema file: %s", err)
	}
	err = json.Unmarshal(schemaFile, &schema)
	if err != nil {
		t.Fatalf("failed to unmarshal schema file: %s", err)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		t.Fatalf("failed to resolve schema file: %s", err)
	}

	file, err := os.ReadFile("./results.json")
	if err != nil {
		t.Log("failed to read results.json file, regenerating")
		Execute(true, true, false, "./results.json", -1)
		file, err = os.ReadFile("./results.json")
		if err != nil {
			t.Fatalf("failed to read results.json even after re-generating: %s", err)
		}
	}
	var results map[string]any
	err = json.Unmarshal(file, &results)
	if err != nil {
		t.Fatalf("failed to unmarshal results file: %s", err)
	}
	err = resolved.Validate(results)
	if err != nil {
		t.Fatalf("Results file failed validation: %s", err)
	}
}
