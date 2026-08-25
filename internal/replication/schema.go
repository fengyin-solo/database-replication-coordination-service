package replication

import "datasync/internal/state"

type SchemaService struct{ Registry *state.SchemaRegistry }

func (s *SchemaService) Refresh(name string, builder state.SchemaBuilder) (*state.Schema, error) {
	return s.Registry.Build(name, builder)
}

func (s *SchemaService) Lookup(name string) (*state.Schema, bool) {
	return s.Registry.Get(name)
}
