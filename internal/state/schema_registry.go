package state

import (
	"fmt"
	"sync"
)

type Schema struct {
	Name   string
	Fields map[string]string
}

type SchemaBuilder interface {
	Build(*Schema) error
}

type SchemaRegistry struct {
	mu      sync.RWMutex
	schemas map[string]*Schema
}

func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{schemas: make(map[string]*Schema)}
}

// Build publishes only a fully initialized schema and converts builder panics to errors.
func (r *SchemaRegistry) Build(name string, builder SchemaBuilder) (schema *Schema, err error) {
	candidate := &Schema{Name: name, Fields: make(map[string]string)}
	defer func() {
		if recovered := recover(); recovered != nil {
			schema = nil
			err = fmt.Errorf("build schema %s: %v", name, recovered)
		}
	}()
	if err := builder.Build(candidate); err != nil {
		return nil, err
	}
	r.mu.Lock()
	r.schemas[name] = cloneSchema(candidate)
	r.mu.Unlock()
	return cloneSchema(candidate), nil
}

func cloneSchema(schema *Schema) *Schema {
	if schema == nil {
		return nil
	}
	copyValue := &Schema{Name: schema.Name, Fields: make(map[string]string, len(schema.Fields))}
	for key, value := range schema.Fields {
		copyValue.Fields[key] = value
	}
	return copyValue
}

func (r *SchemaRegistry) Get(name string) (*Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	schema, ok := r.schemas[name]
	return cloneSchema(schema), ok
}
