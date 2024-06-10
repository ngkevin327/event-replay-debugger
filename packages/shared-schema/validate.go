package sharedschema

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed captured-event.schema.json
var schemaFS embed.FS

var compiled *jsonschema.Schema

func init() {
	raw, err := schemaFS.ReadFile("captured-event.schema.json")
	if err != nil {
		panic(fmt.Sprintf("load schema: %v", err))
	}
	c, err := jsonschema.CompileString("captured-event.schema.json", string(raw))
	if err != nil {
		panic(fmt.Sprintf("compile schema: %v", err))
	}
	compiled = c
}

// ValidateCapturedEvent checks raw JSON against the canonical CapturedEvent schema.
func ValidateCapturedEvent(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}
	if err := compiled.Validate(v); err != nil {
		return fmt.Errorf("schema validation: %w", err)
	}
	return nil
}
